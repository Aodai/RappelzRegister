package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var config configuration

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
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
	go http.ListenAndServe(fmt.Sprintf(":%d", config.WebPort), nil)
	fmt.Printf("Can be accessed at http://127.0.0.1:%d/\n", config.WebPort)
	wg.Wait()
}
