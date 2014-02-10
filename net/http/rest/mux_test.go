package rest

import (
	"github.com/rafaeljusto/shelter/net/http/rest/context"
	"github.com/rafaeljusto/shelter/net/http/rest/handler"
	"net/http"
	"regexp"
	"testing"
)

func TestFindRoute(t *testing.T) {
	caller := 0
	handler.Routes = make(map[*regexp.Regexp]handler.Handler)

	handler.HandleFunc(regexp.MustCompile(`^/domain/([[:alnum:]]|\-|\.)+$`), func(r *http.Request, context *context.Context) {
		caller = 1
	})

	handler.HandleFunc(regexp.MustCompile("^/domains(/)?.*$"), func(r *http.Request, context *context.Context) {
		caller = 2
	})

	handler.HandleFunc(regexp.MustCompile(`^/domain/([[:alnum:]]|\-|\.)+/verification$`), func(r *http.Request, context *context.Context) {
		caller = 3
	})

	data := []struct {
		URI      string
		CallerId int
	}{
		{URI: "/domain/example.com.br.", CallerId: 1},
		{URI: "/domains", CallerId: 2},
		{URI: "/domains/", CallerId: 2},
		{URI: "/domains/?page=1&pagesize=10&orderby=fqdn:desc@lastmodified:asc", CallerId: 2},
		{URI: "/domain/example.com.br./verification", CallerId: 3},
	}

	for _, item := range data {
		handler := mux.findRoute(item.URI)

		if handler == nil {
			t.Fatal("Did not found a valid route for %s", item.URI)
		}

		handler(nil, nil)
		if caller != item.CallerId {
			t.Errorf("Not calling the correct handler for %s", item.URI)
		}
	}
}

func TestUnknownRoute(t *testing.T) {
	handler.Routes = make(map[*regexp.Regexp]handler.Handler)
	handler.HandleFunc(regexp.MustCompile(`^/domains$`), func(r *http.Request, context *context.Context) {
	})

	if mux.findRoute("/domain") != nil {
		t.Error("Not returning nil when there's no handler for the URI")
	}
}

func BenchmarkFindRoute(b *testing.B) {
	handler.Routes = make(map[*regexp.Regexp]handler.Handler)

	handler.HandleFunc(regexp.MustCompile(`^/domain/([[:alnum:]]|\-|\.)+$`), func(r *http.Request, context *context.Context) {})

	for i := 0; i < b.N; i++ {
		mux.findRoute("/domain/example.com.br.")
	}
}
