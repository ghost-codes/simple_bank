package api

import (
	"net/http"
	"time"

	db "github.com/ghost-codes/simplebank/db/sqlc"
	"github.com/ghost-codes/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(context *gin.Context) {
	var req CreateUserRequest
	err := context.BindJSON(&req)
	if err != nil {

		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hasedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.CreateuserParams{
		Username:       req.Username,
		HashedPassword: hasedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.Createuser(context, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "uniue_violation":
				context.JSON(http.StatusForbidden, errorResponse(err))
				return
			}

		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := UserResponse{
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	context.JSON(http.StatusOK, res)
}
