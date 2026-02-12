package grpc_transport

import (
	"context"
	"errors"
	"fmt"
	"marketplace/internal/payment-service/domain"
	pp "marketplace/pkg/proto/payment"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type PaymentGrpcHandler struct {
	pp.UnimplementedPaymentServiceServer
	paymentRepo   domain.PaymentRepository
	stripeService domain.StripeService
}

func NewPaymentGrpcHandler(repo domain.PaymentRepository, stripeService domain.StripeService) *PaymentGrpcHandler {
	return &PaymentGrpcHandler{
		paymentRepo:   repo,
		stripeService: stripeService,
	}
}
func (h *PaymentGrpcHandler) Register(gRPCServer *grpc.Server) {
	pp.RegisterPaymentServiceServer(gRPCServer, h)
}

func (h *PaymentGrpcHandler) CreatePaymentSession(ctx context.Context, req *pp.CreatePaymentRequest) (*pp.CreatePaymentResponse, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}
	amount := req.GetAmount()
	if amount <= 0 {
		return nil, errors.New("The amount cannot be 0.")
	}
	if req.GetUserEmail() == "" || req.GetUserName() == "" {
		return nil, errors.New("User email and name cannot be empty.")
	}

	paymentSessionRequest := domain.CreatePaymentSessionRequest{
		OrderID:   orderID,
		UserID:    userID,
		Amount:    amount,
		UserEmail: req.GetUserEmail(),
		UserName:  req.GetUserName(),
	}
	paymentURL, err := h.stripeService.CreatePaymentSession(paymentSessionRequest)
	if err != nil {
		fmt.Println("error creating payment session", err)
		return nil, err
	}

	return &pp.CreatePaymentResponse{
		PaymentUrl: paymentURL,
		SessionId:  "session id",
	}, nil

}
