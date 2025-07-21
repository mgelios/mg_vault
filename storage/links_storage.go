package storage

import (
	"context"
	"log/slog"
	"mg_vault/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InitRootLinkCategoryIfAbsent() error {
	slog.Info("Init root link called")

	_, err := GetRootLinkCategory()

	if err != nil {
		slog.Info(err.Error())
		err = CreateLinkCategory(model.LinkCategory{
			Name:       "Home",
			Links:      []model.Link{},
			LinkGroups: []model.LinkGroup{},
			Parent:     primitive.NilObjectID,
		})
	}
	if err == nil {
		slog.Info("Sucessfuly created new root link category")
	}
	return err
}

func GetRootLinkCategory() (model.LinkCategory, error) {
	collection := mongo_client.Database("mg_vault").Collection("link_categories")
	filter := bson.D{{"parent", primitive.NilObjectID}}
	var result model.LinkCategory
	err := collection.FindOne(context.Background(), filter).Decode(&result)

	return result, err
}

func GetLinkCategoryById(id string) (model.LinkCategory, error) {
	collection := mongo_client.Database("mg_vault").Collection("link_categories")
	idFilter, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", idFilter}}
	var result model.LinkCategory
	err := collection.FindOne(context.Background(), filter).Decode(&result)

	return result, err
}

func GetLinkSubcategoriesByParentId(id string) ([]model.LinkCategory, error) {
	collection := mongo_client.Database("mg_vault").Collection("link_categories")
	idFilter, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"parent", idFilter}}
	cursor, err := collection.Find(context.Background(), filter)

	if err != nil {
		slog.Error(err.Error())
	}

	defer cursor.Close(context.Background())

	var subcategories []model.LinkCategory
	err = cursor.All(context.Background(), &subcategories)
	return subcategories, err
}

func CreateLinkCategoryWithParent(parent_id string) (string, error) {
	parentObjectId, _ := primitive.ObjectIDFromHex(parent_id)
	collection := mongo_client.Database("mg_vault").Collection("link_categories")

	result, err := collection.InsertOne(context.Background(), model.LinkCategory{
		Name:       "Name Placeholder",
		Links:      []model.Link{},
		LinkGroups: []model.LinkGroup{},
		Parent:     parentObjectId,
	})

	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func CreateLinkCategory(linkCategory model.LinkCategory) error {
	collection := mongo_client.Database("mg_vault").Collection("link_categories")
	_, err := collection.InsertOne(context.Background(), linkCategory)

	return err
}

func UpdateLinkCategory(linkCategory model.LinkCategory) error {
	collection := mongo_client.Database("mg_vault").Collection("link_categories")
	id, _ := primitive.ObjectIDFromHex(linkCategory.Id)
	categoryUpdate := model.LinkCategoryUpdate{
		Name:       linkCategory.Name,
		Links:      linkCategory.Links,
		LinkGroups: linkCategory.LinkGroups,
	}
	_, err := collection.UpdateByID(context.Background(), id, bson.M{"$set": categoryUpdate})

	return err
}
