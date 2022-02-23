package controllers

import (
	model "PraktikumPBP/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()

	defer db.Close()

	query := "SELECT * FROM products"
	id := r.URL.Query()["id"]
	if id != nil {
		query += " WHERE id = " + id[0]
	}

	rows, err := db.Query(query)
	var response model.ProductsResponse

	if err != nil {
		response.Status = 400
		response.Message = err.Error()
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var product model.Product
	var products []model.Product

	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			response.Message += err.Error() + "\n"
		} else {
			products = append(products, product)
		}
	}

	if len(products) != 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = products
	} else if response.Message == "" {
		response.Status = 400
		response.Message = "Data Not Found"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response model.ErrorResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	vars := mux.Vars(r)
	productId := vars["id"]
	query, errQuery := db.Exec(`DELETE FROM products WHERE id = ?;`, productId)
	RowsAffected, err := query.RowsAffected()

	if RowsAffected == 0 {
		response.Status = 400
		response.Message = "Product not found"
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if errQuery == nil {
		response.Status = 200
		response.Message = "Success Delete Data"
		w.WriteHeader(200)
	} else {
		response.Status = 400
		response.Message = "Error Delete Data"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var response model.ProductResponse
	err := r.ParseForm()

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	var product model.Product
	product.Name = r.Form.Get("name")
	product.Price, _ = strconv.Atoi(r.Form.Get("price"))

	log.Println(product.Name)
	log.Println(product.Price)

	_, errQuery := db.Exec("INSERT INTO products (name, price) VALUES (?,?)", product.Name, product.Price)

	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
		response.Data = product
	} else {
		response.Status = 400
		response.Message = "Error Insert Data"
		w.WriteHeader(400)
		log.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response model.ErrorResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	vars := mux.Vars(r)
	productId := vars["id"]

	data, _ := db.Query(`SELECT * FROM products WHERE id = ?;`, productId)

	if data == nil {
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", productId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var product model.Product
	product.Name = r.Form.Get("name")
	product.Price, _ = strconv.Atoi(r.Form.Get("price"))

	_, errQuery := db.Query(`UPDATE products SET name = ?, price = ? WHERE id = ?;`, product.Name, product.Price, productId)

	if errQuery == nil {
		response.Status = 200
		response.Message = "Success Update Data"
		w.WriteHeader(200)
	} else {
		response.Status = 400
		response.Message = "Error Update Data"
		w.WriteHeader(400)
		log.Println(errQuery)

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
