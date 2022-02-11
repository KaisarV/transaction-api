package main

import (
	"log"
	"net/http"

	controllers "PraktikumPBP/controllers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", controllers.InsertUser).Methods("POST")
	router.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE")

	router.HandleFunc("/products", controllers.InsertProduct).Methods("POST")
	router.HandleFunc("/products", controllers.GetAllProducts).Methods("GET")
	router.HandleFunc("/products/{id}", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", controllers.DeleteProduct).Methods("DELETE")

	router.HandleFunc("/transactions", controllers.InsertTransaction).Methods("POST")
	router.HandleFunc("/transactions", controllers.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions/{id}", controllers.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{id}", controllers.DeleteTransaction).Methods("DELETE")

	log.Println("Starting on Port")

	err := http.ListenAndServe(":8080", router)
	log.Fatal(err)
}
