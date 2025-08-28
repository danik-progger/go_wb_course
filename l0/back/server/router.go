package server

import (
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`pong`))
}

func (s *Server) SetupRoutes() {
	r := s.Router
	r.HandleFunc("/ping", ping).Methods("GET")
	r.HandleFunc("/orders", s.handleAddOrder).Methods("POST")
	r.HandleFunc("/orders/{order_uid}", s.handleGetOrderById).Methods("GET")
}
