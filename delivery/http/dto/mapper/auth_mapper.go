package mapper

import (
	"github.com/Josey34/goshop/delivery/http/dto/request"
	"github.com/Josey34/goshop/delivery/http/dto/response"
	"github.com/Josey34/goshop/service"
	ucauth "github.com/Josey34/goshop/usecase/auth"
)

func ToRegisterInput(req request.RegisterRequest) ucauth.RegisterInput {
	return ucauth.RegisterInput{
		Name:       req.Name,
		Email:      req.Email,
		Password:   req.Password,
		Phone:      req.Phone,
		Street:     req.Street,
		City:       req.City,
		Province:   req.Province,
		PostalCode: req.PostalCode,
	}
}

func ToLoginInput(req request.LoginRequest) ucauth.LoginInput {
	return ucauth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToAuthResponse(sl *service.LoginOutput) response.AuthResponse {
	return response.AuthResponse{
		Token:      sl.Token,
		CustomerID: sl.Customer.ID,
		Email:      sl.Customer.Email.Value(),
	}
}
