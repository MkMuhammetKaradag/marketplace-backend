// internal/api-gateway/grpc_client/auth_client.go (Yeni dosya)

package grpc_client

import (
	"context"
	"log"
	"time"

	// Kendi oluşturduğunuz proto paketini import edin
	pb "marketplace/pkg/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// İstemciyi uygulamada global olarak erişilebilir tutmak için.
var AuthValidatorClient pb.AuthValidatorClient

// User Servisine olan bağlantıyı temsil eder.
var conn *grpc.ClientConn

// Gateway uygulaması başlangıcında çağrılacak fonksiyon
func InitAuthClient(grpcAddress string) error {
	var err error

	// Güvenliksiz bağlantı (Genellikle internal mikroservisler için kabul edilebilir)
	conn, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	// Bağlantı üzerinden gRPC istemcisini oluşturun
	AuthValidatorClient = pb.NewAuthValidatorClient(conn)
	log.Printf("✅ Gateway, User Servisine gRPC ile bağlandı: %s", grpcAddress)
	return nil
}

// Uygulama kapanırken bağlantıyı kapatmak için
func CloseAuthClient() {
	if conn != nil {
		conn.Close()
	}
}

// AuthMiddleware'in çağıracağı ana doğrulama fonksiyonu
func ValidateToken(token string) (isValid bool, userID string, permissions int64) {
	// 3 saniyelik timeout ile bir context oluşturun
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req := &pb.TokenRequest{Token: token}

	// User Servisindeki gRPC metodunu çağırın
	resp, err := AuthValidatorClient.ValidateToken(ctx, req)

	if err != nil {
		log.Printf("🔒 gRPC doğrulama çağrısı başarısız: %v", err)
		return false, "", 0
	}

	// Geri dönen cevabı kontrol edin
	return resp.GetIsValid(), resp.GetUserId(), resp.GetPermissions()
}
