package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()

	defer db.Close()

	query := "SELECT * FROM transactions"

	rows, err := db.Query(query)
	var response TransactionsResponse

	if err != nil {
		response.Status = 400
		response.Message = "Data Not Found"
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var transaction Transaction
	var transactions []Transaction

	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductId, &transaction.Quantity); err != nil {
			log.Println(err.Error())
		} else {
			transactions = append(transactions, transaction)
		}
	}

	if len(transactions) != 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = transactions
	} else {
		response.Status = 400
		response.Message = "Error Array Size Not Correct"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response ErrorResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	vars := mux.Vars(r)
	transactionId := vars["id"]
	data, _ := db.Query(`SELECT * FROM transactions WHERE id = ?;`, transactionId)

	if data == nil {
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", transactionId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	_, errQuery := db.Query(`DELETE FROM transactions WHERE id = ?;`, transactionId)

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

func InsertTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var response TransactionResponse
	err := r.ParseForm()

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	var transaction Transaction

	transaction.UserID, _ = strconv.Atoi(r.Form.Get("userid"))
	transaction.ProductId, _ = strconv.Atoi(r.Form.Get("productid"))
	transaction.Quantity, _ = strconv.Atoi(r.Form.Get("qty"))

	_, errQuery := db.Exec("INSERT INTO transactions (userid, productid, quantity) VALUES (?,?,?)", transaction.UserID, transaction.ProductId, transaction.Quantity)

	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
		response.Data = transaction
	} else {
		response.Status = 400
		response.Message = "Error Insert Data"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response ErrorResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	vars := mux.Vars(r)
	transactionId := vars["id"]

	data, _ := db.Query(`SELECT * FROM transactions WHERE id = ?;`, transactionId)

	if data == nil {
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", transactionId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var transaction Transaction
	transaction.UserID, _ = strconv.Atoi(r.Form.Get("userid"))
	transaction.ProductId, _ = strconv.Atoi(r.Form.Get("productid"))
	transaction.Quantity, _ = strconv.Atoi(r.Form.Get("qty"))

	_, errQuery := db.Query(`UPDATE transactions SET userid = ?, productid = ?, quantity = ? WHERE id = ?;`, transaction.UserID, transaction.ProductId, transaction.Quantity, transactionId)

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
