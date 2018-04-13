package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

//WorkRequest is type for work request
type WorkRequest struct {
	Name     string
	Duration time.Duration
}

//Worker is custom type for each worker
type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

//WorkQueue channel to store incoming work request
var WorkQueue = make(chan WorkRequest, 100)

//Collector handler for taking up the work request over HTTP
func Collector(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//parse the delay
	delay, err := time.ParseDuration(r.FormValue("delay"))
	if err != nil {
		http.Error(w, "Bad delay value:  "+err.Error(), http.StatusBadRequest)
		return
	}

	//validate delay value
	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		http.Error(w, "Delay must be between 1 and 10 seconds inclusively", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return
	}

	// queue the work request
	work := WorkRequest{Name: name, Duration: delay}
	WorkQueue <- work
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

//StartDispatcher starts all workers and the dispatcher
func StartDispatcher(nworkers int) {
	// Common WorkeQueue for all workers to comsume the work
	WorkerQueue := make(chan chan WorkRequest, nworkers)

	for i := 0; i < nworkers; i++ {
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				fmt.Println("Recieved work request")
				go func() {
					worker := <-WorkerQueue
					fmt.Println("Dispatching the work request")
					worker <- work
				}()
			}
		}
	}()
}

var (
	//NWorkers is number of workers to start
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	//HTTPAddr is address on which to listen for HTTP requests
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
)

func main() {
	// Parse the command line flags
	flag.Parse()

	// Start the disptacher
	fmt.Println("Starting the dispatcher")
	StartDispatcher(*NWorkers)

	// Register our collector as an HTTP handler function.
	fmt.Println("Registering the collector")
	http.HandleFunc("/work", Collector)

	// Start the HTTP server
	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}
}
