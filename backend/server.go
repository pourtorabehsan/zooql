package main

import (
	"encoding/json"
	"net/http"
)

type server struct {
	zookeepers string
	basePath   string
	znodes     *znodeRepository
	router     *http.ServeMux
}

func newServer(repo *znodeRepository, zookeepers string, basePath string) *server {
	s := &server{
		zookeepers: zookeepers,
		basePath:   basePath,
		znodes:     repo,
		router:     http.NewServeMux(),
	}

	s.router.HandleFunc("/backend/connection", s.handleConnection())
	s.router.HandleFunc("/backend/query", s.handleQuery())

	return s
}

func (s *server) Start() error {
	return http.ListenAndServe(":8000", s.router)
}

func (s *server) handleConnection() http.HandlerFunc {
	type connection struct {
		Zookeepers string `json:"zookeepers"`
		BasePath   string `json:"basePath"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resp := &connection{
			Zookeepers: s.zookeepers,
			BasePath:   s.basePath,
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func (s *server) handleQuery() http.HandlerFunc {
	type queryRequest struct {
		Query string `json:"query"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := &queryRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := s.znodes.query(req.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp)
	}
}
