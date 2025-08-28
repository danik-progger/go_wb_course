package server

import (
	"encoding/json"
	"fmt"
	"io"
	"l0/cache"
	"l0/db"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Port   string
	Router *mux.Router
	DB     *db.DB
	Cache  *cache.Cache
}

func InitServer(port string) *Server {
	database := db.InitDB()
	fmt.Println("🟢 Database connected & migrated")

	router := mux.NewRouter()
	fmt.Println("🟢 Router started")

	// CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	fmt.Println("🟢 CORS middleware added")

	cache := cache.InitCache()
	fmt.Println("🟢 Cache initialized")

	s := &Server{
		Port:   port,
		Router: router,
		DB:     database,
		Cache:  cache,
	}

	s.fillCache()
	fmt.Println("🟢 Cache filled with data from db")

	s.SetupRoutes()
	fmt.Println("🟢 Routes set up")

	return s
}

func (s *Server) fillCache() {
	dbNotEmpty, err := s.DB.HasOrders()
	if err != nil {
		fmt.Printf("🔴 Failed check if DB is empty: %s\n", err)
	}
	if dbNotEmpty {
		data, err := s.DB.GetAllOrders(10, 0)
		if err != nil {
			fmt.Printf("🔴 Failed to get orders from DB: %s\n", err)
		}

		for _, o := range data {
			b, err := json.Marshal(o)
			if err != nil {
				fmt.Printf("🔴 Failed to marshal order for cache uid=%s: %s\n", o.OrderUID, err)
				continue
			}
			s.Cache.Set(o.OrderUID, string(b))
			fmt.Printf("🟢 Cached order uid=%s\n", o.OrderUID)
		}
	}
}

func (s *Server) handleAddOrder(w http.ResponseWriter, r *http.Request) {
	var order db.Order
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("🔴 Error reading request body: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("🟢 AddOrder raw body: %s\n", string(bodyBytes))

	var env struct {
		Action string          `json:"action"`
		Body   json.RawMessage `json:"body"`
	}
	orderBytes := bodyBytes
	if err := json.Unmarshal(bodyBytes, &env); err == nil && len(env.Body) > 0 {
		orderBytes = env.Body
		fmt.Printf("🟢 Detected envelope with action=%s, using body\n", env.Action)
	}

	if err := json.Unmarshal(orderBytes, &order); err != nil {
		fmt.Printf("🔴 Error in AddOrder decode: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("🟢 Decoded")

	if s.Cache.Has(order.OrderUID) {
		fmt.Printf("🟡 Cache: order already cached uid=%s\n", order.OrderUID)
	} else {
		if err := s.DB.AddOrder(&order); err != nil {
			fmt.Printf("🔴 Error in AddOrder: %s\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(order)
		if err != nil {
			fmt.Printf("🔴 Failed to marshal order for cache: %s\n", err)
		} else {
			s.Cache.Set(order.OrderUID, string(b))
			fmt.Printf("🟢 Cached added order uid=%s\n", order.OrderUID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
	fmt.Println("🟢 Successful request AddOrder")
}

func (s *Server) handleGetOrderById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["order_uid"]

	// Try cache first
	if val, ok := s.Cache.Get(id); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(val))
		fmt.Printf("🟢 Cache hit for uid=%s\n", id)
		return
	}
	fmt.Printf("🟡 Cache miss for uid=%s\n", id)

	orders, err := s.DB.GetOrderById(id)
	if orders == nil {
		fmt.Printf("🔴 Error in GetOrderById: %s\n", err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Printf("🔴 Error in GetOrderById: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !s.Cache.Has(orders.OrderUID) {
		if b, err := json.Marshal(orders); err == nil {
			s.Cache.Set(orders.OrderUID, string(b))
			fmt.Printf("🟢 Cached order uid=%s after DB fetch\n", orders.OrderUID)
		} else {
			fmt.Printf("🔴 Failed to marshal order for cache: %s\n", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
	fmt.Println("🟢 Successful request GetOrderById")
}
