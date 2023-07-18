package apii

import (
	"database/sql"
	"net/http"
	"simple_shop/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createUserParams struct {
	Username       string `json:"username" binding:"required"`
	FullName       string `json:"full_name" binding:"required"`
	HashedPassword string `json:"hashed_password" binding:"required"`
	Email          string `json:"email" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := server.store.CreateUser(ctx, sqlc.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		HashedPassword: req.HashedPassword,
		Email:          req.Email,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type getUserParams struct {
	Username string `uri:"username" binding:"required"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserParams
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type listUsersParams struct {
	PageId   uint16 `form:"page_id" binding:"required,min=1"`
	PageSize uint16 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req listUsersParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	users, err := server.store.ListUsers(ctx,
		sqlc.ListUsersParams{
			Limit:  int32(req.PageSize),
			Offset: (int32(req.PageId) - 1) * int32(req.PageSize),
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}
