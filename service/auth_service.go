package service

import (
	"context"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/pkg/jwt"
	ucauth "github.com/Josey34/goshop/usecase/auth"
)

type LoginOutput struct {
	Token    string
	Customer *entity.Customer
}

type AuthService struct {
	registerUC *ucauth.RegisiterUseCase
	loginUC    *ucauth.LoginUseCase
	jwtService *jwt.JWTService
}

func NewAuthService(
	registerUC *ucauth.RegisiterUseCase,
	loginUC *ucauth.LoginUseCase,
	jwtService *jwt.JWTService,
) *AuthService {
	return &AuthService{registerUC, loginUC, jwtService}
}

func (s *AuthService) Register(ctx context.Context, input ucauth.RegisterInput) error {
	return s.registerUC.Execute(ctx, input)
}

func (s *AuthService) Login(ctx context.Context, input ucauth.LoginInput) (*LoginOutput, error) {
	customer, err := s.loginUC.Execute(ctx, input)
	if err != nil {
		return nil, err
	}

	token, err := s.jwtService.Generate(customer.ID, string(customer.Email.Value()))
	if err != nil {
		return nil, err
	}

	return &LoginOutput{Token: token, Customer: customer}, nil
}
