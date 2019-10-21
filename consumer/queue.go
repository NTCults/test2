package main

import (
	"encoding/json"
	"test2/model"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const QueueName = "Alarm"

type Queue struct {
	service Service
	log     logrus.FieldLogger
	msgs    <-chan amqp.Delivery
}

func NewQueue(s Service, log logrus.FieldLogger, url string) *Queue {
	q := &Queue{
		service: s,
		log:     log,
	}
	q.subscribe(url)
	return q
}

func (q *Queue) subscribe(queueURL string) {
	conn, err := amqp.Dial(queueURL)
	checkError(err)

	ch, err := conn.Channel()
	checkError(err)

	queue, err := ch.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	checkError(err)

	q.msgs, err = ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	checkError(err)
}

// Listen starts main loop of the application.
// IMPORTANT NOTE: does not suspend the main goroutine
func (q *Queue) Listen() {
	go func() {
		for d := range q.msgs {
			err := q.handleMessage(d.Body)
			if err != nil {
				q.log.WithError(err).
					WithField("message", d).
					Error("error on handling message")
			}
			if err := d.Ack(false); err != nil {
				q.log.WithError(err).
					Error("Error on acking message")
			} else {
				q.log.Debug("Message successfully acked.")
			}
		}
	}()
	q.log.Printf("Waiting for messages. Press CTRL+C to exit")
}

func (q *Queue) handleMessage(msg []byte) error {
	var event model.Event
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	err := q.service.HandleEvent(&event)
	if err != nil {
		return err
	}

	return nil
}
