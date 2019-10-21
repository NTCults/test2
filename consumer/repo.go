package main

import (
	"context"
	"test2/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo interface {
	FindAlert(ctx context.Context, component, resource, status string) (*model.Alert, error)
	SaveAlert(ctx context.Context, alert *model.Alert) error
	UpdateAlert(ctx context.Context, alert *model.Alert) error
}

type mongoRepo struct {
	client *mongo.Client
}

func NewMongoRepo(dbURL string) Repo {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	checkError(err)

	ctx, cancelCtx := context.WithTimeout(context.Background(), defaultContextTimeout)
	defer cancelCtx()

	err = client.Connect(ctx)
	checkError(err)

	err = client.Ping(ctx, nil)
	checkError(err)

	return &mongoRepo{
		client: client,
	}
}

// FindAlert returns an Alert entity with the specified component, resource and alert if the document is not found, returns nil value
func (r *mongoRepo) FindAlert(ctx context.Context, component, resource, status string) (*model.Alert, error) {
	var result model.Alert

	collection := r.client.Database("test").Collection("alerts")
	filter := bson.D{
		{Key: "component", Value: component},
		{Key: "resource", Value: resource},
		{Key: "status", Value: status},
	}

	if err := collection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *mongoRepo) SaveAlert(ctx context.Context, alert *model.Alert) error {
	collection := r.client.Database("test").Collection("alerts")
	_, err := collection.InsertOne(context.TODO(), alert)
	return err
}

func (r *mongoRepo) UpdateAlert(ctx context.Context, alert *model.Alert) error {
	collection := r.client.Database("test").Collection("alerts")

	filter := bson.D{
		{Key: "component", Value: alert.Component},
		{Key: "resource", Value: alert.Resource},
	}

	update := bson.D{
		{Key: "$set", Value: alert},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}
