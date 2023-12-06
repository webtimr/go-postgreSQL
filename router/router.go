package router

import (
	"github.com/gorilla/mux"
	"github.com/webtimr/go-postgreSQL.git/middleware"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/stock/{id}", middleware.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stock", middleware.GetAllStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newstock", middleware.CreateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/stock/{id}", middleware.UpdateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deletestock/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")

	return router
}
