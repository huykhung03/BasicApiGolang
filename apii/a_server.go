package apii

import (
	"simple_shop/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP request for our shopping service
type Server struct {
	// * it allows us to interact with database when processing API requests from clients
	store sqlc.Store
	// * it help us send each API request to the correct handler for processing
	router *gin.Engine
}

func NewServer(store sqlc.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	// * add routes to router
	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUser)
	router.GET("/users", server.listUsers)

	server.router = router

	return server
}

// Start rusn the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{
		"err": err.Error(),
	}
}
