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

	userRoutes.GET("/users/:id", server.getUser)    //* Finished With tests (token and changed response... No Etag)
	adminRoutes.GET("/users", server.listUsers)     //! Admin Only
	userRoutes.PUT("/users/:id", server.updateUser) //* Finished With tests (token and changed response... No Etag)
	router.DELETE("/users/:id", server.deleteUser)  //! Admin Only

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

	router.POST("/categories", server.createCategory) //! Admin Only
	router.GET("/categories/:id", server.getCategory)
	router.GET("/categories", server.listCategories)

	router.POST("/inventories", server.createInventory) //! Admin Only
	router.GET("/inventories/:id", server.getInventory)
	router.GET("/inventories", server.listInventories)

	router.POST("/discounts", server.createDiscount) //! Admin Only
	router.GET("/discounts/:id", server.getDiscount)
	router.GET("/discounts", server.listDiscount)

	router.POST("/products", server.createProduct) //! Admin Only
	router.GET("/products/:id", server.getProduct)
	router.GET("/products", server.listProducts)
	router.PUT("/products/:id", server.updateProduct)    //! Admin Only
	router.DELETE("/products/:id", server.deleteProduct) //! Admin Only

	userRoutes.POST("/shoppingsessions", server.createShoppingSession)
	userRoutes.GET("/shoppingsessions/:id", server.getShoppingSession)

	userRoutes.POST("/cartitems", server.createCartItem)
	router.GET("/cartitems/:id", server.getCartItem)
	router.GET("/cartitems", server.listCartItems)

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
}

// TODO: write the default tests for all the methods

// TODO: add update and delete methods

// TODO: modify the verification tags in the apis line required tags

// TODO: add etag logic with tests

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
