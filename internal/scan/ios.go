package scan

import (
	"context"
	"errors"
	"helloworld/internal/datastore"
	log "helloworld/pkg/logger"
	"time"
)

type IOSScanJob struct {
	store datastore.DBFactory
}

func NewIOSScanJob(ds datastore.DBFactory) *IOSScanJob {
	return &IOSScanJob{store: ds}
}

func (job *IOSScanJob) RunJob(ctx context.Context) error {
	log.Ctx(ctx).Info().Msg("ios scan service started")
	time.Sleep(15 * time.Second)
	log.Ctx(ctx).Info().Msg("ios scan service finished")
	//return nil
	return errors.New("ios error")
}

func (job *IOSScanJob) ScanStaticCode(ctx context.Context) error {
	job.store.TaskRecord().Create(ctx, nil)
	return nil
}

func (job *IOSScanJob) ScanBinaryArtifacts(ctx context.Context) error {
	return nil
}
