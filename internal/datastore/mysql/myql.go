package mysql

import (
	"fmt"
	"gorm.io/gorm"
	"helloworld/internal/datastore"
	"helloworld/pkg/db"
	"sync"
)

type mysqlstore struct {
	db *gorm.DB
}

func (d *mysqlstore) TaskRecord() datastore.TaskRecordRepo {
	return newTaskRecord(d)
}

var (
	mysqlFactory datastore.DBFactory
	once         sync.Once
)

func GetMysqlFactoryOr(opts *db.MysqlOptions) (datastore.DBFactory, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store factory")
	}
	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dbIns, err = db.NewMysqlInstance(opts)
		mysqlFactory = &mysqlstore{dbIns}
	})
	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}
	return mysqlFactory, nil
}
