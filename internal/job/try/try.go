package try

import (
	"context"
	"helloworld/internal/datastore"
	log "helloworld/pkg/logger"
)

type TryJob interface {
	RunJob(ctx context.Context) error
}

func CreateTryJobs(ctx context.Context, dbInstance datastore.DBFactory) []TryJob {
	log.Info().Msg("Done running try job.")
	return nil
}
