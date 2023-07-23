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

	router.POST("/create-user", server.createUser)
	router.POST("/create-admin", server.createAdmin)
	router.POST("/login", server.loginUser)

	// authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	user := router.Group("/user").Use(authMiddleware(server.tokenMaker))
	{
		user.GET("/info", server.getUser)
		user.POST("/purchase", server.createPurchase)
		user.GET("/list-bank-accounts", server.listBankAccounts)
		user.GET("/list-purchase-histories", server.listPurchaseHistoresOfUser)
	}

	admin := router.Group("/admin").Use(authMiddleware(server.tokenMaker))
	{
		admin.POST("/create-product", server.createProduct)
		admin.POST("/delete-product", server.deleteProduct)
		admin.GET("/list-products", server.listProducts)
		admin.GET("/:id_product", server.getProduct)
		admin.GET("/list-purchase-histories", server.listPurchaseHistoriesOfAdmin)
	}

	bankAccount := router.Group("/bank-account").Use(authMiddleware(server.tokenMaker))
	{
		bankAccount.POST("/create-bank-account", server.createBankAccount)
		bankAccount.POST("/delete-bank-accounts", server.deleteBankAccount)
	}

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
