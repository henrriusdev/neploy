package service

import (
	"context"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type Visitor interface {
	GetAllInfo(context.Context) ([]model.VisitorInfo, error)
	GetInfoByID(context.Context, string) (model.VisitorInfo, error)
	CreateInfo(context.Context, model.VisitorInfo) error
	UpdateInfo(context.Context, model.VisitorInfo) error
	DeleteInfo(context.Context, string) error
	GetAllTraces(context.Context) ([]model.VisitorTrace, error)
	GetTraceByID(context.Context, string) (model.VisitorTrace, error)
	CreateTrace(context.Context, model.VisitorTrace) error
	UpdateTrace(context.Context, model.VisitorTrace) error
	DeleteTrace(context.Context, string) error
}

type visitor struct {
	info  repository.VisitorInfo
	trace repository.VisitorTrace
}

func NewVisitor(info repository.VisitorInfo, trace repository.VisitorTrace) Visitor {
	return &visitor{info, trace}
}

func (v *visitor) GetAllInfo(ctx context.Context) ([]model.VisitorInfo, error) {
	return v.info.GetAll(ctx)
}

func (v *visitor) GetInfoByID(ctx context.Context, id string) (model.VisitorInfo, error) {
	return v.info.GetByID(ctx, id)
}

func (v *visitor) CreateInfo(ctx context.Context, visitorInfo model.VisitorInfo) error {
	return v.info.Insert(ctx, visitorInfo)
}

func (v *visitor) UpdateInfo(ctx context.Context, visitorInfo model.VisitorInfo) error {
	return v.info.Update(ctx, visitorInfo)
}

func (v *visitor) DeleteInfo(ctx context.Context, id string) error {
	return v.info.Delete(ctx, id)
}

func (v *visitor) GetAllTraces(ctx context.Context) ([]model.VisitorTrace, error) {
	return v.trace.GetAll(ctx)
}

func (v *visitor) GetTraceByID(ctx context.Context, id string) (model.VisitorTrace, error) {
	return v.trace.GetByID(ctx, id)
}

func (v *visitor) CreateTrace(ctx context.Context, visitorTrace model.VisitorTrace) error {
	return v.trace.Insert(ctx, visitorTrace)
}

func (v *visitor) UpdateTrace(ctx context.Context, visitorTrace model.VisitorTrace) error {
	return v.trace.Update(ctx, visitorTrace)
}

func (v *visitor) DeleteTrace(ctx context.Context, id string) error {
	return v.trace.Delete(ctx, id)
}
