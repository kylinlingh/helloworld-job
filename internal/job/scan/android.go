package scan

import (
	"context"
	"errors"
	"helloworld/internal/datastore"
	"helloworld/internal/entity"
	"helloworld/internal/pump/analytics"
	"helloworld/internal/pump/uploadto"
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
	log.Ctx(ctx).Info().Msg("android scan job started")

	a.ScanStaticCode(ctx)

	record := analytics.AnalyticsRecord{
		TimeStamp:  time.Now().Unix(),
		JobID:      "13231",
		TaskID:     "2",
		TaskTag:    "2",
		TaskResult: "41343",
	}
	uploadto.GetUploadInstance().UploadRecord(&record)

	time.Sleep(10 * time.Second)
	log.Ctx(ctx).Info().Msg("android scan job finished")
	//return nil
	return errors.New("android error")
}

func (a *AndroidScanJob) ScanStaticCode(ctx context.Context) error {
	taskRecord := entity.TaskRecord{}
	//err := a.ds.Create(ctx, &taskRecord)
	err := a.store.TaskRecord().Create(ctx, &taskRecord)
	if err != nil {

	}
	a.store.TaskRecord().Create(ctx, &taskRecord)
	return nil
}

func (a *AndroidScanJob) ScanBinaryArtifacts(ctx context.Context) error {
	return nil
}
