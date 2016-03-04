package rest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type handler struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	query          url.Values
	body           io.ReadCloser
}

type handlerMethod func(*handler) error

func makeHandler(method handlerMethod) http.Handler {
	return http.HandlerFunc(func(r http.ResponseWriter, rq *http.Request) {
		h := newHandler(r, rq)
		err := h.invoke(method)
		if err != nil {
			log.WithField("RequestURI", rq.RequestURI).Error(err)
		}
	})
}

func newHandler(w http.ResponseWriter, rq *http.Request) *handler {
	return &handler{
		responseWriter: w,
		request:        rq,
		query:          rq.URL.Query(),
		body:           rq.Body,
	}
}

func (h *handler) invoke(method handlerMethod) error {
	return method(h)
}

func (h *handler) writeJSON(i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	io.WriteString(h.responseWriter, string(b))

	return nil
}

func (h *handler) writeText(text string) {
	io.WriteString(h.responseWriter, text)
}
