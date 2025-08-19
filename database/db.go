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
	query := `INSERT INTO users(username, email, password_hash, created_at) VALUES($1, $2, $3, $4) RETURNING id;`
	var userId int
	err := s.DB.QueryRow(query, user.Username, user.Email, user.Password_hash, user.Created_at).Scan(&userId)
	if err != nil {
		return err
	}
	query2 := `INSERT INTO carts (user_id) VALUES ($1) RETURNING id;`
	_, err = s.DB.Exec(query2, userId)
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

func (s *PostgresStore) AddToCart(userID, productID, quantity int) error {
	var cartID int
	err := s.DB.QueryRow(`SELECT id FROM carts WHERE user_id = $1`, userID).Scan(&cartID)
	if err != nil {
		return fmt.Errorf("cart not found for user %d: %w", userID, err)
	}
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity;
	`
	_, err = s.DB.Exec(query, cartID, productID, quantity)
	return err
}

func (s *PostgresStore) EmptyCart(userID int) error {
	_, err := s.DB.Exec(`
        DELETE FROM cart_items
        WHERE cart_id = (SELECT id FROM carts WHERE user_id = $1)
    `, userID)
	return err
}

func (s *PostgresStore) DeleteFromCart(userID, productID int) error {
	_, err := s.DB.Exec(`
        DELETE FROM cart_items
        WHERE cart_id = (SELECT id FROM carts WHERE user_id = $1)
        AND product_id = $2
    `, userID, productID)
	return err
}

// ORDER FUNCTIONS

func (s *PostgresStore) CreateOrder(userID int, orders []structTypes.OrderRequest) error {
	insertQuery := `INSERT INTO orders (user_id, product_id, quantity, price)
                    VALUES ($1, $2, $3, $4);`

	updateQuery := `UPDATE products
                    SET stock = stock - $1
                    WHERE id = $2 AND stock >= $1;`

	getStockQuery := `SELECT stock FROM products WHERE id = $1;`

	for _, order := range orders {
		var stock int
		err := s.DB.QueryRow(getStockQuery, order.ProductID).Scan(&stock)
		if err != nil {
			return err
		}

		if stock <= 0 {
			return fmt.Errorf("product_id %d is out of stock", order.ProductID)
		}
		if stock < order.Quantity {
			return fmt.Errorf("not enough stock for product_id %d (available: %d, requested: %d)",
				order.ProductID, stock, order.Quantity)
		}

		res, err := s.DB.Exec(updateQuery, order.Quantity, order.ProductID)
		if err != nil {
			return err
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return fmt.Errorf("not enough stock for product_id %d", order.ProductID)
		}

		_, err = s.DB.Exec(insertQuery, userID, order.ProductID, order.Quantity, order.Price)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStore) GetAllOrdersByUserID(userID int) ([]structTypes.OrderResponse, error) {
	var orders []structTypes.OrderResponse

	query := `
		SELECT 
			o.id, o.user_id, o.product_id, 
			p.name, p.description,
			o.quantity, o.price, o.status, o.created_at
		FROM orders o
		JOIN products p ON o.product_id = p.id
		WHERE o.user_id = $1
	`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order structTypes.OrderResponse
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.ProductID,
			&order.ProductName,
			&order.Description,
			&order.Quantity,
			&order.Price,
			&order.Status,
			&order.CreatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *PostgresStore) GetOrderByID(orderID int) (structTypes.OrderResponse, error) {
	var order structTypes.OrderResponse

	query := `
		SELECT 
			o.id, o.user_id, o.product_id, 
			p.name, p.description,
			o.quantity, o.price, o.status, o.created_at
		FROM orders o
		JOIN products p ON o.product_id = p.id
		WHERE o.id = $1
	`
	row := s.DB.QueryRow(query, orderID)

	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.ProductID,
		&order.ProductName,
		&order.Description,
		&order.Quantity,
		&order.Price,
		&order.Status,
		&order.CreatedAt,
	)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (s *PostgresStore) UpdateOrderStatus(orderID int, status string) error {
	query := `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`
	_, err := s.DB.Exec(query, status, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteOrder(orderID int) error {
	// Step 1: get product_id and quantity of the order
	var productID, quantity int
	queryGet := `SELECT product_id, quantity FROM orders WHERE id = $1`
	err := s.DB.QueryRow(queryGet, orderID).Scan(&productID, &quantity)
	if err != nil {
		return err
	}

	updateProductQuery := `UPDATE products SET stock = stock + $1 WHERE id = $2`
	_, err = s.DB.Exec(updateProductQuery, quantity, productID)
	if err != nil {
		return err
	}

	deleteQuery := `DELETE FROM orders WHERE id = $1`
	_, err = s.DB.Exec(deleteQuery, orderID)
	if err != nil {
		return err
	}

	return nil
}

// REVIEWS FUNCTIONS

func (s *PostgresStore) CreateNewReview(review structTypes.ReviewRequest) error {
	query := "insert into reviews (user_id, product_id, rating, comment) values ($1, $2, $3, $4);"
	_, err := s.DB.Exec(query, review.UserID, review.ProductID, review.Rating, review.Comment)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetAllReviewsByProductID(productID int) ([]structTypes.ReviewResponse, error) {
	var reviews []structTypes.ReviewResponse
	query := `SELECT 
			p.id AS product_id,
			p.name AS product_name,
			p.description AS product_description,
			p.price,
			p.stock,
			p.created_at AS product_created_at,
			
			r.id AS review_id,
			r.rating,
			r.comment,
			r.created_at AS review_created_at,
			
			u.id AS user_id,
			u.username,
			u.email
		FROM products p
		LEFT JOIN reviews r ON p.id = r.product_id
		LEFT JOIN users u ON r.user_id = u.id
		WHERE p.id = $1;
		`
	data, err := s.DB.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	for data.Next() {
		var review structTypes.ReviewResponse
		err := data.Scan(
			&review.Product.ID,
			&review.Product.Name,
			&review.Product.Description,
			&review.Product.Price,
			&review.Product.Stock,
			&review.Product.Created_at,

			&review.ID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,

			&review.User.ID,
			&review.User.Username,
			&review.User.Email,
		)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}
	return reviews, nil
}
