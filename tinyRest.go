package tinyRest

import (
	"net/http"
	"strings"
)

func New(s *ResourceSet) *Resource {
	r := &Resource{
		methods: map[string]http.HandlerFunc{},
	}

	if s.Get != nil {
		r.methods[http.MethodGet] = s.Get
	}
	if s.Head != nil {
		r.methods[http.MethodHead] = s.Head
	}
	if s.Post != nil {
		r.methods[http.MethodPost] = s.Post
	}
	if s.Put != nil {
		r.methods[http.MethodPut] = s.Put
	}
	if s.Patch != nil {
		r.methods[http.MethodPatch] = s.Patch
	}
	if s.Connect != nil {
		r.methods[http.MethodConnect] = s.Connect
	}
	if s.Delete != nil {
		r.methods[http.MethodDelete] = s.Delete
	}
	if s.Trace != nil {
		r.methods[http.MethodTrace] = s.Trace
	}

	r.AllowMethods()

	return r
}

type Resource struct {
	methods map[string]http.HandlerFunc
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

	f(w, r)
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
}
