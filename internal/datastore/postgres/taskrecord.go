package postgres

import (
	"context"
	"gorm.io/gorm"
	"helloworld/internal/datastore"
	"helloworld/internal/entity"
)

type taskrecord struct {
	db *gorm.DB
}

func (t *taskrecord) Create(ctx context.Context, tr *entity.TaskRecord) error {
	return t.db.Create(&tr).Error
}

func New(mysql *gorm.DB) datastore.TaskRecordRepo {
	return &taskrecord{db: mysql}
}
