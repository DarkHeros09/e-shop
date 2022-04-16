package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// server.gracefullShutDown(server.router)
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	userRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker, false))
	adminRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker, true))

	userRoutes.GET("/users/:id", server.getUser)        //* Finished With tests (token and changed response... No Etag)
	adminRoutes.GET("/users", server.listUsers)         //! Admin Only # Finished With tests (token and changed response... No Etag)
	userRoutes.PUT("/users/:id", server.updateUser)     //* Finished With tests (token and changed response... No Etag)
	adminRoutes.DELETE("/users/:id", server.deleteUser) //! Admin Only # Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/users/addresses", server.createUserAddress)                 //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/users/addresses/:id", server.getUserAddress)                 //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/users/addresses", server.listUserAddresses)                  //* Finished With tests (token and changed response... No Etag)
	userRoutes.PUT("/users/addresses/:user_id", server.updateUserAddressByUserID) //* Finished With tests (token and changed response... No Etag)
	userRoutes.DELETE("/users/addresses/:id", server.deleteUserAddress)           //* Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/users/payments", server.createUserPayment)         //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/users/payments/:id", server.getUserPayment)         //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/users/payments", server.listUserPayments)           //* Finished With tests (token and changed response... No Etag)
	userRoutes.PUT("/users/payments/:user_id", server.updateUserPayment) //* Finished With tests (token and changed response... No Etag)
	userRoutes.DELETE("/users/payments/:id", server.deleteUserPayment)   //* Finished With tests (token and changed response... No Etag)

	adminRoutes.POST("/products/categories", server.createCategory)       //! Admin Only # Finished With tests (token and changed response... No Etag)
	router.GET("/products/categories/:id", server.getCategory)            //? no auth required # Finished With tests (token and changed response... No Etag)
	router.GET("/products/categories", server.listCategories)             //? no auth required # Finished With tests (token and changed response... No Etag)
	adminRoutes.PUT("/products/categories/:id", server.updateCategory)    //! Admin Only # Finished With tests (token and changed response... No Etag)
	adminRoutes.DELETE("/products/categories/:id", server.deleteCategory) //! Admin Only # Finished With tests (token and changed response... No Etag)

	adminRoutes.POST("/products/inventories", server.createInventory)       //! Admin Only # Finished With tests (token and changed response... No Etag)
	router.GET("/products/inventories/:id", server.getInventory)            //? no auth required # Finished With tests (token and changed response... No Etag)
	router.GET("/products/inventories", server.listInventories)             //? no auth required # Finished With tests (token and changed response... No Etag)
	adminRoutes.PUT("/products/inventories/:id", server.updateInventory)    //! Admin Only # Finished With tests (token and changed response... No Etag)
	adminRoutes.DELETE("/products/inventories/:id", server.deleteInventory) //! Admin Only # Finished With tests (token and changed response... No Etag)

	adminRoutes.POST("/products/discounts", server.createDiscount)       //! Admin Only # Finished With tests (token and changed response... No Etag)
	router.GET("/products/discounts/:id", server.getDiscount)            //? no auth required # Finished With tests (token and changed response... No Etag)
	router.GET("/products/discounts", server.listDiscount)               //? no auth required # Finished With tests (token and changed response... No Etag)
	adminRoutes.PUT("/products/discounts/:id", server.updateDiscount)    //! Admin Only # Finished With tests (token and changed response... No Etag)
	adminRoutes.DELETE("/products/discounts/:id", server.deleteDiscount) //! Admin Only # Finished With tests (token and changed response... No Etag)

	adminRoutes.POST("/products", server.createProduct)       //! Admin Only # Finished With tests (token and changed response... No Etag)
	router.GET("/products/:id", server.getProduct)            //? no auth required # Finished With tests (token and changed response... No Etag)
	router.GET("/products", server.listProducts)              //? no auth required # Finished With tests (token and changed response.)
	adminRoutes.PUT("/products/:id", server.updateProduct)    //! Admin Only # Finished With tests (token and changed response... No Etag)
	adminRoutes.DELETE("/products/:id", server.deleteProduct) //! Admin Only # Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/shopping-sessions", server.createShoppingSession) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/shopping-sessions/:id", server.getShoppingSession) //* Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/cart-items", server.createCartItem)                       //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/cart-items/:session_id", server.getCartItemBySessionID)    //* Finished With tests (token and changed response... No Etag)
	userRoutes.PUT("/cart-items/:session_id", server.updateCartItemBySessionID) //* Finished With tests (token and changed response... No Etag)
	userRoutes.DELETE("/cart-items/:id", server.deleteCartItemBySessionID)      //* Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/order-items", server.createOrderItem) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/order-items/:id", server.getOrderItem) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/order-items", server.listOrderItems)   //* Finished With tests (token and changed response... No Etag)

	userRoutes.POST("/order-details", server.createOrderDetail) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/order-details/:id", server.getOrderDetail) //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/order-details", server.listOrderDetails)   //* Finished With tests (token and changed response... No Etag)

	userRoutes.GET("/payment-details/:id", server.getPaymentDetail)    //* Finished With tests (token and changed response... No Etag)
	userRoutes.GET("/payment-details", server.listPaymentDetails)      //* Finished With tests (token and changed response... No Etag)
	userRoutes.PUT("/payment-details/:id", server.updatePaymentDetail) //* Finished With tests (token and changed response... No Etag)

	server.router = router

}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) gracefullShutDown(router *gin.Engine) {
	srv := &http.Server{
		Addr:    server.config.ServerAddress,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

// DONE: write the default tests for all the methods

// TODO: seperate superadmin from normal admin

// DONE: add update and delete methods

// DONE: modify the json verification tags in the apis line required tags

// TODO: add etag logic with tests

// TODO: add caching logic with tests, try groupcache

// TODO: add refresh token

// DONE: modify the list methods where needed like the listshoppingsession method. video 22 mintue 19.50

// TODO: make etags for put and get, list methods

// DONE: use heraricy in api call

// DONE: use "-" to seprrate

// DONE: add gracefull shutdown logic
