package workers

import (
	"log"
)

type FileRequest struct {
	Uid         int64
	Path        string
	MatchedPath string
}

var FileQueue = make(chan FileRequest)

var FileWorkerQueue chan chan FileRequest

func startFileDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	FileWorkerQueue = make(chan chan FileRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		log.Println("[File] Starting worker", i+1)
		worker := NewFileWorker(i+1, FileWorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-FileQueue:
				go func() {
					worker := <-FileWorkerQueue

					log.Println("Dispatching file request")
					worker <- work
				}()
			}
		}
	}()
}

func DispatchFile(req FileRequest) {
	FileQueue <- req
}
