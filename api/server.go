package api

import (
	"fmt"

	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/DarkHeros09/e-shop/v2/token"
	"github.com/DarkHeros09/e-shop/v2/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our eshop service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// TODO: implement symmetrickey in .env file
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	userRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker, false))
	adminRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker, true))

	userRoutes.GET("/users/:id", server.getUser)        //* Finished With tests (token and changed response... No Etag)
	adminRoutes.GET("/users", server.listUsers)         //! Admin Only # Finished With tests (token and changed response... No Etag)
	userRoutes.PUT("/users/:id", server.updateUser)     //* Finished With tests (token and changed response... No Etag)
	adminRoutes.DELETE("/users/:id", server.deleteUser) //! Admin Only # Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/useraddresses", server.createUserAddress)                      //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/useraddresses/:id", server.getUserAddress)                      //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/useraddressesByUserID/:user_id", server.getUserAddressByUserID) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/useraddresses", server.listUserAddresses)                       //* Finished With tests (token and changed response... No Etag)
	userRoutes.PUT("/useraddresses/:user_id", server.updateUserAddressByUserID)      //* Finished With tests (token and changed response... No Etag)
	userRoutes.DELETE("/useraddresses/:id", server.deleteUserAddress)                //* Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/userpayments", server.createUserPayment)                      //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/userpayments/:id", server.getUserPayment)                      //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/userpaymentsByUserID/:user_id", server.getUserPaymentByUserID) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/userpayments", server.listUserPayments)                        //* Finished With tests (token and changed response... No Etag)
	userRoutes.PUT("/userpayments/:user_id", server.updateUserPaymentByUserID)      //* Finished With tests (token and changed response... No Etag)
	userRoutes.DELETE("/userpayments/:id", server.deleteUserPayment)                //* Finished With tests (token and changed response... No Etag)

	adminRoutes.POST("/categories", server.createCategory)       //! Admin Only # Finished With tests (token and changed response... No Etag)
	router.GET("/categories/:id", server.getCategory)            //? no auth required # Finished With tests (token and changed response... No Etag)
	router.GET("/categories", server.listCategories)             //? no auth required # Finished With tests (token and changed response... No Etag)
	adminRoutes.PUT("/categories/:id", server.updateCategory)    //! Admin Only # Finished With tests (token and changed response... No Etag)
	adminRoutes.DELETE("/categories/:id", server.deleteCategory) //! Admin Only # Finished With tests (token and changed response... No Etag)

	adminRoutes.POST("/inventories", server.createInventory)       //! Admin Only # Finished With tests (token and changed response... No Etag)
	router.GET("/inventories/:id", server.getInventory)            //? no auth required # Finished With tests (token and changed response... No Etag)
	router.GET("/inventories", server.listInventories)             //? no auth required # Finished With tests (token and changed response... No Etag)
	adminRoutes.PUT("/inventories/:id", server.updateInventory)    //! Admin Only # Finished With tests (token and changed response... No Etag)
	adminRoutes.DELETE("/inventories/:id", server.deleteInventory) //! Admin Only # Finished With tests (token and changed response... No Etag)

	adminRoutes.POST("/discounts", server.createDiscount)       //! Admin Only # Finished With tests (token and changed response... No Etag)
	router.GET("/discounts/:id", server.getDiscount)            //? no auth required # Finished With tests (token and changed response... No Etag)
	router.GET("/discounts", server.listDiscount)               //? no auth required # Finished With tests (token and changed response... No Etag)
	adminRoutes.PUT("/discounts/:id", server.updateDiscount)    //! Admin Only # Finished With tests (token and changed response... No Etag)
	adminRoutes.DELETE("/discounts/:id", server.deleteDiscount) //! Admin Only # Finished With tests (token and changed response... No Etag)

	adminRoutes.POST("/products", server.createProduct)       //! Admin Only # Finished With tests (token and changed response... No Etag)
	router.GET("/products/:id", server.getProduct)            //? no auth required # Finished With tests (token and changed response... No Etag)
	router.GET("/products", server.listProducts)              //? no auth required # Finished With tests (token and changed response.)
	adminRoutes.PUT("/products/:id", server.updateProduct)    //! Admin Only # Finished With tests (token and changed response... No Etag)
	adminRoutes.DELETE("/products/:id", server.deleteProduct) //! Admin Only # Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/shoppingsessions", server.createShoppingSession) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/shoppingsessions/:id", server.getShoppingSession) //* Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/cartitems", server.createCartItem)                    //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/cartitems/:session_id", server.getCartItemBySessionID) //* Finished With tests (token and changed response... No Etag)

	userRoutes.PUT("/cartitems/:session_id", server.updateCartItemBySessionID) //* Finished With tests (token and changed response... No Etag)
	userRoutes.DELETE("/cartitems/:id", server.deleteCartItemBySessionID)      //* Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/orderitems", server.createOrderItem) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/orderitems/:id", server.getOrderItem) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/orderitems", server.listOrderItems)   //* Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/orderdetails", server.createOrderDetail) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/orderdetails/:id", server.getOrderDetail) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/orderdetails", server.listOrderDetails)   //* Finished With tests (token and changed response... No Etag)

	// userRoutes.POST("/paymentdetails", server.createPaymentDetail)
	userRoutes.GET("/paymentdetails/:id", server.getPaymentDetail)
	userRoutes.GET("/paymentdetails", server.listPaymentDetails)

	server.router = router
}

// TODO: write the default tests for all the methods

// TODO: seperate superadmin from normal admin

// TODO: add update and delete methods

// TODO: modify the verification tags in the apis line required tags

// TODO: add etag logic with tests

// TODO: set header to application/json

// TODO: add caching logic with tests

// TODO: modify the list methods where needed like the listshoppingsession method. video 22 mintue 19.50

// TODO: make etags for put and get, list methods

// TODO: add gracefull shutdown logic
// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
