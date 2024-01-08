package postgres

import (
	"fmt"
	"gorm.io/gorm"
	"helloworld/internal/datastore"
	"helloworld/pkg/db"
	log "helloworld/pkg/logger"
	"sync"
)

type psqlStore struct {
	db *gorm.DB
}

func (d *psqlStore) TaskRecord() datastore.TaskRecordRepo {
	return newTaskRecord(d)
}

var (
	psqlFactory datastore.DBFactory
	once        sync.Once
)

func GetPsqlFactoryOr(opts *db.PsqlOptions) (datastore.DBFactory, error) {
	if opts == nil && psqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store factory")
	}
	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dbIns, err = db.NewPsqlInstance(opts)
		psqlFactory = &psqlStore{dbIns}
	})
	if psqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get postgresql store fatory, psqlFactory: %+v, error: %w", psqlFactory, err)
	}
	log.Info().Msg("connected to postgres successfully")
	return psqlFactory, nil
}
