package datastore

import (
	"context"
	"helloworld/internal/entity"
)

type TaskRecordRepo interface {
	Create(ctx context.Context, tr *entity.TaskRecord) error
}
