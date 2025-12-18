package messaging

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// QuietKafkaLogger filters out noisy Kafka logs
type QuietKafkaLogger struct{}

func (l QuietKafkaLogger) Printf(format string, v ...interface{}) {
	msg := strings.ToLower(format)
	noisyPatterns := []string{
		"no messages received from kafka",
		"request timed out",
		"request exceeded",
		"allocated time",
	}
	
	for _, pattern := range noisyPatterns {
		if strings.Contains(msg, pattern) {
			return
		}
	}
	log.Printf(format, v...)
}

// KafkaClient handles Kafka producer and consumer operations
type KafkaClient struct {
	config      KafkaConfig
	producer    *kafka.Writer
	mu          sync.RWMutex
	closed      bool
	serviceType pb.ServiceType
	
	// Retry mechanism
	retryProducer *kafka.Writer
	
	// Worker pool for concurrent processing
	workerPool chan struct{}
	wg         sync.WaitGroup
}

// NewKafkaClient initializes a new Kafka client with improved configuration
func NewKafkaClient(config KafkaConfig) (*KafkaClient, error) {
	kc := &KafkaClient{
		config:      config,
		serviceType: config.ServiceType,
		workerPool:  make(chan struct{}, config.MaxConcurrentHandlers), // Worker pool
	}

	// Create topics first
	if err := kc.createTopicsIfNotExists(); err != nil {
		log.Printf("Warning: Failed to create topics: %v", err)
	}

	quietLogger := QuietKafkaLogger{}
	
	// Main producer
	kc.producer = &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		RequiredAcks: kafka.RequireAll,
		BatchSize:    100,
		BatchBytes:   1048576,
		BatchTimeout: 1 * time.Second,
		MaxAttempts:  3,
		Async:        false, // Synchronous for reliability
		Logger:       kafka.LoggerFunc(quietLogger.Printf),
		ErrorLogger:  kafka.LoggerFunc(log.Printf),
	}

	// Retry producer (if enabled)
	if config.EnableRetry && config.RetryTopic != "" {
		kc.retryProducer = &kafka.Writer{
			Addr:         kafka.TCP(config.Brokers...),
			Topic:        config.RetryTopic,
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: 5 * time.Second,
			RequiredAcks: kafka.RequireOne,
			MaxAttempts:  2,
			Logger:       kafka.LoggerFunc(quietLogger.Printf),
			ErrorLogger:  kafka.LoggerFunc(log.Printf),
		}
	}

	log.Printf("✓ Kafka Client initialized [service=%s, topic=%s, workers=%d]",
		config.ServiceType.String(), config.Topic, cap(kc.workerPool))
	
	return kc, nil
}

// createTopicsIfNotExists ensures all required topics exist
func (kc *KafkaClient) createTopicsIfNotExists() error {
	conn, err := kafka.DialContext(context.Background(), "tcp", kc.config.Brokers[0])
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get controller failed: %w", err)
	}

	controllerConn, err := kafka.DialContext(
		context.Background(),
		"tcp",
		fmt.Sprintf("%s:%d", controller.Host, controller.Port),
	)
	if err != nil {
		return fmt.Errorf("dial controller failed: %w", err)
	}
	defer controllerConn.Close()

	topics := []kafka.TopicConfig{
		{
			Topic:             kc.config.Topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}

	if kc.config.EnableRetry && kc.config.RetryTopic != "" {
		topics = append(topics, kafka.TopicConfig{
			Topic:             kc.config.RetryTopic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		})
	}

	if kc.config.DLQTopic != "" {
		topics = append(topics, kafka.TopicConfig{
			Topic:             kc.config.DLQTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	err = controllerConn.CreateTopics(topics...)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("create topics failed: %w", err)
	}

	for _, topic := range topics {
		log.Printf("✓ Topic '%s' ready", topic.Topic)
	}

	return nil
}

// Close gracefully shuts down the Kafka client
func (kc *KafkaClient) Close() error {
	kc.mu.Lock()
	defer kc.mu.Unlock()

	if kc.closed {
		return nil
	}
	kc.closed = true

	// Wait for all workers to finish
	kc.wg.Wait()

	var errs []error
	
	if kc.producer != nil {
		if err := kc.producer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close producer: %w", err))
		}
	}

	if kc.retryProducer != nil {
		if err := kc.retryProducer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close retry producer: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close errors: %v", errs)
	}

	log.Println("✓ Kafka Client closed gracefully")
	return nil
}

// PublishMessage sends a message to Kafka with proper error handling
func (kc *KafkaClient) PublishMessage(ctx context.Context, msg *pb.Message) error {
	kc.mu.RLock()
	if kc.closed {
		kc.mu.RUnlock()
		return fmt.Errorf("client is closed")
	}
	kc.mu.RUnlock()

	// Set message metadata
	if msg.Id == "" {
		msg.Id = uuid.New().String()
	}
	if msg.Created == nil {
		msg.Created = timestamppb.Now()
	}
	
	isCritical := kc.isCriticalMessageType(msg.Type)
	msg.Critical = isCritical
	msg.FromService = kc.serviceType

	// Marshal message
	messageBytes, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(msg.Id),
		Value: messageBytes,
		Headers: []kafka.Header{
			{Key: "EventType", Value: []byte(msg.Type.String())},
			{Key: "FromService", Value: []byte(msg.FromService.String())},
			{Key: "Critical", Value: []byte(fmt.Sprintf("%v", isCritical))},
		},
	}

	// Add timeout to context
	publishCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Publish with retries
	err = kc.producer.WriteMessages(publishCtx, kafkaMsg)
	if err != nil {
		log.Printf("✗ Publish failed [id=%s, type=%s]: %v", 
			msg.Id, msg.Type.String(), err)
		return fmt.Errorf("write message failed: %w", err)
	}

	log.Printf("✓ Published [id=%s, type=%s, critical=%v]",
		msg.Id, msg.Type.String(), isCritical)
	
	return nil
}

