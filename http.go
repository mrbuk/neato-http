package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mrbuk/neato-http/neato"
)

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
	isCleaning, err := robot.IsCleaning()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("content-type", "application/json")
	io.WriteString(w, fmt.Sprintf(`{"cleaning": %t}`, isCleaning))
}

// start is putting the robot in house cleaning mode in the following way
//  - check the state of robot and resume in case it is paused
//  - if the robot is not cleaning start cleaning with map
//  - if the robot is not at the base start cleaning without a map
func start(w http.ResponseWriter) {

	s, err := robot.GetState()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// if robot is busy, we try to be idempotent
	if s.State == neato.StateBusy {
		io.WriteString(w, fmt.Sprintf(`{"result": ""}`, "ok"))
		return
	}

	// if the robot is paused we try to unpause it
	if s.State == neato.StatePaused {
		r, err := robot.ResumeCleaning()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Add("content-type", "application/json")
		io.WriteString(w, fmt.Sprintf(`{"result": "%s"}`, r.Result))
		return
	}

	// start regular cleaning cycle if robot is in idle state
	if s.State == neato.StateIdle {
		// try to clean with maps
		r, err := robot.StartCleaning(false)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// in case robot is not on charge base or the command
		// to clean without maps
		if r.Result == "not_on_charge_base" {
			r, err = robot.StartCleaning(true)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
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

	r, err := robot.SendToBase()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("content-type", "application/json")
	io.WriteString(w, fmt.Sprintf(`{"result": "%s"}`, r.Result))
}
