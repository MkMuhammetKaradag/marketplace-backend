package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PublishMessage, verilen protobuf mesajını Kafka'ya gönderir.
// Neden Bu Yapı? 
// 1. Otomatik ID: Mesajın ID'si yoksa UUID atar, böylece sistemde izlenebilirlik (traceability) sağlar.
// 2. Metadata: Mesajın hangi servisten çıktığını ve ne zaman oluşturulduğunu otomatik ekler.
// 3. Kritiklik Kontrolü: Mesaj tipine göre otomatik 'Critical' etiketi basar.
func (kc *KafkaClient) PublishMessage(ctx context.Context, msg *pb.Message) error {
	kc.mu.RLock()
	if kc.closed {
		kc.mu.RUnlock()
		return fmt.Errorf("kafka client is already closed")
	}
	kc.mu.RUnlock()

	// Mesaj Hazırlığı: Mesajın kimliğini ve zamanını damgalıyoruz.
	if msg.Id == "" {
		msg.Id = uuid.New().String()
	}
	if msg.Created == nil {
		msg.Created = timestamppb.Now()
	}

	// Servis Bilgisi: Hangi servisin bu mesajı bastığını kaydediyoruz.
	msg.FromService = kc.serviceType
	isCritical := kc.isCriticalMessageType(msg.Type)
	msg.Critical = isCritical

	// Serileştirme: Protobuf formatına çeviriyoruz.
	// Neden? JSON'a göre çok daha az yer kaplar ve çok daha hızlıdır.
	messageBytes, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf: %w", err)
	}

	// Kafka Mesajı Oluşturma: Key olarak ID kullanıyoruz.
	// Neden? Aynı ID'li mesajların (eğer birden fazla partition varsa) aynı partition'a gitmesini sağlar.
	kafkaMsg := kafka.Message{
		Key:   []byte(msg.Id),
		Value: messageBytes,
		Headers: []kafka.Header{
			{Key: "EventType", Value: []byte(msg.Type.String())},
			{Key: "FromService", Value: []byte(msg.FromService.String())},
		},
	}

	// Yazma İşlemi: 10 saniyelik bir timeout ile yazıyoruz.
	publishCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = kc.producer.WriteMessages(publishCtx, kafkaMsg)
	if err != nil {
		log.Printf("✗ [Producer] Publish failed [id=%s, type=%s]: %v", msg.Id, msg.Type.String(), err)
		return err
	}

	log.Printf("✓ [Producer] Published [id=%s, type=%s, critical=%v]", msg.Id, msg.Type.String(), isCritical)
	return nil
}