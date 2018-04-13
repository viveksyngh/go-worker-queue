package main

import (
	"fmt"
	"time"
)

//Worker is custom type for each worker
type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

//NewWorker creates a new worker instance and returns it
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}
	return worker
}

//Start starts the worker
func (w *Worker) Start() {
	go func() {
		for {
			// Worker adds itself into the worker queue
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Recieve a work request
				fmt.Printf("worker%d: Recieved work request, delaying for %f seconds\n", w.ID, work.Duration.Seconds())

				time.Sleep(work.Duration)
				fmt.Printf("worker%d: Hello, %s!\n", w.ID, work.Name)

			case <-w.QuitChan:
				// Stop the worker
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

//Stop tells the workers to stop listening
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
