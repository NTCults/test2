package main

import (
	"context"
	"test2/model"
	"testing"

	"github.com/stretchr/testify/require"
)

var repo = NewMongoRepo(defaultMongoURL)

func TestRepoSave(t *testing.T) {
	testAlert := &model.Alert{
		Component: randStr10(),
		Resource:  randStr10(),
	}
	err := repo.SaveAlert(context.TODO(), testAlert)
	require.NoError(t, err)
}

func TestRepoFind(t *testing.T) {
	repo := NewMongoRepo(defaultMongoURL)
	testAlert := &model.Alert{
		Component:    randStr10(),
		Resource:     randStr10(),
		FirstMessage: "SomeTestMessage",
		Status:       model.StatusOngoing,
	}
	err := repo.SaveAlert(context.TODO(), testAlert)
	require.NoError(t, err)

	t.Run("entry exist", func(t *testing.T) {
		alert, err := repo.FindAlert(context.TODO(), testAlert.Component, testAlert.Resource, testAlert.Status)
		require.NoError(t, err)
		require.Equal(t, "SomeTestMessage", alert.FirstMessage)
	})

	t.Run("entry do not exist", func(*testing.T) {
		alert, err := repo.FindAlert(context.TODO(), "someNonExistingComponent", "someNonExistingResource", model.StatusOngoing)
		require.NoError(t, err)
		require.Nil(t, alert)
	})
}

func TestRepoUpdate(t *testing.T) {
	repo := NewMongoRepo(defaultMongoURL)
	testAlert := &model.Alert{
		Component:    randStr10(),
		Resource:     randStr10(),
		FirstMessage: randStr10(),
		Status:       model.StatusOngoing,
	}

	err := repo.SaveAlert(context.TODO(), testAlert)
	require.NoError(t, err)

	alert, err := repo.FindAlert(context.TODO(), testAlert.Component, testAlert.Resource, testAlert.Status)
	require.NoError(t, err)

	alert.LastMessage = "newMessage"
	err = repo.UpdateAlert(context.TODO(), alert)
	require.NoError(t, err)

	updatedAlert, err := repo.FindAlert(context.TODO(), testAlert.Component, testAlert.Resource, testAlert.Status)
	require.NoError(t, err)
	require.Equal(t, "newMessage", updatedAlert.LastMessage)
}
