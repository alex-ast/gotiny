package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

const CONFIG_FILE = "config.yaml"

var quit = make(chan struct{})


func main() {
	config, err := LoadConfig(CONFIG_FILE)
	if err != nil {
		log.Fatal("FATAL: Can't read config.", err)
	}

	<-quit
	log.Println("Shut down")
}
