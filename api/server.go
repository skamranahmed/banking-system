package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/skamranahmed/banking-system/config"
	db "github.com/skamranahmed/banking-system/db/sqlc"
	"github.com/skamranahmed/banking-system/token"
)

// Server : will serve the HTTP requests for our API
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer : will create a new Server and also setup the routes
func NewServer(store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSigningKey)
	if err != nil {
		return nil, fmt.Errorf("unable to initialise token maker, err: %v", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
	}

	// get the binding engine that gin is using
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	// gin router
	router := gin.Default()

	// setup routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// authenticated routes
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start runs the HTTP server on the provided port
func (server *Server) Start(port string) error {
	address := fmt.Sprintf(":%s", port)
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
