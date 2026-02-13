package messaging

import (
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

// NewKafkaClient, iyileştirilmiş bir yapılandırma ile yeni bir client başlatır.
// Neden? Bağlantı kurulumu sırasında otomatik olarak Topic oluşturma ve
// Producer ayarlarını yapmak, manuel hata riskini ortadan kaldırır.
func NewKafkaClient(config KafkaConfig) (*KafkaClient, error) {
	kc := &KafkaClient{
		config:      config,
		serviceType: config.ServiceType,
		workerPool:  make(chan struct{}, config.MaxConcurrentHandlers),
	}

	// Topic yönetimi: Servis ayağa kalkarken topic yoksa oluşturur.
	// Avantajı: Yeni bir servis eklediğinizde Kafka panelinden manuel topic oluşturma zahmetinden kurtarır.
	if err := kc.createTopicsIfNotExists(); err != nil {
		log.Printf("Warning: Failed to create topics: %v", err)
	}

	// Producer ayarları: RequiredAcks: kafka.RequireAll kullanıldı.
	// Neden? Mesajın tüm kopyalarına (replicalara) yazıldığından emin olmak için.
	// Bu, veri kaybını (data loss) önlemek için en güvenli ayardır.
	kc.producer = &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false, // Veri güvenliği için senkron gönderim tercih edildi.
	}

	if config.EnableRetry && config.RetryTopic != "" {
		kc.retryProducer = &kafka.Writer{
			Addr:         kafka.TCP(config.Brokers...),
			Topic:        config.RetryTopic,
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: 5 * time.Second,
			RequiredAcks: kafka.RequireOne, // Retry için 'One' yeterlidir, performans sağlar
			MaxAttempts:  3,
		}
		log.Printf("✓ Retry producer initialized on topic: %s", config.RetryTopic)
	} else {
		log.Printf("⚠ Warning: Retry mechanism is DISABLED (Check EnableRetry or RetryTopic config)")
	}

	return kc, nil
}

// Close, Kafka bağlantılarını ve aktif işçileri (workers) güvenli bir şekilde kapatır.
// Neden? Bekleyen mesajların yarım kalmasını önlemek ve 'zombi process' oluşumunu engellemek için.
func (kc *KafkaClient) Close() error {
	kc.mu.Lock()
	if kc.closed {
		kc.mu.Unlock()
		return nil
	}
	kc.closed = true
	kc.mu.Unlock()

	// WaitGroup (wg) kullanarak içeride hala işlenen mesajların bitmesini bekleriz.
	kc.wg.Wait()

	if kc.producer != nil {
		return kc.producer.Close()
	}
	return nil
}

type QuietKafkaLogger struct{}

func (l QuietKafkaLogger) Printf(format string, v ...interface{}) {
	msg := strings.ToLower(format)
	// Bu kelimeler geçiyorsa loglama (Gürültüyü kes)
	noisy := []string{"no messages received", "timed out", "exceeded"}
	for _, p := range noisy {
		if strings.Contains(msg, p) {
			return
		}
	}
	log.Printf(format, v...)
}
