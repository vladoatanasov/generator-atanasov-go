package rest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"github.com/gorilla/handlers"
	"github.com/vladoatanasov/logrus_amqp"
	"gopkg.in/couchbaselabs/gocb.v1"
)

var (
	log = logrus.New()
)

//Config is the serialized config.json file
type Config struct {
	SyncEndpoint   string `json:"syncEndpoint"`
	CBEndpoint     string `json:"cbEndpoint"`
	Port           int    `json:"port"`
	Bucket         string `json:"bucket"`
	BucketPassword string `json:"bucketPassword"`
	Amqp           AMQP   `json:"amqp"`
	Cors           bool   `json:"cors"`
	SSLCert        string `json:"SSLCert"`
	SSLKey         string `json:"SSLKey"`
	SSLPort        int    `json:"sslPort"`
}

// AMQP holds connection data for rabbitMQ
type AMQP struct {
	Server     string `json:"server"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routingKey"`
}

var (
	configFileDescription = "The name of the config file.  Defaults to 'config.json'"
	configFileName        = kingpin.Arg("config file name", configFileDescription).Default("config.json").String()
	config                Config
	bucket                *gocb.Bucket
)

func init() {
	// parse config file
	kingpin.Parse()
	if *configFileName == "" {
		log.Panic("Config file name missing")
		return
	}
	configFile, err := os.Open(*configFileName)
	if err != nil {
		log.Panic(err)
	}
	defer configFile.Close()

	configReader := bufio.NewReader(configFile)
	err = parseConfigFile(configReader)
	if err != nil {
		log.Panic(err)
	}

	log.Level = logrus.DebugLevel

	if (config.Amqp != AMQP{}) {
		hook := logrus_amqp.NewAMQPHook(config.Amqp.Server, config.Amqp.Username, config.Amqp.Password, config.Amqp.Exchange, config.Amqp.RoutingKey)
		log.Hooks.Add(hook)
		log.Info("Forwarding logs to RabbitMQ")
	}

	// init couchbase connection
	cluster, err := gocb.Connect(config.CBEndpoint)
	if err != nil {
		log.WithFields(logrus.Fields{
			"context": "couchbase",
			"topic":   "connect",
		}).Error(err)
	}

	bucket, err = cluster.OpenBucket(config.Bucket, config.BucketPassword)
	if err != nil {
		log.WithFields(logrus.Fields{
			"context": "couchbase",
			"topic":   "open bucket",
		}).Error(err)
	}

}

func parseConfigFile(r io.Reader) error {
	config = Config{}

	decoder := json.NewDecoder(r)

	if err := decoder.Decode(&config); err != nil {
		return err
	}

	return nil
}

//StartServer ...
func StartServer() {
	router := createRouter()

	http.Handle("/", router)

	var cors http.Handler
	var message string

	if config.Cors {
		cors = handlers.CORS()(router)
		message = " and CORS enabled"
	}

	if config.SSLKey != "" && config.SSLCert != "" && config.SSLPort > 0 {
		go func() {
			log.Infof("Starting server on port %d with SSL %s", config.SSLPort, message)
			err := http.ListenAndServeTLS(fmt.Sprintf(":%d", config.SSLPort), config.SSLCert, config.SSLKey, cors)
			if err != nil {
				log.Fatal(err)
			}

		}()
	}

	log.Infof("Starting server on port %d %s", config.Port, message)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), cors)
	if err != nil {
		log.Fatal(err)
	}

}
