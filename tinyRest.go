package tinyRest

import (
	"net/http"
	"strings"
)

type Middlewares []func(http.Handler) http.Handler

func (m Middlewares) Use(middlewares ...func(http.Handler) http.Handler) Middlewares {
	m = append(m, middlewares...)
	return m
}

func New(s *ResourceSet) *Resource {
	r := &Resource{
		methods: map[string]http.Handler{},
	}

	if s.Get != nil {
		r.methods[http.MethodGet] = Chain(s.Get, s.Middlewares...)
	}
	if s.Head != nil {
		r.methods[http.MethodHead] = Chain(s.Head, s.Middlewares...)
	}
	if s.Post != nil {
		r.methods[http.MethodPost] = Chain(s.Post, s.Middlewares...)
	}
	if s.Put != nil {
		r.methods[http.MethodPut] = Chain(s.Put, s.Middlewares...)
	}
	if s.Patch != nil {
		r.methods[http.MethodPatch] = Chain(s.Patch, s.Middlewares...)
	}
	if s.Connect != nil {
		r.methods[http.MethodConnect] = Chain(s.Connect, s.Middlewares...)
	}
	if s.Delete != nil {
		r.methods[http.MethodDelete] = Chain(s.Delete, s.Middlewares...)
	}
	if s.Trace != nil {
		r.methods[http.MethodTrace] = Chain(s.Trace, s.Middlewares...)
	}

	r.AllowMethods()

	return r
}

type Resource struct {
	methods map[string]http.Handler
	cors    *string
}

func (s *Resource) AllowMethods() string {
	if s.cors == nil {
		c := []string{}
		for method, _ := range s.methods {
			c = append(c, method)
		}

		if f, _ := s.methods[http.MethodOptions]; f == nil {
			c = append(c, http.MethodOptions)
		}

		cors := strings.Join(c, ",")
		s.cors = &cors
	}

	return *s.cors
}

func (s *Resource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", s.AllowMethods())

	if r.Method == http.MethodOptions {
		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Length", "0")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	f, _ := s.methods[r.Method]
	if f == nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	f.ServeHTTP(w, r)
}

type ResourceSet struct {
	Get     http.HandlerFunc
	Head    http.HandlerFunc
	Post    http.HandlerFunc
	Put     http.HandlerFunc
	Patch   http.HandlerFunc
	Connect http.HandlerFunc
	Delete  http.HandlerFunc
	Trace   http.HandlerFunc

	Middlewares []func(http.Handler) http.Handler
}

func Chain(f func(http.ResponseWriter, *http.Request), middlewares ...func(http.Handler) http.Handler) http.Handler {
	var fn http.Handler
	fn = http.HandlerFunc(f)

	for _, middleware := range middlewares {
		fn = middleware(fn)
	}

	return fn
}
