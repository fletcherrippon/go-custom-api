package http

import (
	"net/http"
	"strings"
)

type Request struct {
	Method string
	Path   string
}

type Router struct {
	root *Route
}

func NewRouter() *Router {
	return &Router{
		root: NewRoute(),
	}
}

func (r *Router) SplitPath(path string) []string {
	return strings.Split(path, "/")
}

func (r *Router) AddRoute(method, path string, handler http.HandlerFunc) {
	segments := strings.Split(path, "/")
	node := r.root

	for _, segment := range segments {
		if segment == "" {
			continue
		}
		child, ok := node.children[segment]
		if !ok {
			child = NewRoute()
			node.children[segment] = child
		}
		node = child
	}

	node.handlers[method] = handler
}

func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.AddRoute("GET", path, handler)
}

func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.AddRoute("POST", path, handler)
}

func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.AddRoute("PUT", path, handler)
}

func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.AddRoute("DELETE", path, handler)
}

func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.AddRoute("PATCH", path, handler)
}

func (r *Router) Head(path string, handler http.HandlerFunc) {
	r.AddRoute("HEAD", path, handler)
}

func (r *Router) Options(path string, handler http.HandlerFunc) {
	r.AddRoute("OPTIONS", path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	segments := strings.Split(req.URL.Path, "/")
	node := r.root

	for _, segment := range segments {
		if segment == "" {
			continue
		}
		child, ok := node.children[segment]
		if !ok {
			http.NotFoundHandler().ServeHTTP(w, req)
			return
		}
		node = child
	}

	handler, ok := node.handlers[req.Method]

	if !ok {
		http.NotFoundHandler().ServeHTTP(w, req)
		return
	}

	handler(w, req)
}
