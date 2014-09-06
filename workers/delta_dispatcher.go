package workers

import (
	"log"
)

type DeltaRequest struct {
	Uid int64
}

var DeltaQueue = make(chan DeltaRequest)

var DeltaWorkerQueue chan chan DeltaRequest

func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	DeltaWorkerQueue = make(chan chan DeltaRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		log.Println("Starting worker", i+1)
		worker := NewDeltaWorker(i+1, DeltaWorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-DeltaQueue:
				go func() {
					worker := <-DeltaWorkerQueue

					log.Println("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}

func DispatchDelta(req DeltaRequest) {
	DeltaQueue <- req
}
