package messaging

import (
	"context"
	"sync"

	pb "marketplace/pkg/proto/events"

	"github.com/segmentio/kafka-go"
)

// MessageHandler, Kafka'dan gelen mesajların işlenmesi için gereken fonksiyon imzasını tanımlar.
// Neden? Handler'ı bir tip olarak tanımlamak, farklı servislerde aynı imzayı zorunlu kılarak
// kodun polimorfik (çok biçimli) çalışmasını sağlar.
type MessageHandler func(context.Context, *pb.Message) error

// KafkaClient, Kafka üretici (producer) ve tüketici (consumer) operasyonlarını yöneten ana yapıdır.
// Neden? Tüm Kafka operasyonlarını tek bir struct altında toplamak, servislerin Kafka detaylarını
// bilmeden (encapsulation) mesaj alıp göndermesini sağlar.
type KafkaClient struct {
	config        KafkaConfig
	producer      *kafka.Writer  // Ana mesaj gönderici
	retryProducer *kafka.Writer  // Hatalı mesajları tekrar gönderen yardımcı
	mu            sync.RWMutex   // Thread-safety (eşzamanlı erişim güvenliği) için
	closed        bool           // Client'ın kapanıp kapanmadığını takip eder
	serviceType   pb.ServiceType // Hangi servisin bu client'ı kullandığı bilgisi

	// Worker pool (İşçi havuzu), sistemin aynı anda kaç mesajı işleyeceğini sınırlar.
	// Neden? Eğer binlerce mesaj gelirse, hepsine aynı anda goroutine açmak sistemi kilitler.
	// Bu havuz, kaynak tüketimini (CPU/RAM) kontrol altında tutar.
	workerPool chan struct{}
	wg         sync.WaitGroup // Uygulama kapanırken aktif işlerin bitmesini beklemek için
}

// KafkaConfig, Kafka client'ın çalışma parametrelerini içerir.
// Detaylı konfigürasyon, kodun farklı ortamlarda (Dev/Prod) değişmeden çalışmasını sağlar.
// type KafkaConfig struct {
// 	Brokers               []string
// 	Topic                 string
// 	RetryTopic            string
// 	DLQTopic              string
// 	EnableRetry           bool
// 	MaxRetries            int
// 	MaxConcurrentHandlers int // Worker pool kapasitesi
// 	ServiceType           pb.ServiceType
// 	CriticalMessageTypes  []pb.MessageType
// 	AllowedMessageTypes   map[pb.ServiceType][]pb.MessageType
// }
