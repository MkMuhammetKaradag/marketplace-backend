package grpc_transport

import (
	"context"
	"log"
	"marketplace/internal/user-service/domain"
	pb "marketplace/pkg/proto/auth"
	"time"

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

const SessionDuration = 24 * time.Hour
const MaxSessionLifetime = 30 * 24 * time.Hour

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

	if time.Since(session.CreatedAt) > MaxSessionLifetime {
		log.Printf("Token Maksimum Süreyi Aştı. UserID: %s", session.UserID)

		_ = h.SessionRepo.DeleteSession(ctx, token)

		return &pb.ValidationResponse{
			IsValid: false,
			Message: "Maksimum oturum süresi doldu, lütfen tekrar giriş yapın.",
		}, nil
	}

	refreshThreshold := SessionDuration / 5
	ttl, err := h.SessionRepo.GetTTL(ctx, token) // Bu metodu SessionRepo'ya eklemelisiniz.
	if err != nil {
		log.Printf("TTL alma hatası (devam ediliyor): %v", err)
		// Hata olsa bile oturum şu an geçerli, devam edebiliriz.
	} else if ttl > 0 && ttl < refreshThreshold {
		// Kalan süre eşikten azsa, süreyi yenile
		refreshErr := h.SessionRepo.RefreshSession(ctx, token, SessionDuration)
		if refreshErr != nil {
			// Yenileme başarısız olsa bile kullanıcıyı engelleme, sadece logla.
			log.Printf("Oturum yenileme başarısız: %v", refreshErr)
		} else {
			log.Printf("Oturum yenilendi. UserID: %s", session.UserID)
		}
	}

	// Oturum geçerli, kullanıcı ID'sini dön
	log.Printf("Token doğrulama başarılı. UserID: %s", session.UserID)
	return &pb.ValidationResponse{
		IsValid: true,
		UserId:  session.UserID,
		Message: "Oturum aktif.",
		Role:    string(session.Role),
	}, nil
}
