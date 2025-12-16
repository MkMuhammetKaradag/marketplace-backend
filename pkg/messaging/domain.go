package messaging

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

type Message struct {
	Id   string
	Type MessageType

	FromService ServiceType
	ToServices  []ServiceType
	Priority    int32
	Headers     map[string]string
	Critical    bool
	RetryCount  int32
}
