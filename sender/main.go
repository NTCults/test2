package main

import (
	"test2/model"

	"encoding/json"

	"github.com/streadway/amqp"
)

const defaultURL = "amqp://guest:guest@localhost:5672"

func main() {
	conn, err := amqp.Dial(defaultURL)
	checkError(err)
	defer conn.Close()

	ch, err := conn.Channel()
	checkError(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"test.queue", // name
		true,         // durable
		false,        // delete when usused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	checkError(err)

	err = ch.ExchangeDeclare(
		"events", // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	checkError(err)

	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		"events", // exchange
		false,
		nil,
	)
	checkError(err)

	err = ch.Publish(
		"events", // exchange
		"",       // routing key
		false,    // mandatorya
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        createMessage(testPayload),
		})
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func createMessage(payload *model.Event) []byte {
	data, err := json.Marshal(payload)
	checkError(err)
	return data
}

var testPayload = &model.Event{
	Source:    "snmp",
	Component: "server1",
	Resource:  "CPU",
	Crit:      2,
	Message:   "CPU Load > 80%",
	Timestamp: 123456789,
}
