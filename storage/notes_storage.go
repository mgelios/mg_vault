package storage

import (
	"context"
	"log/slog"
	"mg_vault/model"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getNotesCollectionWithFilter(filter bson.M) ([]model.Note, error) {
	slog.Debug("Getting notes")
	collection := mongo_client.Database("mg_vault").Collection("notes")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		slog.Error("Error during notes extraction")
		return nil, err
	}
	var results []model.Note
	err = cursor.All(context.Background(), &results)
	return results, err
}

func CreateNote(note model.Note) error {
	slog.Debug(note.Name)
	collection := mongo_client.Database("mg_vault").Collection("notes")
	_, err := collection.InsertOne(context.Background(), note)
	if err == nil {
		addPathEntry(note.Path, note.Author)
	}
	return err
}

func GetAllNotesForUser(userId string) ([]model.Note, error) {
	filter := bson.M{"author": userId}
	return getNotesCollectionWithFilter(filter)
}

func GetAllNotesForUserInPath(userId string, path []string) ([]model.Note, error) {
	filter := bson.M{
		"author": userId,
		"path": bson.M{
			"$eq": path,
		},
	}
	return getNotesCollectionWithFilter(filter)
}

func GetNoteById(id string) (model.Note, error) {
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
	oldNote, _ := GetNoteById(note.Id)
	noteUpdate := model.NoteUpdate{
		Name:    note.Name,
		Author:  note.Author,
		Content: note.Content,
		Tags:    note.Tags,
		Path:    note.Path,
	}
	_, err := collection.UpdateByID(context.Background(), id, bson.M{"$set": noteUpdate})

	if err == nil {
		updatePathEntry(oldNote.Path, note.Path, note.Author)
	}
	return err
}

func DeleteNoteByUserAndId(id string, userId string) error {
	collection := mongo_client.Database("mg_vault").Collection("notes")
	noteToDelete, _ := GetNoteById(id)
	idFilter, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", idFilter}}
	_, err := collection.DeleteOne(context.Background(), filter)

	if err == nil {
		removePathEntry(noteToDelete.Path, userId)
	}
	return err
}

func GetNotesTreeForUser(userId string) (model.NotesTree, error) {
	collection := mongo_client.Database("mg_vault").Collection("notes_tree")
	filter := bson.D{{"author", userId}}
	var result model.NotesTree
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	return result, err
}

func UpdateNotesTree(notesTree model.NotesTree) error {
	collection := mongo_client.Database("mg_vault").Collection("notes_tree")
	id, _ := primitive.ObjectIDFromHex(notesTree.Id)
	notesTreeUpdate := model.NotesTreeUpdate{
		Author: notesTree.Author,
		Root:   notesTree.Root,
	}
	_, err := collection.UpdateByID(context.Background(), id, bson.M{"$set": notesTreeUpdate})
	return err
}

func createNotesTreeForUser(userId string) error {
	collection := mongo_client.Database("mg_vault").Collection("notes_tree")
	notes_tree := model.NotesTree{
		Author: userId,
		Root: model.NotesTreeNode{
			ChildNodes:  map[string]*model.NotesTreeNode{},
			Breadcrumbs: []string{},
		},
	}
	_, err := collection.InsertOne(context.Background(), notes_tree)

	return err
}

func deleteNotesTreeForUser(userId string) error {
	collection := mongo_client.Database("mg_vault").Collection("notes_tree")
	filter := bson.D{{"author", userId}}
	_, err := collection.DeleteMany(context.Background(), filter)
	return err
}

func BuildNotesTree(userId string) {
	deleteNotesTreeForUser(userId)
	slog.Info("Deleted all trees for user")
	createNotesTreeForUser(userId)
	slog.Info("Create tree for user")
	notes, err := GetAllNotesForUser(userId)
	if err == nil {
		for i := 0; i < len(notes); i++ {
			addPathEntry(notes[i].Path, userId)
			slog.Info("Added new path entry")
		}
	}
}

func updatePathEntry(oldPath []string, newPath []string, userId string) {
	removePathEntry(oldPath, userId)
	addPathEntry(newPath, userId)
}

func addPathEntry(path []string, userId string) {
	notesTree, err := GetNotesTreeForUser(userId)
	if err != nil {
		BuildNotesTree(userId)
	} else {
		currentNode := &notesTree.Root
		for i := 0; i < len(path); i++ {
			if currentNode.ChildNodes[path[i]] != nil {
				currentNode = currentNode.ChildNodes[path[i]]
				currentNode.Entries += 1
			} else {
				newNode := model.NotesTreeNode{
					ChildNodes:        map[string]*model.NotesTreeNode{},
					Entries:           1,
					Breadcrumbs:       append(currentNode.Breadcrumbs, path[i]),
					BreadcrumbsString: strings.Join([]string(append(currentNode.Breadcrumbs, path[i])), ","),
				}
				currentNode.ChildNodes[path[i]] = &newNode
				currentNode = &newNode
			}
		}
		UpdateNotesTree(notesTree)
	}
}

func removePathEntry(path []string, userId string) {
	notesTree, _ := GetNotesTreeForUser(userId)
	currentNode := &notesTree.Root
	for i := 0; i < len(path); i++ {
		currentNode.ChildNodes[path[i]].Entries -= 1
		if currentNode.ChildNodes[path[i]].Entries < 1 {
			delete(currentNode.ChildNodes, path[i])
			break
		}
		currentNode = currentNode.ChildNodes[path[i]]
	}
	UpdateNotesTree(notesTree)
}
