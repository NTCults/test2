package main

import (
	"context"
	"test2/model"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	var event = model.Event{
		Source:    "snmp",
		Component: "server1",
		Resource:  "CPU",
		Crit:      3,
		Message:   "CPU Load > 80%",
		Timestamp: 123456789,
	}

	repo := NewMongoRepo(defaultMongoURL)
	service := NewService(repo, getLogger())
	queue := NewQueue(service, getLogger(), deafultRabbitURL)
	go queue.Listen()

	t.Run("new alert must be saved", func(t *testing.T) {
		sendEventToQueue(&event)
		time.Sleep(time.Second)

		alert, err := repo.FindAlert(context.TODO(), event.Component, event.Resource, model.StatusOngoing)
		require.NoError(t, err)
		require.NotNil(t, alert)

		require.Equal(t, event.Component, alert.Component)
		require.Equal(t, event.Resource, alert.Resource)
		require.Equal(t, event.Timestamp, alert.StartTime)
		require.Equal(t, event.Timestamp, alert.LastTime)
		require.Equal(t, event.Message, alert.LastMessage)
		require.Equal(t, event.Message, alert.FirstMessage)
		require.Equal(t, model.StatusOngoing, alert.Status)
		require.Equal(t, event.Crit, alert.Crit)
	})

	t.Run("alert must be updated", func(t *testing.T) {
		updateEvent := event
		updateEvent.Crit = 2
		updateEvent.Message = "new message"

		sendEventToQueue(&updateEvent)
		time.Sleep(time.Second)

		alert, err := repo.FindAlert(context.TODO(), event.Component, event.Resource, model.StatusOngoing)
		require.NoError(t, err)
		require.NotNil(t, alert)

		require.Equal(t, event.Component, alert.Component)
		require.Equal(t, event.Resource, alert.Resource)
		require.Equal(t, event.Timestamp, alert.StartTime)
		require.Equal(t, event.Timestamp, alert.LastTime)
		require.Equal(t, event.Message, alert.FirstMessage)
		require.Equal(t, model.StatusOngoing, alert.Status)
		// modified data
		require.Equal(t, updateEvent.Message, alert.LastMessage)
		require.Equal(t, updateEvent.Crit, alert.Crit)
	})

	t.Run("alert must be marked as resolved", func(t *testing.T) {
		resolveEvent := event
		resolveEvent.Crit = 0
		resolveEvent.Message = "resolved"
		resolveEvent.Timestamp = 321

		sendEventToQueue(&resolveEvent)
		time.Sleep(time.Second)

		alert, err := repo.FindAlert(context.TODO(), event.Component, event.Resource, model.StatusResolved)
		require.NoError(t, err)
		require.NotNil(t, alert)

		require.Equal(t, event.Component, alert.Component)
		require.Equal(t, event.Resource, alert.Resource)
		require.Equal(t, event.Timestamp, alert.StartTime)
		require.Equal(t, event.Message, alert.FirstMessage)
		// modified data
		require.Equal(t, model.StatusResolved, alert.Status)
		require.Equal(t, resolveEvent.Message, alert.LastMessage)
		require.Equal(t, resolveEvent.Crit, alert.Crit)
		require.Equal(t, resolveEvent.Timestamp, alert.LastTime)
	})

	t.Run("alert must be ignored", func(t *testing.T) {
		resolveEvent := event
		resolveEvent.Crit = 0

		sendEventToQueue(&resolveEvent)
		time.Sleep(time.Second)

		alert, err := repo.FindAlert(context.TODO(), event.Component, event.Resource, model.StatusResolved)
		require.NoError(t, err)
		require.NotNil(t, alert)

		require.Equal(t, event.Component, alert.Component)
		require.Equal(t, event.Resource, alert.Resource)
		require.Equal(t, event.Timestamp, alert.StartTime)
		require.Equal(t, event.Message, alert.FirstMessage)
		// modified data
		require.Equal(t, model.StatusResolved, alert.Status)
		require.Equal(t, resolveEvent.Message, alert.LastMessage)
		require.Equal(t, resolveEvent.Crit, alert.Crit)
		require.Equal(t, resolveEvent.Timestamp, alert.LastTime)
	})

	t.Run("alert must be ongoing again", func(t *testing.T) {
		sendEventToQueue(&event)
		time.Sleep(time.Second)

		alert, err := repo.FindAlert(context.TODO(), event.Component, event.Resource, model.StatusOngoing)
		require.NoError(t, err)
		require.NotNil(t, alert)

		require.Equal(t, event.Component, alert.Component)
		require.Equal(t, event.Resource, alert.Resource)
		require.Equal(t, event.Timestamp, alert.StartTime)
		require.Equal(t, event.Timestamp, alert.LastTime)
		require.Equal(t, event.Message, alert.LastMessage)
		require.Equal(t, event.Message, alert.FirstMessage)
		require.Equal(t, model.StatusOngoing, alert.Status)
		require.Equal(t, event.Crit, alert.Crit)
	})
}

func sendEventToQueue(event *model.Event) {
	conn, err := amqp.Dial(deafultRabbitURL)
	checkError(err)
	defer conn.Close()

	ch, err := conn.Channel()
	checkError(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
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
			Body:        createMessage(event),
		})
	checkError(err)
}
