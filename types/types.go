package structTypes

import (
	"net/http"
	"time"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

type Storage interface {
	GetData()
	CreateUser(*UserAccount) error
	GetAllProducts() ([]Product, error)
	GetProductByID(int) (Product, error)
	GetCartByID(int) ([]CartProduct, error)
	AddToCart(int, int, int) error
	EmptyCart(int) error
	DeleteFromCart(int, int) error
	GetAllOrdersByUserID(int) ([]OrderResponse, error)
	GetOrderByID(int) (OrderResponse, error)
	CreateOrder(int, []OrderRequest) error
	UpdateOrderStatus(int, string) error
	DeleteOrder(int) error
}

type ErrorMSG struct {
	Error string `json:"error"`
}

type UserAccount struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password_hash string    `json:"password_hash"`
	Created_at    time.Time `json:"created_at"`
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Created_at  time.Time `json:"created_at"`
}

type CartProduct struct {
	CartItemID         int     `json:"cart_item_id"`
	ProductID          int     `json:"product_id"`
	ProductName        string  `json:"product_name"`
	ProductDescription string  `json:"product_description"`
	Quantity           int     `json:"quantity"`
	Price              float64 `json:"price_at_time"`
	TotalPrice         float64 `json:"total_price"`
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderRequest struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderResponse struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ProductID   int       `json:"product_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
