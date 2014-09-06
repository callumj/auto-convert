package workers

type FileRequest struct {
	Path string
}

var FileQueue = make(chan FileRequest)

func DispatchFile(req FileRequest) {
	FileQueue <- req
}
