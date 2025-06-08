package storage

import (
	"context"
	"log/slog"
	"mg_vault/model"
)

func initRootLinkCategoryIfAbsent() error {
	slog.Debug("Init root link called")
	collection := mongo_client.Database("mg_vault").Collection("link_categories")
	_, err := collection.InsertOne(context.Background(), model.LinkCategory{})
	if err == nil {
		slog.Debug("Sucessfuly created new root link category")
	}
	return err

}
