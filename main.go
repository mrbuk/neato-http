package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mrbuk/neato-http/neato"
)

var (
	robot *neato.NeatoRobot
)

func main() {
	serialNumber := os.Getenv("NEATO_ROBOT_SERIALNUMBER")
	secret := os.Getenv("NEATO_ROBOT_SECRET")

	robot = &neato.NeatoRobot{
		SerialNumber: serialNumber,
		Secret:       secret,
	}

	http.HandleFunc("/houseCleaning", HouseCleaningFunc)
	log.Println("Starting neato-http server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
