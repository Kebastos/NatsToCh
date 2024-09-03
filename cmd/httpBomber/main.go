package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
)

var (
	server     = nats.DefaultURL
	subject    = "test"
	msgCount   = 1000000
	msg        = "awesome test message and message with awesome data"
	goroutines = 10
)

func main() {
	nc, err := nats.Connect(server)
	if err != nil {
		log.Fatalf("Не удалось подключиться к NATS: %v", err)
	}
	defer nc.Close()

	var f = func() {
		for i := 0; i < msgCount; i++ {
			err := nc.Publish(subject, []byte(fmt.Sprintf("Сообщение %s", msg)))
			if err != nil {
				log.Fatalf("Не удалось отправить сообщение: %v", err)
			}
		}
	}

	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go f()
	}

	wg.Wait()
	fmt.Println("All workers done")
}
