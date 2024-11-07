package storage

import (
	"context"
	"log/slog"
	"mg_vault/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllQuickNotesForUser(userId string) ([]model.QuickNote, error) {
	collection := mongo_client.Database("mg_vault").Collection("quick_notes")
	filter := bson.D{{"author", userId}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		slog.Error("Error during quick notes extraction")
		return nil, err
	}
	var results []model.QuickNote
	err = cursor.All(context.Background(), &results)
	return results, err
}

func GetQuickNoteForUserWithId(userId string, id string) (model.QuickNote, error) {
	collection := mongo_client.Database("mg_vault").Collection("quick_notes")
	idFilter, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", idFilter}, {"author", userId}}
	var result model.QuickNote
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	return result, err
}

func CreateQuickNote(qnote model.QuickNote) error {
	collection := mongo_client.Database("mg_vault").Collection("quick_notes")
	_, err := collection.InsertOne(context.Background(), qnote)
	return err
}

func UpdateQuickNote(qnote model.QuickNote) error {
	collection := mongo_client.Database("mg_vault").Collection("quick_notes")
	id, _ := primitive.ObjectIDFromHex(qnote.Id)
	quickNoteUpdate := model.QuickNoteUpdate{
		Name:    qnote.Name,
		Author:  qnote.Author,
		Content: qnote.Content,
	}
	_, err := collection.UpdateByID(context.Background(), id, bson.M{"$set": quickNoteUpdate})
	return err
}
