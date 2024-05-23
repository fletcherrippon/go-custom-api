package http

import (
	"net/http"
)

type Route struct {
	children map[string]*Route
	handlers map[string]http.HandlerFunc
}

func NewRoute() *Route {
	return &Route{
		children: make(map[string]*Route),
		handlers: make(map[string]http.HandlerFunc),
	}
}
