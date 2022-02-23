package controllers

import (
	model "PraktikumPBP/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var response model.UsersResponse

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

	var user model.User
	var users []model.User

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
	} else if response.Message == "" {
		response.Status = 400
		response.Message = "Data Not Found"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
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
	userId := vars["id"]
	query, errQuery := db.Exec(`DELETE FROM users WHERE id = ?;`, userId)
	RowsAffected, err := query.RowsAffected()

	if RowsAffected == 0 {
		response.Status = 400
		response.Message = "User not found"
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
		log.Println(errQuery.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response model.UserResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	var user model.User

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
		log.Println(errQuery.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	var response model.UserResponse

	if err != nil {
		response.Status = 400
		response.Message = "Error Parsing Data"
		w.WriteHeader(400)
		log.Println(err.Error())
		return
	}

	vars := mux.Vars(r)
	userId := vars["id"]

	var user model.User
	user.Name = r.Form.Get("name")
	user.Age, _ = strconv.Atoi(r.Form.Get("age"))
	user.Address = r.Form.Get("address")
	user.Password = r.Form.Get("password")

	rows, _ := db.Query("SELECT * FROM users WHERE id = ?", userId)
	var prevDatas []model.User
	var prevData model.User

	for rows.Next() {
		if err := rows.Scan(&prevData.ID, &prevData.Name, &prevData.Age, &prevData.Address, &prevData.Password); err != nil {
			log.Println(err.Error())
		} else {
			prevDatas = append(prevDatas, prevData)
		}
	}

	if len(prevDatas) > 0 {
		if user.Name == "" {
			user.Name = prevDatas[0].Name
		}
		if user.Age == 0 {
			user.Age = prevDatas[0].Age
		}
		if user.Address == "" {
			user.Address = prevDatas[0].Address
		}
		if user.Password == "" {
			user.Password = prevDatas[0].Password
		}

		_, errQuery := db.Exec(`UPDATE users SET name = ?, age = ?, address = ?, password = ? WHERE id = ?;`, user.Name, user.Age, user.Address, user.Password, userId)

		if errQuery == nil {
			response.Status = 200
			response.Message = "Success Update Data"
			id, _ := strconv.Atoi(userId)
			user.ID = id
			response.Data = user
			w.WriteHeader(200)
		} else {
			response.Status = 400
			response.Message = "Error Update Data"
			w.WriteHeader(400)
			log.Println(errQuery)
		}

	} else {
		response.Status = 400
		response.Message = "Data Not Found"
		w.WriteHeader(400)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
