package messaging

import (
	"encoding/json"
	"fmt"
	"time"
)

type ServiceType int32

const (
	ServiceType_UNKNOWN_SERVICE     ServiceType = 0
	ServiceType_API_GATEWAY_SERVICE ServiceType = 1
	ServiceType_USER_SERVICE        ServiceType = 2
	ServiceType_SELLER_SERVICE      ServiceType = 3
	ServiceType_PRODUCT_SERVICE     ServiceType = 4
	ServiceType_ORDER_SERVICE       ServiceType = 5
	ServiceType_RETRY_SERVICE       ServiceType = 6
)

var serviceTypeToString = map[ServiceType]string{
	ServiceType_UNKNOWN_SERVICE:     "unknown",
	ServiceType_API_GATEWAY_SERVICE: "api-gateway",
	ServiceType_USER_SERVICE:        "user",
	ServiceType_SELLER_SERVICE:      "seller",
	ServiceType_PRODUCT_SERVICE:     "product",
	ServiceType_ORDER_SERVICE:       "order",
	ServiceType_RETRY_SERVICE:       "retry",
}

func (x ServiceType) String() string {
	if s, ok := serviceTypeToString[x]; ok {
		return s
	}
	return "unknown"
}

type MessageType int32

const (
	MessageType_UNKNOWN_MESSAGE_TYPE MessageType = 0
	MessageType_USER_CREATED         MessageType = 1
	MessageType_USER_DELETED         MessageType = 2
	MessageType_USER_UPDATED         MessageType = 3
	MessageType_SELLER_APPROVED      MessageType = 4
	MessageType_SELLER_REJECTED      MessageType = 5

	// MessageType_USER_CREATED MessageType = 6
)

var messageTypeToString = map[MessageType]string{
	MessageType_UNKNOWN_MESSAGE_TYPE: "unknown",
	MessageType_USER_CREATED:         "user-created",
	MessageType_USER_DELETED:         "user-deleted",
	MessageType_USER_UPDATED:         "user-updated",
	MessageType_SELLER_APPROVED:      "seller-approved",
	MessageType_SELLER_REJECTED:      "seller-rejected",
}

func (x MessageType) String() string {
	if s, ok := messageTypeToString[x]; ok {
		return s
	}
	return "unknown"
}

type Message struct {
	Id          string            `json:"id"`
	Type        MessageType       `json:"type"`
	Created     time.Time         `json:"created"`
	FromService ServiceType       `json:"from_service"`
	ToServices  []ServiceType     `json:"to_services"`
	Priority    int32             `json:"priority"`
	Headers     map[string]string `json:"headers"`
	Critical    bool              `json:"critical"`
	RetryCount  int32             `json:"retry_count"`

	Payload interface{} `json:"payload"`
}

func (m *Message) MarshalJSON() ([]byte, error) {

	// 1. Alias (Takma Ad) Tipi Tanımlama
	// Orijinal Message struct'ının tüm alanlarını kopyalar, metotlarını değil.
	type Alias Message

	// 2. Alias'ı kullanarak JSON'a dönüştürme
	// Bu, Go'nun json.Marshal fonksiyonunun Alias tipini (yani Message'ın verisini)
	// MarshalJSON metodunu tekrar çağırmadan işlemesini sağlar.

	// NOT: Payload'ın interface{} olması sorun yaratmaz,
	// içerdiği gerçek veri (örneğin SellerApprovedEvent) JSON'a çevrilir.
	messageBytes, err := json.Marshal((*Alias)(m))

	if err != nil {
		return nil, fmt.Errorf("messaging.Message JSON'a çevrilemedi (Marshal hatası): %w", err)
	}

	return messageBytes, nil
}
