package httpServer

import (
	"fmt"
	"net/http"
	"wb_l0/cache"
)

type Server struct {
	cache *cache.Cache
}

func InitServer(cache *cache.Cache) *Server {
	return &Server{
		cache: cache,
	}
}

func (s *Server) StartServer() error {
	http.HandleFunc("/", s.HandleSearch)
	//http.HandleFunc("/order", s.HandleSearch)

	fmt.Println("Server is listening...")
	return http.ListenAndServe(":8181", nil)
}
