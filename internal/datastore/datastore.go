package datastore

type DBFactory interface {
	TaskRecord() TaskRecordRepo
}
