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

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()

	defer db.Close()

	query := "SELECT * FROM transactions"
	id := r.URL.Query()["id"]
	if id != nil {
		query += " WHERE id = " + id[0]
	}

	rows, err := db.Query(query)
	var response model.TransactionsResponse

	if err != nil {
		response.Status = 400
		response.Message = err.Error()
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var transaction model.Transaction
	var transactions []model.Transaction

	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductId, &transaction.Quantity); err != nil {
			response.Message += err.Error() + "\n"
		} else {
			transactions = append(transactions, transaction)
		}
	}

	if len(transactions) != 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = transactions
	} else if response.Message == "" {
		response.Status = 400
		response.Message = "Data Not Found"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
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
	transactionId := vars["id"]
	query, errQuery := db.Exec(`DELETE FROM transactions WHERE id = ?;`, transactionId)
	RowsAffected, err := query.RowsAffected()

	if RowsAffected == 0 {
		response.Status = 400
		response.Message = "Transaction not found"
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
		response.Message = "Failed Delete Data"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var response model.TransactionResponse
	err := r.ParseForm()

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	var transaction model.Transaction

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
	var response model.ErrorResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	vars := mux.Vars(r)
	transactionId := vars["id"]

	data, err := db.Query(`SELECT * FROM transactions WHERE id = ?;`, transactionId)

	if data == nil {
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", transactionId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var transaction model.Transaction
	transaction.UserID, _ = strconv.Atoi(r.Form.Get("userid"))
	transaction.ProductId, _ = strconv.Atoi(r.Form.Get("productid"))
	transaction.Quantity, _ = strconv.Atoi(r.Form.Get("qty"))

	var transactionOld model.Transaction

	if err := data.Scan(&transactionOld.ID, &transactionOld.UserID, &transactionOld.ProductId, &transactionOld.Quantity); err != nil {
		log.Println(err.Error())
	}

	if transaction.UserID == 0 {
		transaction.UserID = transactionOld.UserID
	}

	if transaction.ProductId == 0 {
		transaction.ProductId = transactionOld.ProductId
	}

	if transaction.Quantity == 0 {
		transaction.Quantity = transactionOld.Quantity
	}

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

func GetDetailUserTransaction(w http.ResponseWriter, r *http.Request) {

	db := connect()
	defer db.Close()

	var transactionDetails []model.TransactionDetail

	query := "SELECT t.ID , u.ID, u.Name, u.Age, u.Address, p.ID, p.Name, p.Price, t.Quantity FROM transactions t JOIN users u ON t.UserId = u.ID JOIN products p ON p.ID = t.ProductID"

	id := r.URL.Query()["id"]
	if id != nil {
		query += " WHERE u.id = " + id[0]
	}

	rows, err := db.Query(query)

	var response model.TransactionDetailsResponse

	if err != nil {
		response.Status = 400
		response.Message = err.Error()
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var transactionDetail model.TransactionDetail
	var user model.User
	var product model.Product

	for rows.Next() {
		if err := rows.Scan(&transactionDetail.ID, &user.ID, &user.Name, &user.Age, &user.Address, &product.ID, &product.Name, &product.Price, &transactionDetail.Quantity); err != nil {
			log.Println(err.Error())
		} else {
			transactionDetail.User = user
			transactionDetail.Product = product
			transactionDetails = append(transactionDetails, transactionDetail)
		}
	}

	if len(transactionDetails) != 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = transactionDetails
	} else {
		response.Status = 400
		response.Message = "Error"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
