package messaging

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	config      KafkaConfig
	producer    *kafka.Writer
	mu          sync.Mutex
	closed      bool
	serviceType ServiceType
}

func NewKafkaClient(config KafkaConfig) (*KafkaClient, error) {
	kc := &KafkaClient{
		config:      config,
		serviceType: config.ServiceType,
	}

	if err := kc.createTopicsIfNotExists(); err != nil {
		log.Println("Failed to create topics:", err)
	}

	producer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireAll,
		BatchSize:    100,
		BatchBytes:   1048576,
		BatchTimeout: 1 * time.Second,
		MaxAttempts:  3,

		// Logger:      kafka.LoggerFunc(quietLogger.Printf),
		// ErrorLogger: kafka.LoggerFunc(log.Printf),
	}
	kc.producer = producer

	log.Printf("Kafka Client initialized for service: %s, main topic: %s", config.ServiceType.String(), config.Topic)
	return kc, nil

}
func (kc *KafkaClient) createTopicsIfNotExists() error {
	conn, err := kafka.Dial("tcp", kc.config.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to controller: %w", err)
	}
	defer controllerConn.Close()

	// Oluşturulacak topic'ler
	topics := []kafka.TopicConfig{
		{
			Topic:             kc.config.Topic, // main-events
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}

	// Retry topic varsa ekle
	if kc.config.EnableRetry && kc.config.RetryTopic != "" {
		topics = append(topics, kafka.TopicConfig{
			Topic:             kc.config.RetryTopic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		})
	}

	// DLQ topic varsa ekle
	if kc.config.DLQTopic != "" {
		topics = append(topics, kafka.TopicConfig{
			Topic:             kc.config.DLQTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	err = controllerConn.CreateTopics(topics...)
	if err != nil {
		// Topic zaten varsa hata vermemeli
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("failed to create topics: %w", err)
		}
		log.Println("Some topics already exist, continuing...")
	}

	// Oluşturulan topic'leri logla
	for _, topic := range topics {
		log.Printf("Topic '%s' created or already exists", topic.Topic)
	}

	return nil
}
func (kc *KafkaClient) Close() error {
	kc.mu.Lock()
	defer kc.mu.Unlock()

	kc.closed = true

	if kc.producer != nil {
		if err := kc.producer.Close(); err != nil {
			return fmt.Errorf("failed to close Kafka producer: %w", err)
		}
	}
	log.Println("Kafka Client closed.")
	return nil
}

func (kc *KafkaClient) PublishMessage(ctx context.Context, msg *Message) error {
	if msg.Id == "" {
		msg.Id = uuid.New().String()
	}
	if msg.Created.IsZero() {
		msg.Created = time.Now()
	}
	isCritical := kc.isCriticalMessageType(msg.Type)
	if isCritical {
		msg.Critical = true
	}
	msg.FromService = kc.serviceType

	// Birden fazla servise gönderim simülasyonu
	// Kafka'da "routing key" olmadığı için, mesajı ana topice gönderip
	// tüketicilerin kendi filtrelemesini yapmasını bekleriz.
	// Ancak kritik mesajlar için farklı topic'lere yönlendirme yapılabilir.

	// Eğer mesaj belirli servislere gitmeli ise, bu bilgiyi mesajın kendisinde taşırız.
	// Tek bir topic'e gönderip, tüketicinin mesajı alıp almayacağına karar vermesi daha basittir.

	// Eğer hiç ToServices belirtilmemişse veya birden fazla servise gitmesi bekleniyorsa,
	// ana topic'e göndeririz.
	targetTopic := kc.config.Topic // Varsayılan olarak ana topic

	// Eğer mesaj sadece tek bir servise özel ve bu servis kritik bir mesaj bekliyorsa,
	// veya retry/DLQ senaryoları için, farklı topic'ler kullanılabilir.
	// Bu karmaşıklık genelde Kafka'da basit bir topic yapısıyla giderilir.
	// Bu örnekte, basitlik adına tüm mesajları ana topice gönderip tüketicinin filtrelemesini sağlıyoruz.
	// Kritik mesajlar için bir `critical-events` topic'i veya retry için `retry-events` topic'i olabilir.
	fmt.Println("Publishing message: ", msg)
	messageBytes, err := msg.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(msg.Id), // Mesaj anahtarı, aynı anahtara sahip mesajlar aynı bölüme gider
		Value: messageBytes,
		Headers: []kafka.Header{
			{Key: "EventType", Value: []byte(msg.Type.String())},
			{Key: "FromService", Value: []byte(msg.FromService.String())},
			// Diğer header'ları buraya ekleyebilirsiniz
		},
	}

	err = kc.producer.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		if isCritical {
			log.Printf("Failed to publish critical message to Kafka for ID %s. Saving to storage. Error: %v", msg.Id, err)

		}
		return fmt.Errorf("failed to write message to Kafka topic %s: %w", targetTopic, err)
	}

	log.Printf("Published message [ID: %s, Type: %s] to topic %s", msg.Id, msg.Type.String(), targetTopic)
	return nil
}
func (kc *KafkaClient) isCriticalMessageType(msgType MessageType) bool {
	for _, t := range kc.config.CriticalMessageTypes {
		if t == msgType {
			return true
		}
	}
	return false
}
