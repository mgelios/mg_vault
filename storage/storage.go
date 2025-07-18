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
	slog.Info("Init of Mongo Client")
	var clientOption *options.ClientOptions
	if os.Getenv("MG_ENV") == "prod" {
		slog.Info("Setting mongo for prod env")
		clientOption = options.Client().ApplyURI("mongodb://mongodb:27017")
		credential := options.Credential{
			Username: os.Getenv("MG_MONGO_USERNAME"),
			Password: os.Getenv("MG_MONGO_PASSWORD"),
		}
		clientOption = clientOption.SetAuth(credential)
	} else {
		slog.Info("Setting mongo for dev env")
		clientOption = options.Client().ApplyURI("mongodb://localhost:19000")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		slog.Error(err.Error())
	}
	client.Database("mg_vault").CreateCollection(ctx, "user")
	client.Database("mg_vault").CreateCollection(ctx, "notes")
	client.Database("mg_vault").CreateCollection(ctx, "notes_tree")
	client.Database("mg_vault").CreateCollection(ctx, "quick_notes")
	client.Database("mg_vault").CreateCollection(ctx, "link_categories")

	if err != nil {
		slog.Error("Error while creating init collections")
	}
	slog.Info("Mongo init finished")
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
