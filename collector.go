package main

import (
	"net/http"
	"time"
)

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
