// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Active    bool      `json:"active"`
	TypeID    int64     `json:"type_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`
}

type AdminType struct {
	ID        int64     `json:"id"`
	AdminType string    `json:"admin_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartItem struct {
	ID        int64 `json:"id"`
	SessionID int64 `json:"session_id"`
	ProductID int64 `json:"product_id"`
	// must be positive
	Quantity  int32     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Discount struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// default is 0
	DiscountPercent string `json:"discount_percent"`
	// default is false
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderDetail struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	// must be positive
	Total     string        `json:"total"`
	PaymentID sql.NullInt64 `json:"payment_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type OrderItem struct {
	ID        int64 `json:"id"`
	OrderID   int64 `json:"order_id"`
	ProductID int64 `json:"product_id"`
	// must be positive
	Quantity  int32     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaymentDetail struct {
	ID int64 `json:"id"`
	// default is 0
	OrderID int64 `json:"order_id"`
	// must be positive
	Amount    int32     `json:"amount"`
	Provider  string    `json:"provider"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sku         string `json:"sku"`
	CategoryID  int64  `json:"category_id"`
	InventoryID int64  `json:"inventory_id"`
	// must be positive
	Price string `json:"price"`
	// default is false
	Active     bool      `json:"active"`
	DiscountID int64     `json:"discount_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProductCategory struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// default is false
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductInventory struct {
	ID int64 `json:"id"`
	// must be positive
	Quantity int32 `json:"quantity"`
	// default is true
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShoppingSession struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
	// must be positive
	Total     string    `json:"total"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Telephone int32     `json:"telephone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserAddress struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	Telephone   int32  `json:"telephone"`
}

type UserPayment struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	PaymentType string    `json:"payment_type"`
	Provider    string    `json:"provider"`
	AccountNo   int32     `json:"account_no"`
	Expiry      time.Time `json:"expiry"`
}

type UserSession struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}
