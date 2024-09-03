package handlers

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
)

func receiveNatsMsgHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Десериализуем тело запроса в структуру NatsMsg
	var msg nats.Msg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Обрабатываем полученный nats.Msg
	log.Printf("Received nats.Msg: Subject=%s, Data=%s\n", msg.Subject, string(msg.Data))

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message received successfully"))
}
