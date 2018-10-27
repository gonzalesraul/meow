package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gonzalesraul/meow/event"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/kit/retry"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

var cfg Config

func main() {
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	hub := newHub()
	connectEventStream(hub)
	defer event.Close()

	go hub.run()
	http.HandleFunc("/pusher", hub.handleWebSocket)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func connectEventStream(hub *Hub) {
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := event.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
		if err != nil {
			log.Println(err)
			return err
		}

		err = es.OnMeowCreated(func(m event.MeowCreatedMessage) {
			log.Printf("Meow received: %v\n", m)
			hub.broadcast(newMeowCreatedMessage(m.ID, m.Body, m.CreatedAt), nil)
		})
		if err != nil {
			log.Println(err)
			return err
		}
		event.SetEventStore(es)
		return nil
	})
}
