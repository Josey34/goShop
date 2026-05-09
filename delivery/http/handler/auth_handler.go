package handler

import (
	"context"
	"net/http"

	"github.com/Josey34/goshop/delivery/http/dto/mapper"
	"github.com/Josey34/goshop/delivery/http/dto/request"
	"github.com/Josey34/goshop/service"
	ucauth "github.com/Josey34/goshop/usecase/auth"
	"github.com/gin-gonic/gin"
)

type authService interface {
	Register(ctx context.Context, input ucauth.RegisterInput) error
	Login(ctx context.Context, input ucauth.LoginInput) (*service.LoginOutput, error)
}

type AuthHandler struct {
	service authService
}

func NewAuthHandler(svc authService) *AuthHandler {
	return &AuthHandler{service: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Register(c.Request.Context(), mapper.ToRegisterInput(req)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	out, err := h.service.Login(c.Request.Context(), mapper.ToLoginInput(req))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, mapper.ToAuthResponse(out))
}
