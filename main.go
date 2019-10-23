package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var config configuration

func main() {

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		fmt.Printf("Received signal: %v\n", <-sigCh)
		//defer db.Close()
		//fmt.Println("Connection to DB closed!")
		os.Exit(1)
	}()

	file, err := os.Open("Site.conf")
	if err != nil {
		fmt.Println("Error: Could not open config file!")
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error: Could not decode config file!")
	}

	http.HandleFunc("/", logging(index))
	http.ListenAndServe(fmt.Sprintf(":%d", config.WebPort), nil)
}
