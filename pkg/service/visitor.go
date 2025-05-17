package service

import (
	"context"

	"neploy.dev/pkg/model"
	"neploy.dev/pkg/repository"
)

type Visitor interface {
	GetAllTraces(context.Context) ([]model.VisitorTrace, error)
	GetTraceByID(context.Context, string) (model.VisitorTrace, error)
	CreateTrace(context.Context, model.VisitorTrace) error
	UpdateTrace(context.Context, model.VisitorTrace) error
	DeleteTrace(context.Context, string) error
}

type visitor struct {
	trace *repository.VisitorTrace
}

func NewVisitor(trace *repository.VisitorTrace) Visitor {
	return &visitor{trace}
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
