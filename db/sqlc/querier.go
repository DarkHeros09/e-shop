// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
)

type Querier interface {
	CreateAdmin(ctx context.Context, arg CreateAdminParams) (Admin, error)
	CreateAdminType(ctx context.Context, adminType string) (AdminType, error)
	CreateCartItem(ctx context.Context, arg CreateCartItemParams) (CartItem, error)
	CreateDiscount(ctx context.Context, arg CreateDiscountParams) (Discount, error)
	CreateOrderDetailAndPaymentDetail(ctx context.Context, arg CreateOrderDetailAndPaymentDetailParams) (OrderDetail, error)
	CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error)
	CreatePaymentDetail(ctx context.Context, arg CreatePaymentDetailParams) (PaymentDetail, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateProductCategory(ctx context.Context, arg CreateProductCategoryParams) (ProductCategory, error)
	CreateProductInventory(ctx context.Context, quantity int32) (ProductInventory, error)
	CreateShoppingSession(ctx context.Context, arg CreateShoppingSessionParams) (ShoppingSession, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserAddress(ctx context.Context, arg CreateUserAddressParams) (UserAddress, error)
	CreateUserPayment(ctx context.Context, arg CreateUserPaymentParams) (UserPayment, error)
	DeleteAdmin(ctx context.Context, id int64) error
	DeleteAdminTypeByID(ctx context.Context, id int64) error
	DeleteAdminTypeByType(ctx context.Context, adminType string) error
	DeleteCartItem(ctx context.Context, id int64) error
	DeleteDiscount(ctx context.Context, id int64) error
	DeleteOrderDetail(ctx context.Context, id int64) error
	DeleteOrderItem(ctx context.Context, id int64) error
	DeletePaymentDetail(ctx context.Context, id int64) error
	DeleteProduct(ctx context.Context, id int64) error
	DeleteProductCategory(ctx context.Context, id int64) error
	DeleteProductInventory(ctx context.Context, id int64) error
	DeleteShoppingSession(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	DeleteUserAddress(ctx context.Context, id int64) error
	DeleteUserPayment(ctx context.Context, arg DeleteUserPaymentParams) error
	GetAdmin(ctx context.Context, id int64) (Admin, error)
	GetAdminByEmail(ctx context.Context, email string) (Admin, error)
	GetAdminType(ctx context.Context, id int64) (AdminType, error)
	GetCartItemByID(ctx context.Context, id int64) (CartItem, error)
	GetCartItemBySessionID(ctx context.Context, sessionID int64) (CartItem, error)
	GetDiscount(ctx context.Context, id int64) (Discount, error)
	GetOrderDetail(ctx context.Context, id int64) (OrderDetail, error)
	GetOrderItem(ctx context.Context, arg GetOrderItemParams) (OrderItem, error)
	GetPaymentDetail(ctx context.Context, arg GetPaymentDetailParams) (PaymentDetail, error)
	GetProduct(ctx context.Context, id int64) (Product, error)
	GetProductCategory(ctx context.Context, id int64) (ProductCategory, error)
	GetProductInventory(ctx context.Context, id int64) (ProductInventory, error)
	GetProductInventoryForUpdate(ctx context.Context, id int64) (ProductInventory, error)
	GetShoppingSession(ctx context.Context, id int64) (ShoppingSession, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserAddress(ctx context.Context, arg GetUserAddressParams) (UserAddress, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserPayment(ctx context.Context, arg GetUserPaymentParams) (UserPayment, error)
	ListAdminTypes(ctx context.Context, arg ListAdminTypesParams) ([]AdminType, error)
	ListAdmins(ctx context.Context, arg ListAdminsParams) ([]Admin, error)
	ListCartItem(ctx context.Context, arg ListCartItemParams) ([]CartItem, error)
	ListDiscounts(ctx context.Context, arg ListDiscountsParams) ([]Discount, error)
	ListOrderDetails(ctx context.Context, arg ListOrderDetailsParams) ([]OrderDetail, error)
	ListOrderItems(ctx context.Context, arg ListOrderItemsParams) ([]OrderItem, error)
	ListPaymentDetails(ctx context.Context, arg ListPaymentDetailsParams) ([]ListPaymentDetailsRow, error)
	ListProductCategories(ctx context.Context, arg ListProductCategoriesParams) ([]ProductCategory, error)
	ListProductInventories(ctx context.Context, arg ListProductInventoriesParams) ([]ProductInventory, error)
	ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error)
	ListShoppingSessions(ctx context.Context, arg ListShoppingSessionsParams) ([]ShoppingSession, error)
	ListUserAddresses(ctx context.Context, arg ListUserAddressesParams) ([]UserAddress, error)
	ListUserPayments(ctx context.Context, arg ListUserPaymentsParams) ([]UserPayment, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateAdmin(ctx context.Context, arg UpdateAdminParams) (Admin, error)
	UpdateAdminType(ctx context.Context, arg UpdateAdminTypeParams) (AdminType, error)
	UpdateCartItem(ctx context.Context, arg UpdateCartItemParams) (CartItem, error)
	UpdateDiscount(ctx context.Context, arg UpdateDiscountParams) (Discount, error)
	UpdateOrderDetail(ctx context.Context, arg UpdateOrderDetailParams) (OrderDetail, error)
	UpdateOrderItem(ctx context.Context, arg UpdateOrderItemParams) (OrderItem, error)
	UpdatePaymentDetail(ctx context.Context, arg UpdatePaymentDetailParams) (PaymentDetail, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateProductCategory(ctx context.Context, arg UpdateProductCategoryParams) (ProductCategory, error)
	UpdateProductInventory(ctx context.Context, arg UpdateProductInventoryParams) (ProductInventory, error)
	UpdateProductQuantity(ctx context.Context, arg UpdateProductQuantityParams) (ProductInventory, error)
	UpdateShoppingSession(ctx context.Context, arg UpdateShoppingSessionParams) (ShoppingSession, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserAddress(ctx context.Context, arg UpdateUserAddressParams) (UserAddress, error)
	UpdateUserAddressByUserID(ctx context.Context, arg UpdateUserAddressByUserIDParams) (UserAddress, error)
	UpdateUserPayment(ctx context.Context, arg UpdateUserPaymentParams) (UserPayment, error)
}

var _ Querier = (*Queries)(nil)
