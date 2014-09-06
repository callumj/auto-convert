package workers

import (
	"github.com/callumj/auto-convert/shared"
	"log"
)

func NewDeltaWorker(id int, deltaTaskQueue chan chan DeltaRequest) DeltaWorker {
	// Create, and return the worker.
	worker := DeltaWorker{
		ID:             id,
		DeltaTask:      make(chan DeltaRequest),
		DeltaTaskQueue: deltaTaskQueue,
		Shutdown:       make(chan bool)}

	return worker
}

type DeltaWorker struct {
	ID             int
	DeltaTask      chan DeltaRequest
	DeltaTaskQueue chan chan DeltaRequest
	Shutdown       chan bool
}

func (dw DeltaWorker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			dw.DeltaTaskQueue <- dw.DeltaTask

			select {
			case work := <-dw.DeltaTask:
				// Receive a work request.
				log.Printf("[%d] Will process %v\n", dw.ID, work.Uid)
				acc := shared.FetchAccount(shared.Account{Uid: work.Uid})
				if acc != nil {
					GetChangedFiles(acc)
				}
			case <-dw.Shutdown:
				// We have been asked to stop.
				log.Printf("[%d] Shutting down\n", dw.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (dw DeltaWorker) Stop() {
	go func() {
		dw.Shutdown <- true
	}()
}
