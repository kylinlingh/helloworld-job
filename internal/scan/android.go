package scan

import (
	"context"
	"errors"
	"helloworld/internal/datastore"
	"helloworld/internal/entity"
	log "helloworld/pkg/logger"
	"time"
)

type AndroidScanJob struct {
	store datastore.DBFactory
}

func NewAndroidScanJob(ds datastore.DBFactory) *AndroidScanJob {
	return &AndroidScanJob{store: ds}
}

func (a *AndroidScanJob) RunJob(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("android scan service started")
	time.Sleep(10 * time.Second)
	log.Ctx(ctx).Info().Msg("android scan service finished")
	//return nil
	return errors.New("android error")
}

func (a *AndroidScanJob) ScanStaticCode(ctx context.Context) error {
	taskRecord := entity.TaskRecord{}
	//err := a.ds.Create(ctx, &taskRecord)
	err := a.store.TaskRecord().Create(ctx, &taskRecord)
	if err != nil {

	}
	return nil
}

func (a *AndroidScanJob) ScanBinaryArtifacts(ctx context.Context) error {
	return nil
}
