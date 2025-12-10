package grpc_transport

import (
	"context"
	"log"
	"marketplace/internal/user-service/domain"
	pb "marketplace/pkg/proto/auth"

	"google.golang.org/grpc"
)

// AuthGrpcHandler'ı app/application.go'dan buraya taşıyabilirsiniz veya
// sadece ValidateToken metodunu bu pakette uygulayabilirsiniz.
type AuthGrpcHandler struct {
	pb.UnimplementedAuthValidatorServer
	SessionRepo domain.SessionRepository
}

func NewAuthGrpcHandler(repo domain.SessionRepository) *AuthGrpcHandler {
	return &AuthGrpcHandler{
		SessionRepo: repo,
	}
}

// GrpcServerRegistrar arayüzüne uyan Register metodu (Application.go'da kalabilir veya buraya taşınır)
func (h *AuthGrpcHandler) Register(gRPCServer *grpc.Server) {
	pb.RegisterAuthValidatorServer(gRPCServer, h)
}

// *** ASIL DOĞRULAMA MANTIĞI BURASI ***
func (h *AuthGrpcHandler) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.ValidationResponse, error) {
	token := req.GetToken()
	if token == "" {
		return &pb.ValidationResponse{IsValid: false, Message: "Token eksik."}, nil
	}

	// SessionRepository'yi kullanarak oturumu kontrol et
	session, err := h.SessionRepo.GetSessionData(ctx, token)
	if err != nil {
		// Oturum bulunamadı (expired veya geçersiz)
		log.Printf("Token doğrulama başarısız: %v", err)
		return &pb.ValidationResponse{
			IsValid: false,
			Message: "Geçersiz veya süresi dolmuş oturum.",
		}, nil
	}

	// Oturum geçerli, kullanıcı ID'sini dön
	log.Printf("Token doğrulama başarılı. UserID: %s", session.UserID)
	return &pb.ValidationResponse{
		IsValid: true,
		UserId:  session.UserID,
		Message: "Oturum aktif.",
	}, nil
}
