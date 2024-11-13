package storage

import (
	"context"
	"log/slog"
	"mg_vault/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateNote(note model.Note) error {
	slog.Debug(note.Name)
	collection := mongo_client.Database("mg_vault").Collection("notes")
	_, err := collection.InsertOne(context.Background(), note)
	return err
}

func GetAllNotesForUser(userId string) ([]model.Note, error) {
	slog.Debug("Getting quick notes")
	collection := mongo_client.Database("mg_vault").Collection("notes")
	filter := bson.D{{"author", userId}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		slog.Error("Error during notes extraction")
		return nil, err
	}
	var results []model.Note
	err = cursor.All(context.Background(), &results)
	return results, err
}

func GetNoteForUserWithId(userId string, id string) (model.Note, error) {
	collection := mongo_client.Database("mg_vault").Collection("notes")
	idFilter, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", idFilter}}
	var result model.Note
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	return result, err
}

func UpdateNote(note model.Note) error {
	collection := mongo_client.Database("mg_vault").Collection("notes")
	id, _ := primitive.ObjectIDFromHex(note.Id)
	noteUpdate := model.NoteUpdate{
		Name:    note.Name,
		Author:  note.Author,
		Content: note.Content,
		Tags:    note.Tags,
		Path:    note.Path,
	}
	_, err := collection.UpdateByID(context.Background(), id, bson.M{"$set": noteUpdate})
	return err
}

func DeleteNoteById(id string) error {
	collection := mongo_client.Database("mg_vault").Collection("notes")
	idFilter, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", idFilter}}
	_, err := collection.DeleteOne(context.Background(), filter)
	return err
}
