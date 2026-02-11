// internal/user-service/transport/http/usecase/signup.go
package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"marketplace/internal/user-service/domain"
	pb "marketplace/pkg/proto/events"
	"math/big"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type SignUpUseCase interface {
	Execute(ctx context.Context, user *domain.User) error
}
type signUpUseCase struct {
	userRepository domain.UserRepository
	messaging      domain.Messaging
}

type SignUpRequest struct {
	Username string
	Email    string
	Password string
}

func NewSignUpUseCase(repository domain.UserRepository, messaging domain.Messaging) SignUpUseCase {
	return &signUpUseCase{
		userRepository: repository,
		messaging:      messaging,
	}
}

func (u *signUpUseCase) Execute(ctx context.Context, user *domain.User) error {

	activationID := uuid.New() // Artık DB'nin RETURNING yapmasını beklemiyoruz
	code, _ := generateRandomCode(6)

	user.ActivationCode = code
	user.ActivationID = activationID // Domain modeline ID'yi set et

	// 2. Mesajı (Protobuf) hazırla
	// Dikkat: Burada senin yazdığın pb.Message yapısını koruyoruz
	data := &pb.UserActivationEmailData{
		ActivationId:   activationID.String(), // Kendi ürettiğimiz ID
		Email:          user.Email,
		Username:       user.Username,
		ActivationCode: code,
	}

	message := &pb.Message{
		Id:          uuid.New().String(),
		Type:        pb.MessageType_USER_ACTIVATION_EMAIL,
		FromService: pb.ServiceType_USER_SERVICE,
		Critical:    true,
		ToServices:  []pb.ServiceType{pb.ServiceType_NOTIFICATION_SERVICE},
		Payload:     &pb.Message_UserActivationEmailData{UserActivationEmailData: data},
	}

	// 3. Mesajı byte dizisine çevir (DB'ye kaydetmek için)
	payload, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	// 4. Repository'ye gönder
	// Artık repository'ye ID'yi de biz söylüyoruz
	_, _, err = u.userRepository.SignUpWithOutbox(ctx, user, payload)
	if err != nil {
		return err
	}
	return nil
}
func generateRandomCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length")
	}

	max := big.NewInt(10)
	code := make([]byte, length)

	for i := range code {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("random number generation failed: %w", err)
		}
		code[i] = byte(n.Int64()) + '0'
	}

	return string(code), nil
}
