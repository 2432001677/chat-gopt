package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type Qa struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Ip       string             `bson:"ip"`
	Question string             `bson:"question"`
	Answer   string             `bson:"answer"`
	Time     time.Time          `bson:"time"`
}

var client *mongo.Client

func GetMongo() *mongo.Database {
	if client != nil {
		return client.Database("chatpyt")
	}
	uri := os.Getenv("MONGO_URI")
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client.Database("chatpyt")
}
func CloseMongo() {
	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}
