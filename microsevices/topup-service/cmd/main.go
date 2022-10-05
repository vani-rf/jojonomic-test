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
	"github.com/vani-rf/jojonomic-test/topup-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load environment variable")
	}

	kafkaConn := createKafkaConn(os.Getenv("KAFKA_URL"), os.Getenv("KAFKA_TOPIC"))
	defer kafkaConn.Close()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), nil)
	if err != nil {
		log.Fatal("Error connect to db")
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/topup", HandleBuyback(kafkaConn, db)).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("start servet at ", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func HandleBuyback(kafkaConn *kafka.Conn, db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req models.TopupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.TopupResponse{
				IsError: true,
				Message: err.Error(),
			})
			return
		}

		rekening, err := getRekening(db, req.Norek)
		if err != nil {
			code := http.StatusInternalServerError
			if err == gorm.ErrRecordNotFound {
				code = http.StatusNotFound
			}

			w.WriteHeader(code)
			json.NewEncoder(w).Encode(models.TopupResponse{
				IsError: true,
				Message: "rekening tidak ditemukan",
			})
			return
		}

		harga, err := getHarga(db, req.Amount)
		if err != nil {
			code := http.StatusInternalServerError
			if err == gorm.ErrRecordNotFound {
				code = http.StatusNotFound
			}

			w.WriteHeader(code)
			json.NewEncoder(w).Encode(models.TopupResponse{
				IsError: true,
				Message: "harga tidak ditemukan",
			})
			return
		}

		reffID, err := shortid.Generate()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.TopupResponse{
				IsError: true,
				Message: err.Error(),
			})
			return
		}

		buybackParams := models.TopupParams{
			GoldWeight:         req.GoldWeight,
			Amount:             req.Amount,
			Norek:              req.Norek,
			ReffID:             reffID,
			HargaTopup:         harga.HargaTopup,
			HargaBuyback:       harga.HargaBuyback,
			CurrentGoldBalance: rekening.GoldBalance,
		}

		payloadBytes, err := json.Marshal(&buybackParams)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.TopupResponse{
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
			json.NewEncoder(w).Encode(models.TopupResponse{
				IsError: true,
				ReffID:  reffID,
				Message: "Kafka not ready",
			})
			return
		}

		json.NewEncoder(w).Encode(models.TopupResponse{
			IsError: false,
			ReffID:  reffID,
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

func getRekening(db *gorm.DB, norek string) (*models.Rekening, error) {
	rekening := models.Rekening{}
	if err := db.Model(rekening).Where("norek = ?", norek).First(&rekening).Error; err != nil {
		return nil, err
	}

	return &rekening, nil
}

func getHarga(db *gorm.DB, buybakcAmount float64) (*models.Harga, error) {
	harga := models.Harga{}
	if err := db.Model(harga).Where("harga_topup = ?", buybakcAmount).Order("created_at desc").First(&harga).Error; err != nil {
		return nil, err
	}

	return &harga, nil
}
