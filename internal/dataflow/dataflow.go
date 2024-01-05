package dataflow

import (
	"helloworld/internal/dataflow/storage/downloadfrom"
	"helloworld/internal/dataflow/storage/uploadto"
)

type DataFlowService struct {
	upStorage uploadto.UploadStorage
	dlStorage downloadfrom.DownloadStroage
}

func New(u uploadto.UploadStorage, d downloadfrom.DownloadStroage) {

}

func CreateDownloadService() {

}

func CreateUploadService() {

}
