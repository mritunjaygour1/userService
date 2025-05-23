package main

import (
	"log"
	"net/http"
	"userService/handler"
	"userService/service"

	"github.com/gorilla/mux"
)

func main() {
	service := service.NewUserService()
	userHandler := handler.NewUserHandlerService(service)
	healthHandler := handler.NewHealthHandlerService()

	router := mux.NewRouter()

	// healthz
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods(http.MethodGet)
	// user
	router.HandleFunc("/users/v1", userHandler.CreateUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/users/v1/{id}", userHandler.CreateUserHandler).Methods(http.MethodGet)
	router.HandleFunc("/users/v1/{id}", userHandler.CreateUserHandler).Methods(http.MethodPut)
	router.HandleFunc("/users/v1/{id}", userHandler.CreateUserHandler).Methods(http.MethodDelete)

	log.Println("server started with 8080 port...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
