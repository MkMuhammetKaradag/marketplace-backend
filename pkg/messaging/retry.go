package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "marketplace/pkg/proto/events"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// handleFailure, mesaj işleme sırasında oluşan hataları yönetir.
// Neden Değişti? Artık reader ve kafkaMsg almıyor; çünkü mesaj ana akışta zaten commit edildi.
// Bu fonksiyon sadece mesajın tekrar denenip denenmeyeceğine veya DLQ'ya gidip gitmeyeceğine karar verir.
func (kc *KafkaClient) handleFailure(ctx context.Context, message *pb.Message, err error) {
	if err != nil {
		message.LastError = err.Error()
	}

	// Yeniden deneme (Retry) limiti dolmadıysa tekrar gönder
	if kc.shouldRetry(message) {
		message.RetryCount++
		log.Printf("⟳ [Failure] Retrying [id=%s, count=%d]", message.Id, message.RetryCount)
		kc.sendToRetry(ctx, message)
	} else {
		// Limit dolduysa mesajı Dead Letter Queue (DLQ) topic'ine at
		log.Printf("⚠ [Failure] Max retries reached, sending to DLQ [id=%s]", message.Id)
		kc.sendToDLQ(ctx, message, err)
	}
}

// sendToRetry, mesajı gecikmeli olarak tekrar işlenmek üzere Retry Topic'ine gönderir.
// Neden? 'Exponential Backoff' algoritması kullanarak sistemin (veya veritabanının)
// toparlanması için zaman tanır. İlk hata 5sn, ikinci 10sn, üçüncü 20sn bekletir.
func (kc *KafkaClient) sendToRetry(ctx context.Context, msg *pb.Message) {
	if kc.retryProducer == nil {
		log.Printf("✗ [Retry] CRITICAL: Retry producer nil! Sending [id=%s] directly to DLQ", msg.Id)
		kc.sendToDLQ(ctx, msg, fmt.Errorf("retry producer not configured"))
		return
	}

	// Gecikme süresini hesapla ve mesajın üzerine 'RetryAfter' olarak damgala
	delaySeconds := kc.calculateRetryDelay(int(msg.RetryCount))
	retryTime := time.Now().Add(time.Duration(delaySeconds) * time.Second)
	msg.RetryAfter = timestamppb.New(retryTime)

	messageBytes, _ := proto.Marshal(msg)

	retryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := kc.retryProducer.WriteMessages(retryCtx, kafka.Message{
		Key:   []byte(msg.Id),
		Value: messageBytes,
		Headers: []kafka.Header{
			{Key: "RetryCount", Value: []byte(fmt.Sprintf("%d", msg.RetryCount))},
			{Key: "RetryAfter", Value: []byte(retryTime.Format(time.RFC3339))},
		},
	})

	if err != nil {
		log.Printf("✗ [Retry] Write failed: %v", err)
		// Kafka yazamazsa yine DLQ'ya yedekle
		kc.sendToDLQ(ctx, msg, err)
	} else {
		log.Printf("⟳ [Retry] Success: Scheduled for %v (Delay: %ds)",
			retryTime.Format("15:04:05"), delaySeconds)
	}
}

// calculateRetryDelay, 'Exponential Backoff' stratejisi ile bekleme süresi üretir.
// Neden? Hata anında servisi mesaj yağmuruna tutmak yerine (thundering herd),
// aradaki süreyi katlayarak açar.
func (kc *KafkaClient) calculateRetryDelay(retryCount int) int {
	// Formül: 5 * 2^(retryCount-1) -> 5, 10, 20, 40...
	if retryCount <= 1 {
		return 5
	}
	delay := 5 << (retryCount - 1)

	// Maksimum 5 dakika (300sn) sınır koyuyoruz
	if delay > 300 {
		return 300
	}
	return delay
}

// shouldRetry, yapılandırmadaki MaxRetries değerine göre kontrol yapar.
func (kc *KafkaClient) shouldRetry(msg *pb.Message) bool {
	return kc.config.EnableRetry && int(msg.RetryCount) < kc.config.MaxRetries
}

// sendToDLQ, hata alan mesajları DLQ topic'ine gönderir.
func (kc *KafkaClient) sendToDLQ(ctx context.Context, msg *pb.Message, errReason error) {
	if kc.config.DLQTopic == "" {
		log.Printf("✗ [DLQ] Not configured for [id=%s]", msg.Id)
		return
	}

	// DLQ için geçici bir writer oluşturuyoruz (Genelde DLQ trafiği azdır)
	dlqProducer := &kafka.Writer{
		Addr:         kafka.TCP(kc.config.Brokers...),
		Topic:        kc.config.DLQTopic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	defer dlqProducer.Close()

	messageBytes, _ := proto.Marshal(msg)

	kafkaMsg := kafka.Message{
		Key:   []byte(msg.Id),
		Value: messageBytes,
		Headers: []kafka.Header{
			{Key: "ErrorReason", Value: []byte(errReason.Error())},
			{Key: "OriginalTopic", Value: []byte(kc.config.Topic)},
			{Key: "FailedAt", Value: []byte(time.Now().Format(time.RFC3339))},
		},
	}

	if err := dlqProducer.WriteMessages(ctx, kafkaMsg); err != nil {
		log.Printf("✗ [DLQ] Send failed: %v", err)
	} else {
		log.Printf("⚠ [DLQ] Message moved to DLQ: %s", msg.Id)
	}
}

// ConsumeDLQWithRecovery, DLQ'daki mesajları kurtarmaya çalışır.
func (kc *KafkaClient) ConsumeDLQWithRecovery(ctx context.Context, handler MessageHandler) error {
	if kc.config.DLQTopic == "" {
		return fmt.Errorf("DLQ topic not configured")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     kc.config.Brokers,
		GroupID:     kc.serviceType.String() + "-dlq-recovery",
		Topic:       kc.config.DLQTopic,
		StartOffset: kafka.FirstOffset, // En baştan başla ki hiçbir şey kaçmasın
	})
	defer reader.Close()

	for {
		m, err := reader.FetchMessage(ctx)
		if err != nil {
			return err
		}

		var message pb.Message
		if err := proto.Unmarshal(m.Value, &message); err != nil {
			reader.CommitMessages(ctx, m)
			continue
		}

		// KRİTİK MESAJ KONTROLÜ: Eğer kritikse hem diske yaz hem kurtar
		if kc.isCriticalMessageType(message.Type) {
			log.Printf("🆘 [Recovery] Critical message found: %s", message.Id)
			kc.saveCriticalMessageToStorage(&message)
		}

		// Handler ile tekrar dene
		if err := handler(ctx, &message); err == nil {
			log.Printf("✨ [Recovery] Success for id: %s", message.Id)
			reader.CommitMessages(ctx, m)
		}
	}
}

// handleFailureAfterDelay, gecikmeli (retry) olarak işlenen bir mesaj tekrar hata aldığında ne olacağını belirler.
func (kc *KafkaClient) handleFailureAfterDelay(ctx context.Context, message *pb.Message, err error) {
	if err != nil {
		message.LastError = err.Error()
	}

	if kc.shouldRetry(message) {
		message.RetryCount++
		log.Printf("⟳ [Retry] Re-scheduling after delayed failure [id=%s, count=%d]", message.Id, message.RetryCount)
		kc.sendToRetry(ctx, message)
	} else {
		log.Printf("⚠ [Retry] Max retries reached for delayed message [id=%s]", message.Id)
		kc.sendToDLQ(ctx, message, err)
	}
}