// MessageHandler defines the function signature for message processing
type MessageHandler func(context.Context, *pb.Message) error

// ConsumeMessages starts consuming messages with concurrent processing
func (kc *KafkaClient) ConsumeMessages(
	ctx context.Context,
	handler MessageHandler,
	topic *string,
	groupID *string,
) error {
	consumerGroupID := kc.getConsumerGroupID(groupID)
	consumerTopic := kc.getConsumerTopic(topic)

	quietLogger := QuietKafkaLogger{}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        kc.config.Brokers,
		GroupID:        consumerGroupID,
		Topic:          consumerTopic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		MaxWait:        1 * time.Second,
		CommitInterval: 1 * time.Second,
		StartOffset:    kafka.LastOffset, // Start from latest for new consumers
		Logger:         kafka.LoggerFunc(quietLogger.Printf),
		ErrorLogger:    kafka.LoggerFunc(quietLogger.Printf),
	})
	defer reader.Close()

	log.Printf("✓ Consumer started [service=%s, topic=%s, group=%s]",
		kc.serviceType.String(), consumerTopic, consumerGroupID)

	for {
		select {
		case <-ctx.Done():
			log.Println("✓ Consumer shutting down gracefully")
			kc.wg.Wait() // Wait for all in-flight messages
			return nil
		default:
			if err := kc.processMessage(ctx, reader, handler); err != nil {
				if err == context.Canceled {
					return nil
				}
				// Log error but continue consuming
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// processMessage handles a single message with worker pool
func (kc *KafkaClient) processMessage(
	ctx context.Context,
	reader *kafka.Reader,
	handler MessageHandler,
) error {
	// Fetch message with timeout
	fetchCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	m, err := reader.FetchMessage(fetchCtx)
	if err != nil {
		if err == context.Canceled || err == context.DeadlineExceeded {
			return err
		}
		if !strings.Contains(err.Error(), "request timed out") {
			log.Printf("✗ Fetch error: %v", err)
		}
		return err
	}

	// Unmarshal message
	var message pb.Message
	if err := proto.Unmarshal(m.Value, &message); err != nil {
		log.Printf("✗ Unmarshal failed [partition=%d, offset=%d]: %v",
			m.Partition, m.Offset, err)
		return reader.CommitMessages(ctx, m) // Skip bad message
	}

	// Filter messages not for this service
	if !kc.shouldProcessMessage(&message) {
		return reader.CommitMessages(ctx, m)
	}

	// Process with worker pool (concurrent processing)
	kc.workerPool <- struct{}{} // Acquire worker
	kc.wg.Add(1)

	go func() {
		defer func() {
			<-kc.workerPool // Release worker
			kc.wg.Done()
		}()

		kc.handleMessage(ctx, reader, m, &message, handler)
	}()

	return nil
}

// handleMessage processes a single message with retry/DLQ logic
func (kc *KafkaClient) handleMessage(
	ctx context.Context,
	reader *kafka.Reader,
	kafkaMsg kafka.Message,
	message *pb.Message,
	handler MessageHandler,
) {
	log.Printf("⚙ Processing [id=%s, type=%s, retry=%d]",
		message.Id, message.Type.String(), message.RetryCount)

	// Execute handler with timeout
	handlerCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := handler(handlerCtx, message)

	if err != nil {
		log.Printf("✗ Handler failed [id=%s]: %v", message.Id, err)
		kc.handleFailure(ctx, message, reader, kafkaMsg)
		return
	}

	log.Printf("✓ Processed successfully [id=%s]", message.Id)
	
	// Commit offset
	if err := reader.CommitMessages(context.Background(), kafkaMsg); err != nil {
		log.Printf("✗ Commit failed [id=%s]: %v", message.Id, err)
	}
}

// handleFailure implements retry and DLQ logic
func (kc *KafkaClient) handleFailure(
	ctx context.Context,
	message *pb.Message,
	reader *kafka.Reader,
	kafkaMsg kafka.Message,
) {
	// Check if should retry
	if kc.shouldRetry(message) {
		message.RetryCount++
		kc.sendToRetry(ctx, message)
	} else {
		// Send to DLQ
		kc.sendToDLQ(ctx, message)
	}

	// Always commit original message to avoid reprocessing
	if err := reader.CommitMessages(context.Background(), kafkaMsg); err != nil {
		log.Printf("✗ Commit failed after error [id=%s]: %v", message.Id, err)
	}
}

// sendToRetry sends message to retry topic
func (kc *KafkaClient) sendToRetry(ctx context.Context, msg *pb.Message) {
	if kc.retryProducer == nil {
		log.Printf("✗ Retry producer not configured [id=%s]", msg.Id)
		kc.sendToDLQ(ctx, msg)
		return
	}

	// Calculate exponential backoff delay
	delaySeconds := kc.calculateRetryDelay(int(msg.RetryCount))
	msg.RetryAfter = timestamppb.New(time.Now().Add(time.Duration(delaySeconds) * time.Second))

	messageBytes, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("✗ Marshal retry message failed [id=%s]: %v", msg.Id, err)
		return
	}

	retryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	kafkaMsg := kafka.Message{
		Key:   []byte(msg.Id),
		Value: messageBytes,
		Headers: []kafka.Header{
			{Key: "RetryCount", Value: []byte(fmt.Sprintf("%d", msg.RetryCount))},
			{Key: "RetryAfter", Value: []byte(msg.RetryAfter.AsTime().Format(time.RFC3339))},
		},
	}

	if err := kc.retryProducer.WriteMessages(retryCtx, kafkaMsg); err != nil {
		log.Printf("✗ Send to retry failed [id=%s]: %v", msg.Id, err)
		kc.sendToDLQ(ctx, msg)
		return
	}

	log.Printf("⟳ Sent to retry [id=%s, count=%d, delay=%ds]",
		msg.Id, msg.RetryCount, delaySeconds)
}

// sendToDLQ sends failed messages to Dead Letter Queue
func (kc *KafkaClient) sendToDLQ(ctx context.Context, msg *pb.Message) {
	if kc.config.DLQTopic == "" {
		log.Printf("✗ DLQ not configured [id=%s]", msg.Id)
		return
	}

	dlqProducer := &kafka.Writer{
		Addr:         kafka.TCP(kc.config.Brokers...),
		Topic:        kc.config.DLQTopic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 5 * time.Second,
		RequiredAcks: kafka.RequireOne,
	}
	defer dlqProducer.Close()

	messageBytes, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("✗ Marshal DLQ message failed [id=%s]: %v", msg.Id, err)
		return
	}

	dlqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	kafkaMsg := kafka.Message{
		Key:   []byte(msg.Id),
		Value: messageBytes,
		Headers: []kafka.Header{
			{Key: "OriginalTopic", Value: []byte(kc.config.Topic)},
			{Key: "FailedTimestamp", Value: []byte(time.Now().Format(time.RFC3339))},
			{Key: "FinalRetryCount", Value: []byte(fmt.Sprintf("%d", msg.RetryCount))},
		},
	}

	if err := dlqProducer.WriteMessages(dlqCtx, kafkaMsg); err != nil {
		log.Printf("✗ Send to DLQ failed [id=%s]: %v", msg.Id, err)
	} else {
		log.Printf("⚠ Sent to DLQ [id=%s]", msg.Id)
	}
}

// Helper methods

func (kc *KafkaClient) isCriticalMessageType(msgType pb.MessageType) bool {
	for _, t := range kc.config.CriticalMessageTypes {
		if t == msgType {
			return true
		}
	}
	return false
}

func (kc *KafkaClient) shouldProcessMessage(msg *pb.Message) bool {
	// Check ToServices filter
	if len(msg.ToServices) > 0 {
		found := false
		for _, svc := range msg.ToServices {
			if svc == kc.serviceType {
				found = true
				break
			}
		}
		if !found {
			log.Printf("⊘ Skipped [id=%s, not for %s]", msg.Id, kc.serviceType.String())
			return false
		}
	}

	// Check allowed message types
	if !kc.isAllowedMessageType(kc.serviceType, msg.Type) {
		log.Printf("⊘ Unauthorized type [id=%s, type=%s]", msg.Id, msg.Type.String())
		return false
	}

	return true
}

func (kc *KafkaClient) isAllowedMessageType(svcType pb.ServiceType, msgType pb.MessageType) bool {
	allowed, ok := kc.config.AllowedMessageTypes[svcType]
	if !ok {
		return false
	}
	for _, t := range allowed {
		if t == msgType {
			return true
		}
	}
	return false
}

func (kc *KafkaClient) shouldRetry(msg *pb.Message) bool {
	if !kc.config.EnableRetry {
		return false
	}
	return int(msg.RetryCount) < kc.config.MaxRetries
}

func (kc *KafkaClient) calculateRetryDelay(retryCount int) int {
	// Exponential backoff: 2^retryCount seconds (max 5 minutes)
	delay := 1 << retryCount // 2^retryCount
	if delay > 300 {
		delay = 300
	}
	return delay
}

func (kc *KafkaClient) getConsumerGroupID(groupID *string) string {
	if groupID != nil {
		return *groupID
	}
	return kc.serviceType.String() + "-group"
}

func (kc *KafkaClient) getConsumerTopic(topic *string) string {
	if topic != nil {
		return *topic
	}
	return kc.config.Topic
}