package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"salesTracker/db"
	"salesTracker/db/repo"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Server struct {
	Router  *chi.Mux
	salesDB *db.DB
}

func NewServer(salesDB *db.DB) *Server {
	s := &Server{
		Router:  chi.NewRouter(),
		salesDB: salesDB,
	}

	// CRUD operations
	s.Router.Get("/items", s.getItemsHandler)
	s.Router.Get("/items/{id}", s.getItemHandler)
	s.Router.Post("/items", s.createItemHandler)
	s.Router.Put("/items/{id}", s.updateItemHandler)
	s.Router.Delete("/items/{id}", s.deleteItemHandler)

	// Analytics
	s.Router.Get("/analytics", s.getAnalyticsHandler)
	s.Router.Get("/analytics/avg", s.getAverageHandler)
	s.Router.Get("/analytics/sum", s.getSumHandler)
	s.Router.Get("/analytics/median", s.getMedianHandler)
	s.Router.Get("/analytics/percentile", s.getPercentileHandler)

	return s
}

func (s *Server) createItemHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Category string  `json:"category"`
		Date     string  `json:"date"`
		Amount   float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse date
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	date := pgtype.Date{Time: parsedDate, Valid: true}
	amount := pgtype.Numeric{}
	if err := amount.Scan(req.Amount); err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	sale, err := s.salesDB.CreateSale(req.Category, date, amount)
	if err != nil {
		http.Error(w, "Error creating item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sale)
}

func (s *Server) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	sales, err := s.salesDB.GetAll()
	if err != nil {
		http.Error(w, "Error fetching items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sales)
}

func (s *Server) getItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	sale, err := s.salesDB.GetById(id)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sale)
}

