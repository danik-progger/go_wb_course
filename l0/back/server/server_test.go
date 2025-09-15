package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"l0/db"
	"l0/mocks"

	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (*Server, *mocks.MockDBInterface, *mocks.MockCacheInterface) {
	ctrl := gomock.NewController(t)
	mockDB := mocks.NewMockDBInterface(ctrl)
	mockCache := mocks.NewMockCacheInterface(ctrl)

	router := mux.NewRouter()
	s := &Server{
		Router: router,
		DB:     mockDB,
		Cache:  mockCache,
	}
	s.SetupRoutes()

	return s, mockDB, mockCache
}

func TestHandleGetOrderById_CacheHit(t *testing.T) {
	s, _, mockCache := setup(t)

	order := db.Order{OrderUID: "test-uid"}
	orderJSON, _ := json.Marshal(order)

	mockCache.EXPECT().Get("test-uid").Return(string(orderJSON), true)

	req, _ := http.NewRequest("GET", "/orders/test-uid", nil)
	rr := httptest.NewRecorder()

	s.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	if rr.Body.String() != string(orderJSON) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(orderJSON))
	}
}

func TestHandleGetOrderById_CacheMiss_DBHit(t *testing.T) {
	s, mockDB, mockCache := setup(t)

	order := &db.Order{OrderUID: "test-uid"}
	orderJSON, _ := json.Marshal(order)

	mockCache.EXPECT().Get("test-uid").Return("", false)
	mockDB.EXPECT().GetOrderById("test-uid").Return(order, nil)
	mockCache.EXPECT().Has("test-uid").Return(false)
	mockCache.EXPECT().Set("test-uid", string(orderJSON))

	req, _ := http.NewRequest("GET", "/orders/test-uid", nil)
	rr := httptest.NewRecorder()

	s.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	if rr.Body.String() != string(orderJSON)+"\n" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(orderJSON)+"\n")
	}
}

func TestHandleGetOrderById_NotFound(t *testing.T) {
	s, mockDB, mockCache := setup(t)

	mockCache.EXPECT().Get("test-uid").Return("", false)
	mockDB.EXPECT().GetOrderById("test-uid").Return(nil, nil)

	req, _ := http.NewRequest("GET", "/orders/test-uid", nil)
	rr := httptest.NewRecorder()

	s.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusNotFound)
	}
}

func TestHandleAddOrder_NewOrder(t *testing.T) {
	s, mockDB, mockCache := setup(t)

	order := db.Order{OrderUID: "test-uid"}
	orderJSON, _ := json.Marshal(order)

	mockCache.EXPECT().Has("test-uid").Return(false)
	mockDB.EXPECT().AddOrder(gomock.Any()).Return(nil)
	mockCache.EXPECT().Set("test-uid", gomock.Any())

	req, _ := http.NewRequest("POST", "/orders", strings.NewReader(string(orderJSON)))
	rr := httptest.NewRecorder()

	s.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusCreated)
	}
}

func TestHandleAddOrder_ExistingOrder(t *testing.T) {
	s, _, mockCache := setup(t)

	order := db.Order{OrderUID: "test-uid"}
	orderJSON, _ := json.Marshal(order)

	mockCache.EXPECT().Has("test-uid").Return(true)

	req, _ := http.NewRequest("POST", "/orders", strings.NewReader(string(orderJSON)))
	rr := httptest.NewRecorder()

	s.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusCreated)
	}
}
