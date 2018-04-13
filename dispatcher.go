package main

import "fmt"

//StartDispatcher starts all workers and the dispatcher
func StartDispatcher(nworkers int) {
	// Common WorkeQueue for all workers to comsume the work
	WorkerQueue := make(chan chan WorkRequest, nworkers)

	//Start all the workers
	for i := 0; i < nworkers; i++ {
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	//Start dispatching the work requests
	go func() {
		for {
			select {

			// Recieve work request
			case work := <-WorkQueue:
				fmt.Println("Recieved work request")

				// Dispatch the work request
				go func() {
					worker := <-WorkerQueue
					fmt.Println("Dispatching the work request")
					worker <- work
				}()
			}
		}
	}()
}
