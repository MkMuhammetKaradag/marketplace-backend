// internal/seller-service/transport/http/controller/seller_onboard.go
package controller

import (
	"marketplace/internal/seller-service/domain"
	"marketplace/internal/seller-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type CreateSellerRequest struct {
	StoreName         string `json:"store_name" validate:"required"`
	TaxNumber         string `json:"tax_number" validate:"required"`
	TaxOffice         string `json:"tax_office" validate:"required"`
	LegalBusinessName string `json:"legal_business_name" validate:"required"`
	PhoneNumber       string `json:"phone_number" validate:"required"`
	AddressLine       string `json:"address_line" validate:"required"`
	// AddressLine2 string `json:"address_line2" validate:"required"`
	City string `json:"city" validate:"required"`
	// State string `json:"state" validate:"required"`
	// ZipCode string `json:"zip_code" validate:"required"`
	Country               string `json:"country" validate:"required"`
	BankAccountIban       string `json:"bank_account_iban" validate:"required"`
	BankAccountHolderName string `json:"bank_account_holder_name" validate:"required"`
	BankAccountBic        string `json:"bank_account_bic" validate:"required"`
	Email                 string `json:"email" validate:"required,email"`
}

type CreateSellerResponse struct {
	Message  string `json:"message"`
	SellerId string `json:"seller_id"`
}
type CreateSellerController struct {
	usecase usecase.CreateSellerUseCase
}

func NewCreateSellerController(usecase usecase.CreateSellerUseCase) *CreateSellerController {
	return &CreateSellerController{
		usecase: usecase,
	}
}

func (h *CreateSellerController) Handle(fbrCtx *fiber.Ctx, req *CreateSellerRequest) (*CreateSellerResponse, error) {
	userId := fbrCtx.Get("X-User-ID")
	if userId == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	seller_id, err := h.usecase.Execute(fbrCtx.UserContext(), &domain.Seller{
		UserID:                userId,
		StoreName:             req.StoreName,
		TaxNumber:             req.TaxNumber,
		TaxOffice:             req.TaxOffice,
		LegalBusinessName:     req.LegalBusinessName,
		PhoneNumber:           req.PhoneNumber,
		AddressLine:           req.AddressLine,
		City:                  req.City,
		Country:               req.Country,
		BankAccountIban:       req.BankAccountIban,
		BankAccountHolderName: req.BankAccountHolderName,
		BankAccountBic:        req.BankAccountBic,
		Email:                 req.Email,
	})
	if err != nil {
		return nil, err
	}

	return &CreateSellerResponse{Message: "Your seller application has been successfully received. Your approval process has begun.", SellerId: seller_id}, nil
}
