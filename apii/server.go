package apii

import (
	"fmt"
	"simple_shop/db/sqlc"
	"simple_shop/db/util"
	"simple_shop/token"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func (server *Server) setUpServer() {
	router := gin.Default()

	// * add routes to router
	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUser)
	router.GET("/users", server.listUsers)
	router.POST("/users/login", server.loginUser)

	router.POST("/purchases", server.createPurchase)

	server.router = router
}

func NewServer(config util.Config, store sqlc.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, oke := binding.Validator.Engine().(*validator.Validate); oke {
		v.RegisterValidation("currency", ValidCurrency)
	}

	server.setUpServer()

	return server, nil
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
