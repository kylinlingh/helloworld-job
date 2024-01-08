package mysql

import (
	"fmt"
	"gorm.io/gorm"
	"helloworld/internal/datastore"
	"helloworld/pkg/db"
	log "helloworld/pkg/logger"
	"sync"
)

type mysqlStore struct {
	db *gorm.DB
}

func (d *mysqlStore) TaskRecord() datastore.TaskRecordRepo {
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
		mysqlFactory = &mysqlStore{dbIns}
	})
	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}
	log.Info().Msg("connected to mysql successfully")
	return mysqlFactory, nil
}
