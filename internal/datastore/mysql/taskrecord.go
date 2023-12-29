package mysql

import (
	"context"
	"gorm.io/gorm"
	"helloworld/internal/entity"
)

type taskrecord struct {
	db *gorm.DB
}

func (t *taskrecord) Create(ctx context.Context, tr *entity.TaskRecord) error {
	return t.db.Create(&tr).Error
}

func newTaskRecord(ds *mysqlstore) *taskrecord {
	return &taskrecord{
		ds.db,
	}
}
