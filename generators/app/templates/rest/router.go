package rest

import "github.com/gorilla/mux"

func createRouter() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/", makeHandler((*handler).handleRoot))

	return r
}
