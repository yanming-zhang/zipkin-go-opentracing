// +build go1.7

package svc2

import (
	"fmt"
	"net/http"
	"strconv"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/yanming-zhang/zipkin-go-opentracing/examples/middleware"
)

type httpService struct {
	service Service
}

// sumHandler is our HTTP Handlerfunc for a Sum request.
func (s *httpService) sumHandler(w http.ResponseWriter, req *http.Request) {
	// parse query parameters
	v := req.URL.Query()
	a, err := strconv.ParseInt(v.Get("a"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := strconv.ParseInt(v.Get("b"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// call our Sum binding
	result, err := s.service.Sum(req.Context(), a, b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// return the result
	w.Write([]byte(fmt.Sprintf("%d", result)))
}

// NewHTTPHandler returns a new HTTP handler our svc2.
func NewHTTPHandler(tracer opentracing.Tracer, service Service) http.Handler {
	// Create our HTTP Service.
	svc := &httpService{service: service}

	// Create the mux.
	mux := http.NewServeMux()

	// Create the Sum handler.
	var sumHandler http.Handler
	sumHandler = http.HandlerFunc(svc.sumHandler)

	// Wrap the Sum handler with our tracing middleware.
	sumHandler = middleware.FromHTTPRequest(tracer, "Sum")(sumHandler)

	// Wire up the mux.
	mux.Handle("/sum/", sumHandler)

	// Return the mux.
	return mux
}
