package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/skamranahmed/banking-system/db/sqlc"
	"github.com/skamranahmed/banking-system/utils"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func (server *Server) createUser(c *gin.Context) {
	var req createUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	plainTextPassword := req.Password
	hashedPassword, err := utils.HashPassword(plainTextPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Password: hashedPassword,
		FullName: req.FullName,
		Email:    req.Email,
	}

	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := &createUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}

	c.JSON(http.StatusCreated, response)
	return
}
