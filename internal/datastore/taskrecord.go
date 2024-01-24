package datastore

import (
	"context"
	"helloworld/internal/entity"
	v1 "helloworld/pkg/meta/v1"
)

type TaskRecordRepo interface {
	Create(ctx context.Context, tr *entity.TaskRecord) error
	List(ctx context.Context, opt v1.ListOptions) (*entity.TaskRecordList, error)
}
