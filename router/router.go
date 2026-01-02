package router

import (
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/stocks/{id}", middleware.GetStock).methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/stocks", middleware.GetAllStocks).methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/stocks", middleware.CreateStock).methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/stocks/{id}", middleware.UpdateStock).methods("PUT", "OPTIONS")
	router.HandleFunc("/api/v1/stocks/{id}", middleware.DeleteStock).methods("DELETE", "OPTIONS")

	return router
}
