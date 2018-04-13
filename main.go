package main

import (
	"flag"
	"fmt"
	"net/http"
)

//WorkQueue channel to store incoming work request
var WorkQueue = make(chan WorkRequest, 100)

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
