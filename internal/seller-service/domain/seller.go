package domain

import "time"

type Seller struct {
	ID                string `json:"id"`
	UserID            string `json:"user_id"`
	StoreName         string `json:"store_name" validate:"required"`
	StoreSlug         string `json:"store_slug" validate:"required"`
	TaxNumber         string `json:"tax_number" validate:"required"`
	TaxOffice         string `json:"tax_office" validate:"required"`
	IsApproved        bool   `json:"is_approved" validate:"required"`
	StoreDescription  string `json:"store_description" validate:"required"`
	LegalBusinessName string `json:"legal_business_name" validate:"required"`
	PhoneNumber       string `json:"phone_number" validate:"required"`
	AddressLine       string `json:"address_line" validate:"required"`
	// AddressLine2 string `json:"address_line2" validate:"required"`
	City string `json:"city" validate:"required"`
	// State string `json:"state" validate:"required"`
	// ZipCode string `json:"zip_code" validate:"required"`
	Country               string    `json:"country" validate:"required"`
	BankAccountIban       string    `json:"bank_account_iban" validate:"required"`
	BankAccountHolderName string    `json:"bank_account_holder_name" validate:"required"`
	BankAccountBic        string    `json:"bank_account_bic" validate:"required"`
	Email                 string    `json:"email" validate:"required,email"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	Rating                float64   `json:"rating"`
	TotalSales            int       `json:"total_sales"`
}
