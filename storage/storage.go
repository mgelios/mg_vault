package storage

import (
	"context"
	"log/slog"
	"mg_vault/model"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initMongoClient() *mongo.Client {
	clientOption := options.Client().ApplyURI("mongodb://localhost:19000")
	if os.Getenv("MG_ENV") == "prod" {
		clientOption := options.Client().ApplyURI("mongodb://mongodb:27017")
		credential := options.Credential{
			Username: os.Getenv("MG_MONGO_USERNAME"),
			Password: os.Getenv("MG_MONGO_PASSWORD"),
		}
		clientOption = clientOption.SetAuth(credential)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		slog.Error(err.Error())
	}
	client.Database("mg_vault").CreateCollection(ctx, "user")
	client.Database("mg_vault").CreateCollection(ctx, "notes")
	client.Database("mg_vault").CreateCollection(ctx, "quick_notes")
	if err != nil {
		slog.Error("Error while creating init collections")
	}
	return client
}

var mongo_client *mongo.Client = initMongoClient()

func GetUserById(id string) (model.User, error) {
	collection := mongo_client.Database("auth").Collection("user")
	result := model.User{}
	filter := bson.D{{"id", id}}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	return result, err
}

func GetUserByUsername(username string) (model.User, error) {
	slog.Debug(username)
	collection := mongo_client.Database("mg_vault").Collection("user")
	result := model.User{}
	filter := bson.D{{"username", username}}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	return result, err
}
