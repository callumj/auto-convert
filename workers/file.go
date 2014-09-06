package workers

import (
	//"github.com/callumj/auto-convert/lib"
	//"github.com/callumj/auto-convert/shared"
	"log"
)

func NewFileWorker(id int, fileTaskQueue chan chan FileRequest) FileWorker {
	// Create, and return the worker.
	worker := FileWorker{
		ID:            id,
		FileTask:      make(chan FileRequest),
		FileTaskQueue: fileTaskQueue,
		Shutdown:      make(chan bool)}

	return worker
}

type FileWorker struct {
	ID            int
	FileTask      chan FileRequest
	FileTaskQueue chan chan FileRequest
	Shutdown      chan bool
}

func (fw FileWorker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			fw.FileTaskQueue <- fw.FileTask

			select {
			case file := <-fw.FileTask:
				// Receive a work request.
				log.Printf("[%d] Will process %v for %v\n", fw.ID, file.Path, file.Uid)
				processFile(file)
			case <-fw.Shutdown:
				// We have been asked to stop.
				log.Printf("[%d] Shutting down\n", fw.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (fw FileWorker) Stop() {
	go func() {
		fw.Shutdown <- true
	}()
}
