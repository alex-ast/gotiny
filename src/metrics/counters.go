package metrics

import (
	"log"
	"net/http"
	"time"
)

type MetricsCfg struct {
	// Start HTTP interface to query the counters
	StartHTTP bool `maptructure:"StartHTTP"`
	// Address to listen to for the HTTP
	Addr string `maptructure:"Addr"`
	// Path to handle for the HTTP
	Path string `maptructure:"Path"`
}

type Counters struct {
	cfg MetricsCfg
}
