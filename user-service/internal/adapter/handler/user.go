package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/core/domain"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/core/port"
)

type UserHandler struct {
	svc port.UserService
}

func NewUserHandler(svc port.UserService) *UserHandler {
	return &UserHandler{
		svc,
	}
}

func (uh *UserHandler) RegisterUser(c *gin.Context) {
	var req domain.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("email, password and name are required")})
		return
	}

	user := &domain.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}
	_, err := uh.svc.RegisterUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
	})
}

func (uh *UserHandler) ConfirmEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "token is required",
		})
		return
	}

	_, err := uh.svc.ConfirmEmail(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "email confirmed",
	})
}

func (uh *UserHandler) GetUsers(c *gin.Context) {
	users, err := uh.svc.GetUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": domain.ErrNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, users)
}
