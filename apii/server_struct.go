package apii

import (
	"simple_shop/db/sqlc"
	"simple_shop/db/util"
	"simple_shop/token"

	"github.com/gin-gonic/gin"
)

// Server serves HTTP request for our shopping service
type Server struct {
	config util.Config

	// * it allows us to interact with database when processing API requests from clients
	store sqlc.Store

	tokenMaker token.Maker

	// * it help us send each API request to the correct handler for processing
	router *gin.Engine
}
