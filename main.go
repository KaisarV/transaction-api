package main

import (
	"log"
	"net/http"

	controller "PraktikumPBP/controller"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", controller.InsertUser).Methods("POST")
	router.HandleFunc("/users", controller.GetAllUsers).Methods("GET")
	// router.HandleFunc("/users/login", controller.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", controller.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", controller.DeleteUser).Methods("DELETE")

	router.HandleFunc("/products", controller.InsertProduct).Methods("POST")
	router.HandleFunc("/products", controller.GetAllProducts).Methods("GET")
	router.HandleFunc("/products/{id}", controller.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", controller.DeleteProduct).Methods("DELETE")

	router.HandleFunc("/transactions", controller.InsertTransaction).Methods("POST")
	router.HandleFunc("/transactions", controller.GetAllTransactions).Methods("GET")
	router.HandleFunc("transactions/user", controller.GetDetailUserTransaction).Methods("GET")
	router.HandleFunc("/transactions/{id}", controller.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{id}", controller.DeleteTransaction).Methods("DELETE")

	log.Println("Starting on Port")

	err := http.ListenAndServe(":8080", router)
	log.Fatal(err)
}
