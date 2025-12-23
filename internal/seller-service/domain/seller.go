package domain

import (
	"time"
)

type Seller struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	StoreName string `json:"store_name"`
	StoreSlug string `json:"store_slug"`
	// Pointer kullanarak NULL değerleri yönetiyoruz
	StoreDescription  *string `json:"store_description"`
	Rating            float64 `json:"rating"`
	TotalSales        int     `json:"total_sales"`
	LegalBusinessName string  `json:"legal_business_name"`
	TaxNumber         string  `json:"tax_number"`
	TaxOffice         string  `json:"tax_office"`
	// status sütununu buraya alıyoruz
	Status                string `json:"status"`
	PhoneNumber           string `json:"phone_number"`
	Email                 string `json:"email"`
	AddressLine           string `json:"address_line"`
	City                  string `json:"city"`
	Country               string `json:"country"`
	BankAccountIban       string `json:"bank_account_iban"`
	BankAccountHolderName string `json:"bank_account_holder_name"`
	BankAccountBic        string `json:"bank_account_bic"`
	// NULL gelebilecek alanlar için pointer
	RejectionReason *string               `json:"rejection_reason"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	History         []SellerStatusHistory `json:"history,omitempty"`
}
type SellerStatusHistory struct {
	Status    string    `json:"status"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}
