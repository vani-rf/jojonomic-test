package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/vani-rf/jojonomic-test/input-harga-storage-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load environment variable")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), nil)
	if err != nil {
		log.Fatal("Error connect to db")
	}

	r := getKafkaReader(os.Getenv("KAFKA_URL"), os.Getenv("KAFKA_TOPIC"), os.Getenv("KAFKA_GROUP_ID"))
	ctx := context.Background()
	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}
		log.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Key))
		if err := saveHarga(db, m.Value); err != nil {
			log.Println(err.Error())
			continue
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1e3,  // 1KB
		MaxBytes: 10e6, // 10MB
	})
}

func saveHarga(db *gorm.DB, data []byte) error {
	var harga models.Harga
	if err := json.Unmarshal(data, &harga); err != nil {
		return fmt.Errorf("unmarshall data error : %s", err.Error())
	}

	if err := db.Create(&harga).Error; err != nil {
		return fmt.Errorf("save data error : %s", err.Error())
	}

	log.Printf("success save data : %s", string(data))
	return nil
}
