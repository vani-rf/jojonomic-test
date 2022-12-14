package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/teris-io/shortid"
	"github.com/vani-rf/jojonomic-test/microservices/input-harga-service/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error load from file, read environemt from os environment")
	}

	kafkaConn := createKafkaConn(os.Getenv("KAFKA_URL"), os.Getenv("KAFKA_TOPIC"))
	defer kafkaConn.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/input-harga", HandleInputHarga(kafkaConn)).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("start servet at ", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func HandleInputHarga(kafkaConn *kafka.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req models.InputHargaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.InputHargaResponse{
				IsError: true,
				Message: err.Error(),
			})
			return
		}

		reffId, err := shortid.Generate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.InputHargaResponse{
				IsError: true,
				Message: err.Error(),
			})
			return
		}
		req.ReffId = reffId

		payloadBytes, err := json.Marshal(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.InputHargaResponse{
				IsError: true,
				Message: err.Error(),
			})
			return
		}

		kafkaConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("address-%s", r.RemoteAddr)),
			Value: payloadBytes,
		}
		_, err = kafkaConn.WriteMessages(msg)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.InputHargaResponse{
				IsError: true,
				ReffId:  reffId,
				Message: "Kafka not ready",
			})
			return
		}

		json.NewEncoder(w).Encode(models.InputHargaResponse{
			IsError: false,
			ReffId:  reffId,
		})
	}
}

func createKafkaConn(kafkaURL, topic string) *kafka.Conn {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaURL, topic, 0)
	if err != nil {
		log.Fatal(err.Error())
	}

	return conn
}
