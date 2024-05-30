package http

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type Server struct {
	router *Router
	addr   string
}

func NewServer(addr string) *Server {
	return &Server{
		router: NewRouter(),
		addr:   addr,
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		requestLine, err := reader.ReadString('\n')

		if err != nil {
			if err.Error() == "EOF" {
				return
			}

			fmt.Printf("Failed to read request line: %v\n", err)

			return
		}

		parts := strings.Split(strings.TrimSpace(requestLine), " ")

		method := parts[0]
		path := parts[1]

		req, err := http.NewRequest(method, path, nil)

		if err != nil {
			fmt.Printf("Failed to create request: %v\n", err)
			return
		}

		w := NewResponseWriter(conn)

		s.router.ServeHTTP(w, req)
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		fmt.Printf("Server startup failed: %v\n", err)
		return
	}

	defer listener.Close()

	fmt.Printf("Server is running on %s\n", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) AddRoute(method, path string, handler http.HandlerFunc) {
	s.router.AddRoute(method, path, handler)
}

func (s *Server) Get(path string, handler http.HandlerFunc) {
	s.router.Get(path, handler)
}

func (s *Server) Post(path string, handler http.HandlerFunc) {
	s.router.Post(path, handler)
}

func (s *Server) Put(path string, handler http.HandlerFunc) {
	s.router.Put(path, handler)
}

func (s *Server) Delete(path string, handler http.HandlerFunc) {
	s.router.Delete(path, handler)
}

func (s *Server) Patch(path string, handler http.HandlerFunc) {
	s.router.Patch(path, handler)
}

func (s *Server) Head(path string, handler http.HandlerFunc) {
	s.router.Head(path, handler)
}

func (s *Server) Options(path string, handler http.HandlerFunc) {
	s.router.Options(path, handler)
}

type ResponseWriter struct {
	conn          net.Conn
	headers       http.Header
	statusCode    int
	headerWritten bool
}

func NewResponseWriter(conn net.Conn) *ResponseWriter {
	return &ResponseWriter{
		conn:          conn,
		headers:       make(http.Header),
		statusCode:    http.StatusOK,
		headerWritten: false,
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.headers
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	if !w.headerWritten {
		w.WriteHeader(w.statusCode)
	}
	return w.conn.Write(data)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	if w.headerWritten {
		return
	}
	w.statusCode = statusCode

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))
	w.conn.Write([]byte(statusLine))

	for key, values := range w.headers {
		for _, value := range values {
			header := fmt.Sprintf("%s: %s\r\n", key, value)
			w.conn.Write([]byte(header))
		}
	}

	w.conn.Write([]byte("\r\n"))
	w.headerWritten = true
}
