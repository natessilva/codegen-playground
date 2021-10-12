package apimux

import (
	"fmt"
	"net/http"
)

type Server struct {
	*http.ServeMux
}

func NewServer() *Server {
	return &Server{http.NewServeMux()}
}

func (s *Server) Register(service, method string, handler http.HandlerFunc) {
	s.Handle(fmt.Sprintf("/%s.%s", service, method), handler)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	s.ServeMux.ServeHTTP(w, r)
}
