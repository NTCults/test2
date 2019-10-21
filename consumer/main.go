package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	deafultRabbitURL      = "amqp://guest:guest@localhost:5672"
	defaultMongoURL       = "mongodb://localhost:27017"
	defaultContextTimeout = time.Second * 10
)

func main() {
	mongoURL := getEnvOrDefault("MONGO_URL", defaultMongoURL)
	repo := NewMongoRepo(mongoURL)

	log := getLogger()
	service := NewService(repo, log)

	rabbitURL := getEnvOrDefault("RABBIT_URL", deafultRabbitURL)
	q := NewQueue(service, log, rabbitURL)
	q.Listen()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Print("Server Stopped")
}
