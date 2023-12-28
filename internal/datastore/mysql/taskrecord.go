package mysql

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

func NewTaskRecordRepo(mysql *gorm.DB) datastore.TaskRecordRepo {
	return &taskrecord{db: mysql}
}

func newTaskRecord(ds *mysqlstore) *taskrecord {
	return &taskrecord{
		ds.db,
	}
}
