package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/alex-ast/gotiny/apisrv"
	"github.com/alex-ast/gotiny/cache"
	"github.com/alex-ast/gotiny/db"
	"github.com/alex-ast/gotiny/metrics"
	"github.com/alex-ast/gotiny/web"
)

const CONFIG_FILE = "config.yaml"

var quit = make(chan struct{})

func SetupTermHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Received Ctrl+C/SIGTERM, shutting down")
		close(quit)
	}()
}

func main() {
	config, err := LoadConfig(CONFIG_FILE)
	if err != nil {
		log.Fatal("FATAL: Can't read config.", err)
	}

	counters := metrics.New(config.MetricsCfg)
	if config.MetricsCfg.StartHTTP {
		counters.StartHTTP()
	}

	dbi := db.New(config.DbCfg, counters)
	dbi.Connect()

	cache := cache.New(config.CacheCfg, counters)
	cache.Connect()

	apiServer := apisrv.New(config.ApiCfg, cache, dbi, counters)
	apiServer.Start()

	SetupTermHandler()
	<-quit
	log.Println("Shut down")
}
