package database

import (
	"database/sql"
	"fmt"

	structTypes "github.com/VincentSamuelPaul/production-api/types"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "host=localhost port=5433 user=admin dbname=postgres password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		DB: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	query := `CREATE TABLE if not exists users (
			id serial PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP
		);`
	_, err := s.DB.Exec(query)
	if err != nil {
		return err
	}
	query = `create table if not exists products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		price NUMERIC(10,2) NOT NULL,
		stock INT NOT NULL,
		created_at TIMESTAMP DEFAULT now()
		);`
	_, err = s.DB.Exec(query)
	if err != nil {
		return err
	}
	query = `create table if not exists carts (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id),
		created_at TIMESTAMP DEFAULT now()
		);`
	_, err = s.DB.Exec(query)
	if err != nil {
		return err
	}
	query = `create table if not exists cart_items (
		id SERIAL PRIMARY KEY,
		cart_id INT REFERENCES carts(id),
		product_id INT REFERENCES products(id),
		quantity INT NOT NULL
		);`
	_, err = s.DB.Exec(query)
	if err != nil {
		return err
	}
	query = `create table if not exists orders (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id),
		total NUMERIC(10,2) NOT NULL,
		status TEXT DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT now()
		);`
	_, err = s.DB.Exec(query)
	if err != nil {
		return err
	}
	query = `create table if not exists order_items (
		id SERIAL PRIMARY KEY,
		order_id INT REFERENCES orders(id),
		product_id INT REFERENCES products(id),
		quantity INT NOT NULL,
		price NUMERIC(10,2) NOT NULL
		);`
	_, err = s.DB.Exec(query)
	if err != nil {
		return err
	}
	query = `create table if not exists reviews (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id),
		product_id INT REFERENCES products(id),
		rating INT CHECK (rating >= 1 AND rating <= 5),
		comment TEXT,
		created_at TIMESTAMP DEFAULT now()
		);`
	_, err = s.DB.Exec(query)
	if err != nil {
		return err
	}
	// p := structTypes.Product{
	// 	Name:        "Computer",
	// 	Description: "macbook air 2020",
	// 	Price:       92000,
	// 	Stock:       100,
	// 	Created_at:  time.Now(),
	// }
	// query = `insert into products(name, description, price, stock, created_at) values($1, $2, $3, $4, $5)`
	// _, err = s.DB.Exec(query, p.Name, p.Description, p.Price, p.Stock, p.Created_at)
	// if err != nil {
	// 	return err
	// }
	return nil
}

// AUTH FUNCTIONS

func (s *PostgresStore) CreateUser(user *structTypes.UserAccount) error {
	query := `insert into users(username, email, password_hash, created_at) values($1, $2, $3, $4)`

	_, err := s.DB.Exec(query, user.Username, user.Email, user.Password_hash, user.Created_at)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetData() {
	query := "select * from users;"
	data, err := s.DB.Query(query)
	if err != nil {
		return
	}
	for data.Next() {
		var account structTypes.UserAccount
		data.Scan(
			&account.ID,
			&account.Username,
			&account.Email,
			&account.Password_hash,
			&account.Created_at,
		)
		fmt.Println(account)
	}
}

// PRODUCT FUNCTIONS

func (s *PostgresStore) GetAllProducts() ([]structTypes.Product, error) {
	var Products []structTypes.Product
	query := "select * from products;"
	data, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	for data.Next() {
		var product structTypes.Product
		data.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Created_at,
		)
		Products = append(Products, product)
	}
	return Products, nil
}

func (s *PostgresStore) GetProductByID(id int) (structTypes.Product, error) {
	var product structTypes.Product
	query := fmt.Sprintf("select * from products where id = %d;", id)
	data, err := s.DB.Query(query)
	if err != nil {
		return product, err
	}
	for data.Next() {
		data.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Created_at,
		)
	}
	return product, nil
}

// CART FUNCTIONS

func (s *PostgresStore) GetCartByID(id int) ([]structTypes.CartProduct, error) {
	var cartProducts []structTypes.CartProduct
	query := fmt.Sprintf(`SELECT 
    ci.id AS cart_item_id,
    p.id AS product_id,
    p.name AS product_name,
    p.description,
    ci.quantity,
    ci.price_at_time,
    (ci.quantity * ci.price_at_time) AS total_price
	FROM carts c
	JOIN cart_items ci ON ci.cart_id = c.id
	JOIN products p ON p.id = ci.product_id
	WHERE c.user_id = %d;
	`, id)
	data, err := s.DB.Query(query)
	if err != nil {
		return cartProducts, err
	}
	for data.Next() {
		var cartProduct structTypes.CartProduct
		data.Scan(
			&cartProduct.CartItemID,
			&cartProduct.ProductID,
			&cartProduct.ProductName,
			&cartProduct.ProductDescription,
			&cartProduct.Quantity,
			&cartProduct.Price,
			&cartProduct.TotalPrice,
		)
		cartProducts = append(cartProducts, cartProduct)
	}
	return cartProducts, nil
}
