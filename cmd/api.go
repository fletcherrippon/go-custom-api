package main

import (
	"encoding/json"
	net "net/http"

	"github.com/fletcherrippon/custom-http-go/pkg/http"
)

func main() {
	s := http.NewServer(":3000")

	s.Get("/", func(w net.ResponseWriter, req *net.Request) {
		req.Header.Set("Content-Type", "application/json")

		payload := map[string]string{
			"message": "hello world",
		}

		data, err := json.Marshal(payload)

		if err != nil {
			w.WriteHeader(net.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(data)
	})

	s.Post("/", func(w net.ResponseWriter, req *net.Request) {
		req.Header.Set("Content-Type", "application/json")
		w.Write([]byte("{\"message\": \"posted hello world\"}"))
	})

	s.Get("/home", func(w net.ResponseWriter, req *net.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("this is the home page"))
	})

	s.Get("/about", func(w net.ResponseWriter, req *net.Request) {
		w.Write([]byte("this is the about page"))
	})

	s.Get("/about/me", func(w net.ResponseWriter, req *net.Request) {
		w.Write([]byte("this is the about me page"))
	})

	s.Start()
}
