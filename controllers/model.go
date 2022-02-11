package controllers

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Transaction struct {
	ID        int `json:"id"`
	UserID    int `json:"userid"`
	ProductId int `json:"productid"`
	Quantity  int `json:"qty"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}

type ProductsResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Product `json:"data"`
}

type ProductResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    Product `json:"data"`
}

type TransactionsResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []Transaction `json:"data"`
}

type TransactionResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    Transaction `json:"data"`
}
