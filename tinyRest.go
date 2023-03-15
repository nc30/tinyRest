package tinyRest

import (
	"net/http"
	"strings"
)

var AccessControlDefault = "*"

type ResourceSet struct {
	Get    http.HandlerFunc
	Post   http.HandlerFunc
	Head   http.HandlerFunc
	Put    http.HandlerFunc
	Patch  http.HandlerFunc
	Delete http.HandlerFunc

	cors *string
}

func (s *ResourceSet) AllowMethods() string {
	if s.cors == nil {
		c := []string{}
		if s.Get != nil {
			c = append(c, http.MethodGet)
		}
		if s.Post != nil {
			c = append(c, http.MethodPost)
		}
		if s.Head != nil {
			c = append(c, http.MethodHead)
		}
		if s.Put != nil {
			c = append(c, http.MethodPut)
		}
		if s.Patch != nil {
			c = append(c, http.MethodPatch)
		}
		if s.Delete != nil {
			c = append(c, http.MethodDelete)
		}

		c = append(c, http.MethodOptions)

		cors := strings.Join(c, ",")
		s.cors = &cors
	}

	return *s.cors
}

func (s *ResourceSet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", AccessControlDefault)
	w.Header().Add("Access-Control-Allow-Methods", s.AllowMethods())

	switch r.Method {
	case http.MethodOptions:
		w.Header().Add("Content-Type", "text/plain")
		w.Header().Add("Content-Length", "0")
		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodGet:
		if s.Get == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		s.Get(w, r)
	case http.MethodPost:
		if s.Post == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		s.Post(w, r)
	case http.MethodHead:
		if s.Head == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		s.Head(w, r)
	case http.MethodPut:
		if s.Put == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		s.Put(w, r)
	case http.MethodPatch:
		if s.Patch == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		s.Patch(w, r)
	case http.MethodDelete:
		if s.Delete == nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		s.Delete(w, r)
	}
}
