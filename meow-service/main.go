package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gonzalesraul/meow/db"
	"github.com/gonzalesraul/meow/event"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/kit/retry"
)

//Config represents the envvar list needed to setup the meow service
type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	PostgresHost     string `envconfig:"POSTGRES_HOST"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

var cfg Config

func main() {
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	connectDatabase()
	defer db.Close()
	connectEventStream()
	defer event.Close()
	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func connectDatabase() {
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		addr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresDB)
		repo, err := db.NewPostgres(addr)
		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
}

func connectEventStream() {
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := event.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
		if err != nil {
			log.Println(err)
			return err
		}
		event.SetEventStore(es)
		return nil
	})
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/meows", createMeowHandler).Methods("POST").Queries("body", "{body}")
	return
}
