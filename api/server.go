package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/skamranahmed/banking-system/db/sqlc"
)

// Server : will serve the HTTP requests for our API
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer : will create a new Server and also setup the routes
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	// gin router
	router := gin.Default()

	// get the binding engine that gin is using
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// setup routes
	router.GET("/accounts", server.listAccounts)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start runs the HTTP server on the provided port
func (server *Server) Start(port string) error {
	address := fmt.Sprintf(":%s", port)
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
