package postgres

import (
	"context"
	"gorm.io/gorm"
	"helloworld/internal/entity"
)

type taskrecord struct {
	db *gorm.DB
}

func newTaskRecord(ds *psqlStore) *taskrecord {
	return &taskrecord{
		ds.db,
	}
}

func (t *taskrecord) Create(ctx context.Context, tr *entity.TaskRecord) error {
	return t.db.Create(&tr).Error
}
