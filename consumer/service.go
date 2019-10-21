package main

import (
	"context"
	"test2/model"

	"github.com/sirupsen/logrus"
)

type Service interface {
	HandleEvent(alert *model.Event) error
}

type service struct {
	repo Repo
	log  logrus.FieldLogger
}

func NewService(repo Repo, log logrus.FieldLogger) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

// HandleEvent contains core logic of the service
func (s *service) HandleEvent(event *model.Event) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), defaultContextTimeout)
	defer cancelCtx()

	s.log.WithField("event", event).
		Debug("event has been received")

	alert, err := s.repo.FindAlert(ctx, event.Component, event.Resource, model.StatusOngoing)
	if err != nil {
		return err
	}

	if alert == nil {
		// skip Event because it's crit value is equal to zero
		if event.Crit < 0 {
			return nil
		}
		// create new Alert entry in repo
		return s.createAlert(ctx, event)
	}

	// update existing Alert entry
	if (alert.IsOngoing()) && (event.Crit > 0) {
		return s.updateAlert(ctx, alert, event)
	}

	// mark Alert as resolved because it's crit value evaluated zero
	if (alert.IsOngoing()) && (event.Crit == 0) {
		return s.markAlertAsResolved(ctx, alert, event)
	}

	return nil
}

func (s *service) createAlert(ctx context.Context, event *model.Event) error {
	alert := &model.Alert{
		Component:    event.Component,
		Resource:     event.Resource,
		StartTime:    event.Timestamp,
		LastTime:     event.Timestamp,
		Status:       model.StatusOngoing,
		LastMessage:  event.Message,
		FirstMessage: event.Message,
		Crit:         event.Crit,
	}

	return s.repo.SaveAlert(ctx, alert)
}

func (s *service) markAlertAsResolved(ctx context.Context, alert *model.Alert, event *model.Event) error {
	alert.LastMessage = event.Message
	alert.LastTime = event.Timestamp
	alert.Crit = 0
	alert.Status = model.StatusResolved

	return s.repo.UpdateAlert(ctx, alert)
}

func (s *service) updateAlert(ctx context.Context, alert *model.Alert, event *model.Event) error {
	alert.LastMessage = event.Message
	alert.LastTime = event.Timestamp
	alert.Crit = event.Crit

	return s.repo.UpdateAlert(ctx, alert)
}
