package api

import (
	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our eshop service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.GET("/users", server.listUsers)

	router.POST("/useraddresses", server.createUserAddress)
	router.GET("/useraddresses/:id", server.getUserAddress)
	router.GET("/useraddresses", server.listUserAddresses)

	router.POST("/userpayments", server.createUserPayment)
	router.GET("/userpayments/:id", server.getUserPayment)
	router.GET("/userpayments", server.listUserPayments)

	router.POST("/categories", server.createCategory)
	router.GET("/categories/:id", server.getCategory)
	router.GET("/categories", server.listCategories)

	router.GET("/inventories/:id", server.getInventory)

	router.GET("/discounts/:id", server.getDiscount)

	router.POST("/products", server.createProduct)
	router.GET("/products/:id", server.getProduct)
	router.GET("/products", server.listProducts)

	router.POST("/shoppingsessions", server.createShoppingSession)
	router.GET("/shoppingsessions/:id", server.getShoppingSession)

	router.POST("/orderItems", server.createOrderItem)
	router.GET("/orderItemsByID/:id", server.getOrderItemByID)
	router.GET("/orderItemsByOrderDetailID/:id", server.getOrderItemByOrderDetailID)
	router.GET("/orderItems", server.listOrderItems)

	router.POST("/orderDetails", server.createOrderDetail)
	router.GET("/orderDetails/:id", server.getOrderDetail)
	router.GET("/orderDetails", server.listOrderDetails)

	router.POST("/paymentDetails", server.createPaymentDetail)
	router.GET("/paymentDetails/:id", server.getPaymentDetail)
	router.GET("/paymentDetails", server.listPaymentDetails)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address
// TODO: add gracefull shutdown logic
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
