package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mrbuk/neato-http/neato"
)

func handleError(err error, w http.ResponseWriter) {
	log.Printf("error: %s", err)
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, fmt.Sprintf(`{"result": "%s", "error": "%s"}`, "error", err))
}

// HouseCleaningFunc ask RESTish meaning:
//  - GET will fetch the current cleaning state
//  - POST/PUT will start cleaning
//  - DELETE stop cleaning and move to base
func HouseCleaningFunc(w http.ResponseWriter, r *http.Request) {
	// get the state
	if r.Method == http.MethodGet {
		state(w)
	}

	// start cleaning
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		start(w)
	}

	// stop cleaning
	if r.Method == http.MethodDelete {
		stop(w)
	}
}

// state returns the current cleaning state of the robot
func state(w http.ResponseWriter) {
	log.Println("Checking cleaning state")
	isCleaning, err := robot.IsCleaning()
	if err != nil {
		handleError(err, w)
		return
	}
	s := fmt.Sprintf(`{"cleaning": %t}`, isCleaning)
	log.Printf("state = '%s'\n", s)
	w.Header().Add("content-type", "application/json")
	io.WriteString(w, s)
}

// start is putting the robot in house cleaning mode in the following way
//  - check the state of robot and resume in case it is paused
//  - if the robot is not cleaning start cleaning with map
//  - if the robot is not at the base start cleaning without a map
func start(w http.ResponseWriter) {
	log.Println("start of cleaning cycle")
	s, err := robot.GetState()
	if err != nil {
		handleError(err, w)
		return
	}

	// if robot is busy, we try to be idempotent
	if s.State == neato.StateBusy {
		log.Println("robot is cleaning already")
		w.Header().Add("content-type", "application/json")
		io.WriteString(w, fmt.Sprintf(`{"result": "%s"}`, "ok"))
		return
	}

	// if the robot is paused we try to unpause it
	if s.State == neato.StatePaused {
		log.Println("robot is paused, trying to resume cleaning")
		r, err := robot.ResumeCleaning()
		if err != nil {
			handleError(err, w)
			return
		}
		w.Header().Add("content-type", "application/json")
		io.WriteString(w, fmt.Sprintf(`{"result": "%s"}`, r.Result))
		return
	}

	// start regular cleaning cycle if robot is in idle state
	if s.State == neato.StateIdle {
		// try to clean with maps
		log.Println("robot is idle, cleaning with map")
		r, err := robot.StartCleaning(false)
		if err != nil {
			handleError(err, w)
			return
		}

		// in case robot is not on charge base or the command
		// to clean without maps
		if r.Result == "not_on_charge_base" {
			log.Println("robot is not on charge base, cleaning without map")
			r, err = robot.StartCleaning(true)
			if err != nil {
				handleError(err, w)
				return
			}
		}
		w.Header().Add("content-type", "application/json")
		io.WriteString(w, fmt.Sprintf(`{"result": "%s"}`, r.Result))
	}
}

// stop
func stop(w http.ResponseWriter) {
	// r, err := robot.StopCleaning()
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	log.Println("sending robot back to base")
	r, err := robot.SendToBase()
	if err != nil {
		handleError(err, w)
		return
	}

	w.Header().Add("content-type", "application/json")
	io.WriteString(w, fmt.Sprintf(`{"result": "%s"}`, r.Result))
}
