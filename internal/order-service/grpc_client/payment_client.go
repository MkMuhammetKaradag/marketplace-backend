// internal/order-service/grpc_client/product_client.go

package grpc_client

import (
	"context"

	"marketplace/internal/order-service/domain"
	pPayment "marketplace/pkg/proto/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var PaynebtServiceClient pPayment.PaymentServiceClient
var paymentConn *grpc.ClientConn

type paymentClient struct {
	client pPayment.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentClient(grpcAddress string) (domain.PaymentClient, error) {

	paymentConn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &paymentClient{
		client: pPayment.NewPaymentServiceClient(paymentConn),
		conn:   paymentConn,
	}, nil
}

func (c *paymentClient) CreatePaymentSession(ctx context.Context, orderID, userID, email string, amount float64) (*pPayment.CreatePaymentResponse, error) {

	req := &pPayment.CreatePaymentRequest{
		OrderId:   orderID,
		UserId:    userID,
		Amount:    amount,
		UserEmail: email,
		UserName:  "username",
	}
	return c.client.CreatePaymentSession(ctx, req)
}

func (c *paymentClient) Close() error {
	return c.conn.Close()
}
