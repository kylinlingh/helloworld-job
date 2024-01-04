package dataflow

import (
	"helloworld/internal/dataflow/storage/download"
	"helloworld/internal/dataflow/storage/upload"
)

type DataFlowService struct {
	upStorage upload.UploadStorage
	dlStorage download.DownloadStroage
}

func New(u upload.UploadStorage, d download.DownloadStroage) {

}
