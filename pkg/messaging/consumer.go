package messaging

import (
	"context"
	"log"
	"time"

	pb "marketplace/pkg/proto/events"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

// ConsumeMessages, belirtilen topic'ten mesajları okumaya başlar.
// Neden? Bu fonksiyon bloklayıcıdır (for loop), mesaj geldikçe handler'a iletir.
// İyileştirme: GroupID belirtilmezse servisin adını default grup yapar.
func (kc *KafkaClient) ConsumeMessages(ctx context.Context, handler MessageHandler, topic *string, groupID *string) error {
	consumerGroupID := kc.getConsumerGroupID(groupID)
	consumerTopic := kc.getConsumerTopic(topic)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        kc.config.Brokers,
		GroupID:        consumerGroupID,
		Topic:          consumerTopic,
		StartOffset:    kafka.LastOffset, // Sadece yeni gelen mesajları oku
		CommitInterval: 1 * time.Second,  // Her saniye işlenen mesajları onayla
	})
	defer reader.Close()

	log.Printf("🚀 [Consumer] Started [service=%s, topic=%s, group=%s]", kc.serviceType.String(), consumerTopic, consumerGroupID)

	for {
		select {
		case <-ctx.Done(): // Uygulama kapanıyorsa döngüden çık
			return nil
		default:
			// Mesajı Kafka'dan çekiyoruz
			m, err := reader.FetchMessage(ctx)
			if err != nil {
				if err != context.Canceled {
					log.Printf("✗ [Consumer] Fetch error: %v", err)
				}
				continue
			}

			// Gelen mesajı protobuf nesnesine çeviriyoruz
			var message pb.Message
			if err := proto.Unmarshal(m.Value, &message); err != nil {
				log.Printf("✗ [Consumer] Unmarshal failed: %v", err)
				reader.CommitMessages(ctx, m) // Bozuk mesajı geçmek için commit et
				continue
			}

			// --- DELAY (BEKLETME) MANTIĞI ---
			// Eğer mesajın bir 'RetryAfter' zamanı varsa ve o zaman henüz gelmediyse:
			if message.RetryAfter != nil {
				now := time.Now()
				retryTime := message.RetryAfter.AsTime()

				if now.Before(retryTime) {
					waitDuration := retryTime.Sub(now)
					log.Printf("⏳ [Consumer] Delaying message [id=%s, wait=%v]", message.Id, waitDuration.Round(time.Second))

					// Mesajı Kafka'dan onaylıyoruz (Commit) çünkü hafızada bekleteceğiz.
					reader.CommitMessages(ctx, m)

					// Ayrı bir goroutine'de bekletip sonra işliyoruz.
					kc.wg.Add(1)
					go func(msg pb.Message) {
						defer kc.wg.Done()
						time.Sleep(waitDuration)
						kc.executeWithWorkerPool(ctx, &msg, handler)
					}(message)

					continue
				}
			}

			// Normal Mesaj İşleme: Worker pool'dan izin alarak çalıştır.
			kc.wg.Add(1)
			go func(kafkaMsg kafka.Message, msg pb.Message) {
				defer kc.wg.Done()
				kc.executeWithWorkerPool(ctx, &msg, handler)
				reader.CommitMessages(ctx, kafkaMsg) // İşlem bitince Kafka'ya "okundu" de.
			}(m, message)
		}
	}
}

// executeWithWorkerPool, mesajı işlerken sistem kaynaklarını korur.
// Neden? Aynı anda MaxConcurrentHandlers kadar işin yapılmasını sağlar.
func (kc *KafkaClient) executeWithWorkerPool(ctx context.Context, msg *pb.Message, handler MessageHandler) {
	// Pool'dan bir slot al (eğer doluysa burada bekler)
	kc.workerPool <- struct{}{}
	defer func() { <-kc.workerPool }() // İş bitince slotu boşalt

	log.Printf("⚙ [Worker] Processing [id=%s, type=%s]", msg.Id, msg.Type.String())

	if err := handler(ctx, msg); err != nil {
		log.Printf("✗ [Worker] Handler failed [id=%s]: %v", msg.Id, err)
		// Burada ileride retry.go içinde yazacağımız hata yönetimi devreye girecek
		kc.handleFailure(ctx, msg, err)
	} else {
		log.Printf("✓ [Worker] Processed [id=%s]", msg.Id)
	}
}

