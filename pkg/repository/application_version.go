package repository

import (
	"context"
	"neploy.dev/pkg/model"
	"neploy.dev/pkg/store"
)

type ApplicationVersion struct {
	Base[model.ApplicationVersion]
}

func NewApplicationVersion(db store.Queryable) *ApplicationVersion {
	return &ApplicationVersion{Base[model.ApplicationVersion]{Store: db, Table: "application_versions"}}
}

func (a *ApplicationVersion) Insert(ctx context.Context, version model.ApplicationVersion) error {
	_, err := a.InsertOne(ctx, version)

	return err
}
