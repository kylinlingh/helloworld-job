package uploadto

type UploadStorage interface {
	Connect() bool
	AppendToSetPipelined(string, [][]byte)
}
