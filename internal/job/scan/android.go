package scan

import (
	"context"
	"errors"
	"helloworld/internal/dataflow/datastructure"
	"helloworld/internal/dataflow/upload"
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
	log.Ctx(ctx).Info().Msg("android scan job started")

	a.ScanStaticCode(ctx)

	record := datastructure.AnalyticsRecord{
		TimeStamp:  time.Now().Unix(),
		JobID:      "job-1",
		TaskID:     "2",
		TaskTag:    "2",
		TaskResult: "41343",
	}
	upload.GetUploadService().UploadRecord(&record)

	record1 := datastructure.AnalyticsRecord{
		TimeStamp:  time.Now().Unix(),
		JobID:      "job-2",
		TaskID:     "2",
		TaskTag:    "2",
		TaskResult: "41343",
	}
	upload.GetUploadService().UploadRecord(&record1)

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