func (s *Server) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Category string  `json:"category"`
		Date     string  `json:"date"`
		Amount   float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse date
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	date := pgtype.Date{Time: parsedDate, Valid: true}
	amount := pgtype.Numeric{}
	if err := amount.Scan(req.Amount); err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	sale, err := s.salesDB.UpdateItem(id, req.Category, date, amount)
	if err != nil {
		http.Error(w, "Error updating item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sale)
}

func (s *Server) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	err = s.salesDB.DeleteItem(id)
	if err != nil {
		http.Error(w, "Error deleting item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Item deleted successfully")
}

func (s *Server) getAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	var from, to pgtype.Date
	var err error

	// Parse from date if provided
	if fromDate != "" {
		parsedFrom, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			http.Error(w, "Invalid from date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		from = pgtype.Date{Time: parsedFrom, Valid: true}
	} else {
		// If no date range is provided, use all records
		sales, err := s.salesDB.GetAll()
		if err != nil {
			http.Error(w, "Error fetching analytics data", http.StatusInternalServerError)
			return
		}

		// Calculate analytics
		count := len(sales)

		// Calculate sum
		var sum float64
		for _, sale := range sales {
			if sale.Amount.Valid {
				// Convert pgtype.Numeric to float64
				float64Amount, err := sale.Amount.Float64Value()
				if err == nil {
					sum += float64Amount.Float64
				}
			}
		}

		// Calculate average
		var avg float64
		if count > 0 {
			avg = sum / float64(count)
		}

		// Prepare response
		response := map[string]interface{}{
			"count": count,
			"sum":   sum,
			"avg":   avg,
			"items": sales,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse to date if provided
	if toDate != "" {
		parsedTo, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			http.Error(w, "Invalid to date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		to = pgtype.Date{Time: parsedTo, Valid: true}
	}

	// If both from and to dates are provided, get date interval data
	var sales []repo.Sale
	if from.Valid && to.Valid {
		sales, err = s.salesDB.GetDateInterval(from, to)
		if err != nil {
			http.Error(w, "Error fetching analytics data", http.StatusInternalServerError)
			return
		}
	} else {
		// If only from date is provided, use all records from that date
		sales, err = s.salesDB.GetAll()
		if err != nil {
			http.Error(w, "Error fetching analytics data", http.StatusInternalServerError)
			return
		}
		// Filter the results to only include records from the from date
		filteredSales := make([]repo.Sale, 0)
		for _, sale := range sales {
			if sale.Date.Valid && !sale.Date.Time.Before(from.Time) {
				filteredSales = append(filteredSales, sale)
			}
		}
		sales = filteredSales
	}

	// Calculate analytics
	count := len(sales)

	// Calculate sum
	var sum float64
	for _, sale := range sales {
		if sale.Amount.Valid {
			// Convert pgtype.Numeric to float64
			float64Amount, err := sale.Amount.Float64Value()
			if err == nil {
				sum += float64Amount.Float64
			}
		}
	}

	// Calculate average
	var avg float64
	if count > 0 {
		avg = sum / float64(count)
	}

	// Prepare response
	response := map[string]interface{}{
		"count": count,
		"sum":   sum,
		"avg":   avg,
		"items": sales,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) getAverageHandler(w http.ResponseWriter, r *http.Request) {
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	var from, to pgtype.Date
	var err error

	// Parse from date if provided
	if fromDate != "" {
		parsedFrom, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			http.Error(w, "Invalid from date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		from = pgtype.Date{Time: parsedFrom, Valid: true}
	} else {
		// Default to a very old date to include all records if no from date is provided
		from = pgtype.Date{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	}

	// Parse to date if provided
	if toDate != "" {
		parsedTo, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			http.Error(w, "Invalid to date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		to = pgtype.Date{Time: parsedTo, Valid: true}
	} else {
		// Default to a very future date to include all records if no to date is provided
		to = pgtype.Date{Time: time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC), Valid: true}
	}

	avg, err := s.salesDB.GetAverage(from, to)
	if err != nil {
		http.Error(w, "Error calculating average", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"avg": avg})
}

func (s *Server) getSumHandler(w http.ResponseWriter, r *http.Request) {
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	var from, to pgtype.Date
	var err error

	// Parse from date if provided
	if fromDate != "" {
		parsedFrom, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			http.Error(w, "Invalid from date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		from = pgtype.Date{Time: parsedFrom, Valid: true}
	} else {
		// Default to a very old date to include all records if no from date is provided
		from = pgtype.Date{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	}

	// Parse to date if provided
	if toDate != "" {
		parsedTo, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			http.Error(w, "Invalid to date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		to = pgtype.Date{Time: parsedTo, Valid: true}
	} else {
		// Default to a very future date to include all records if no to date is provided
		to = pgtype.Date{Time: time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC), Valid: true}
	}

	sum, err := s.salesDB.GetSum(from, to)
	if err != nil {
		http.Error(w, "Error calculating sum", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"sum": sum})
}

func (s *Server) getMedianHandler(w http.ResponseWriter, r *http.Request) {
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	var from, to pgtype.Date
	var err error

	// Parse from date if provided
	if fromDate != "" {
		parsedFrom, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			http.Error(w, "Invalid from date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		from = pgtype.Date{Time: parsedFrom, Valid: true}
	} else {
		// Default to a very old date to include all records if no from date is provided
		from = pgtype.Date{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	}

	// Parse to date if provided
	if toDate != "" {
		parsedTo, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			http.Error(w, "Invalid to date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		to = pgtype.Date{Time: parsedTo, Valid: true}
	} else {
		// Default to a very future date to include all records if no to date is provided
		to = pgtype.Date{Time: time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC), Valid: true}
	}

	median, err := s.salesDB.GetMedian(from, to)
	if err != nil {
		http.Error(w, "Error calculating median", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"median": median})
}

func (s *Server) getPercentileHandler(w http.ResponseWriter, r *http.Request) {
	percentileStr := r.URL.Query().Get("value")
	if percentileStr == "" {
		percentileStr = "0.9" // default to 90th percentile
	}

	percentile, err := strconv.ParseFloat(percentileStr, 64)
	if err != nil {
		http.Error(w, "Invalid percentile value", http.StatusBadRequest)
		return
	}

	if percentile < 0 || percentile > 1 {
		http.Error(w, "Percentile must be between 0 and 1", http.StatusBadRequest)
		return
	}

	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")

	var from, to pgtype.Date

	// Parse from date if provided
	if fromDate != "" {
		parsedFrom, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			http.Error(w, "Invalid from date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		from = pgtype.Date{Time: parsedFrom, Valid: true}
	} else {
		// Default to a very old date to include all records if no from date is provided
		from = pgtype.Date{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	}

	// Parse to date if provided
	if toDate != "" {
		parsedTo, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			http.Error(w, "Invalid to date format, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		to = pgtype.Date{Time: parsedTo, Valid: true}
	} else {
		// Default to a very future date to include all records if no to date is provided
		to = pgtype.Date{Time: time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC), Valid: true}
	}

	value, err := s.salesDB.GetPercentile(percentile, from, to)
	if err != nil {
		http.Error(w, "Error calculating percentile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"percentile": value})
}
