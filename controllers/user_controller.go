package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var response UsersResponse

	query := "SELECT * FROM users"
	id := r.URL.Query()["id"]
	if id != nil {
		query += " WHERE id = " + id[0]
	}

	rows, err := db.Query(query)

	if err != nil {
		response.Status = 400
		response.Message = "Data Not Found"
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var user User
	var users []User

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address); err != nil {
			log.Println(err.Error())
		} else {
			users = append(users, user)
		}
	}

	if len(users) != 0 {
		response.Status = 200
		response.Message = "Success"
		response.Data = users
	} else {
		response.Status = 400
		response.Message = "Error Array Size Not Correct"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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
	userId := vars["id"]
	data, _ := db.Query(`SELECT * FROM users WHERE id = ?;`, userId)

	if data == nil {
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", userId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	_, errQuery := db.Query(`DELETE FROM users WHERE id = ?;`, userId)

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

func InsertUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response UserResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	var user User

	user.Name = r.Form.Get("name")
	user.Age, _ = strconv.Atoi(r.Form.Get("age"))
	user.Address = r.Form.Get("address")

	_, errQuery := db.Exec("INSERT INTO users (name, age, address) VALUES (?,?,?)", user.Name, user.Age, user.Address)

	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
		response.Data = user
	} else {
		response.Status = 400
		response.Message = "Error Insert Data"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	userId := vars["id"]

	data, _ := db.Query(`SELECT * FROM users WHERE id = ?;`, userId)

	if data == nil {
		response.Status = 400
		response.Message = fmt.Sprintf("Data using id %s not found", userId)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var user User
	user.Name = r.Form.Get("name")
	user.Age, _ = strconv.Atoi(r.Form.Get("age"))
	user.Address = r.Form.Get("address")

	_, errQuery := db.Query(`UPDATE users SET name = ?, age = ?, address = ? WHERE id = ?;`, user.Name, user.Age, user.Address, userId)

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
