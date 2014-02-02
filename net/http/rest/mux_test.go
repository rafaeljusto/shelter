package rest

import (
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"github.com/rafaeljusto/shelter/net/http/rest/handler"
	"net/http"
	"testing"
)

func TestFindRoute(t *testing.T) {
	caller := 0
	handler.Routes = make(map[string]handler.Handler)

	handler.HandleFunc("/domain/", func(r *http.Request, context *context.Context) {
		caller = 1
	})

	handler.HandleFunc("/domains", func(r *http.Request, context *context.Context) {
		caller = 2
	})

	var mux Mux

	uri := "/domain/example.com.br."
	handler := mux.findRoute(uri)

	if handler == nil {
		t.Fatal("Did not found a valid route")
	}

	handler(nil, nil)
	if caller != 1 {
		t.Error("Not calling the correct handler")
	}
}

func BenchmarkFindRoute(b *testing.B) {
	handler.Routes = make(map[string]handler.Handler)

	handler.HandleFunc("/domain/", func(r *http.Request, context *context.Context) {})

	for i := 0; i < b.N; i++ {
		mux.findRoute("/domain/example.com.br.")
	}
}