// processMessage, Kafka'dan bir mesajı çeker, filtrelerden geçirir ve
// ya hemen işler ya da gecikmeli (RetryAfter) işleme sırasına sokar.
func (kc *KafkaClient) processMessage(ctx context.Context, reader *kafka.Reader, handler MessageHandler) error {
	fetchCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m, err := reader.FetchMessage(fetchCtx)
	if err != nil {
		return err // Context iptali veya timeout
	}

	var message pb.Message
	if err := proto.Unmarshal(m.Value, &message); err != nil {
		log.Printf("✗ [Consumer] Unmarshal failed: %v", err)
		return reader.CommitMessages(ctx, m)
	}

	// 1. Filtreleme: Bu mesaj bizimle mi ilgili?
	if !kc.shouldProcessMessage(&message) {
		return reader.CommitMessages(ctx, m)
	}

	// 2. Gecikme Kontrolü (RetryAfter): Mesajın bekleme süresi doldu mu?
	if message.RetryAfter != nil {
		retryTime := message.RetryAfter.AsTime()
		if time.Now().Before(retryTime) {
			// Henüz zamanı gelmemiş, arka planda bekletip işleyelim
			kc.handleDelayedMessage(ctx, reader, m, &message, handler)
			return nil
		}
	}

	// 3. Normal İşleme: Hemen worker pool'a gönder
	kc.workerPool <- struct{}{}
	kc.wg.Add(1)
	go func() {
		defer func() {
			<-kc.workerPool
			kc.wg.Done()
		}()
		kc.handleMessage(ctx, reader, m, &message, handler)
	}()

	return nil
}

// handleMessage, bir mesajın handler tarafından işlenmesini ve sonucuna göre commit/failure sürecini yönetir.
func (kc *KafkaClient) handleMessage(
	ctx context.Context,
	reader *kafka.Reader,
	kafkaMsg kafka.Message,
	message *pb.Message,
	handler MessageHandler,
) {
	handlerCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 1. Handler'ı çalıştır
	err := handler(handlerCtx, message)

	if err != nil {
		log.Printf("✗ [Worker] Handler failed [id=%s]: %v", message.Id, err)

		kc.handleFailure(ctx, message, err)
	} else {
		log.Printf("✓ [Worker] Processed successfully [id=%s]", message.Id)
	}

	// 2. Mesajı her durumda Kafka'dan onayla (Commit)
	// Neden? Çünkü hata aldıysa zaten Retry topic'ine gönderdik veya DLQ'ya attık.
	// Orijinal topic'te asılı kalıp consumer'ı bloklamamalı.
	if commitErr := reader.CommitMessages(ctx, kafkaMsg); commitErr != nil {
		log.Printf("✗ [Worker] Commit failed [id=%s]: %v", message.Id, commitErr)
	}
}

func (kc *KafkaClient) handleDelayedMessage(ctx context.Context, reader *kafka.Reader, m kafka.Message, msg *pb.Message, handler MessageHandler) {
	waitDuration := msg.RetryAfter.AsTime().Sub(time.Now())

	// Mesajı Kafka'dan siliyoruz çünkü artık sorumluluk bizim hafızamızda (goroutine).
	reader.CommitMessages(ctx, m)

	kc.wg.Add(1)
	go func(copyMsg pb.Message) {
		defer kc.wg.Done()

		select {
		case <-time.After(waitDuration):
			kc.workerPool <- struct{}{}
			defer func() { <-kc.workerPool }()

			hCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := handler(hCtx, &copyMsg); err != nil {
				kc.handleFailureAfterDelay(context.Background(), &copyMsg, err)
			}
		case <-ctx.Done():
			return
		}
	}(*msg)
}
