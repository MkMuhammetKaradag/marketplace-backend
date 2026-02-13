package messaging

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "marketplace/pkg/proto/events"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

// saveCriticalMessageToStorage, DLQ'ya düşen veya kritik olan mesajları disk'e kalıcı olarak yazar.
// Neden? Kafka'da bir sorun olsa bile, kritik verilerin (Örn: Ödeme onayı) kaybolmamasını garanti eder.
// Bu dosyalar daha sonra manuel olarak veya bir script ile tekrar sisteme sokulabilir.
func (kc *KafkaClient) saveCriticalMessageToStorage(msg *pb.Message) {
	os.MkdirAll("critical_messages", 0755)

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("!!! [Storage] Marshal error: %v", err)
		return
	}

	// Hem binary (proto) hem de okunabilir (txt) formatta kaydediyoruz.
	filename := fmt.Sprintf("critical_messages/%s_%s.pb", msg.Type.String(), msg.Id)
	logName := filename + ".txt"

	humanReadable := fmt.Sprintf("ID: %s\nType: %s\nError: %s\nData: %s",
		msg.Id, msg.Type.String(), msg.LastError, msg.String())

	_ = os.WriteFile(filename, data, 0644)
	_ = os.WriteFile(logName, []byte(humanReadable), 0644)

	log.Printf("💾 [Storage] Critical message saved safely: %s", msg.Id)
}

// getConsumerGroupID, grup ID'si yoksa servise özel bir grup üretir.
func (kc *KafkaClient) getConsumerGroupID(groupID *string) string {
	if groupID != nil {
		return *groupID
	}
	return kc.serviceType.String() + "-group"
}

// getConsumerTopic, topic belirtilmemişse varsayılan topic'i döner.
func (kc *KafkaClient) getConsumerTopic(topic *string) string {
	if topic != nil {
		return *topic
	}
	return kc.config.Topic
}

func (kc *KafkaClient) createTopicsIfNotExists() error {
	conn, err := kafka.DialContext(context.Background(), "tcp", kc.config.Brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	// Controller'ı bul (Topic yaratma yetkisi ondadır)
	controller, _ := conn.Controller()
	controllerConn, _ := kafka.DialContext(context.Background(), "tcp",
		fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	defer controllerConn.Close()

	topics := []kafka.TopicConfig{
		{Topic: kc.config.Topic, NumPartitions: 3, ReplicationFactor: 1},
	}

	if kc.config.EnableRetry {
		topics = append(topics, kafka.TopicConfig{
			Topic: kc.config.RetryTopic, NumPartitions: 3, ReplicationFactor: 1})
	}

	if kc.config.DLQTopic != "" {
		topics = append(topics, kafka.TopicConfig{
			Topic: kc.config.DLQTopic, NumPartitions: 1, ReplicationFactor: 1})
	}

	return controllerConn.CreateTopics(topics...)
}

// shouldProcessMessage, mesajın bu servise gelip gelmemesi gerektiğini kontrol eder.
func (kc *KafkaClient) shouldProcessMessage(msg *pb.Message) bool {
	// 1. Hedef Servis Filtresi (ToServices)
	if len(msg.ToServices) > 0 {
		isTarget := false
		for _, svc := range msg.ToServices {
			if svc == kc.serviceType {
				isTarget = true
				break
			}
		}
		if !isTarget {
			return false
		}
	}

	// 2. Yetki Filtresi (AllowedMessageTypes)
	allowed, ok := kc.config.AllowedMessageTypes[kc.serviceType]
	if !ok {
		return false
	}

	for _, t := range allowed {
		if t == msg.Type {
			return true
		}
	}
	return false
}

// isCriticalMessageType, mesajın kritik (asla kaybolmaması gereken) tipte olup olmadığını kontrol eder.
func (kc *KafkaClient) isCriticalMessageType(msgType pb.MessageType) bool {
	for _, t := range kc.config.CriticalMessageTypes {
		if t == msgType {
			return true
		}
	}
	return false
}

// isAllowedMessageType, servis bazlı izin verilen mesaj tiplerini kontrol eder.
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
