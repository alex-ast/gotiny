package metrics

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	//
	apiReqs prometheus.Counter
	apiLat  prometheus.Summary
	//
	cacheHit  prometheus.Counter
	cacheMiss prometheus.Counter
	cacheLat  prometheus.Summary
	//
	dbHit  prometheus.Counter
	dbMiss prometheus.Counter
	dbLat  prometheus.Summary
}

func New(cfg MetricsCfg) *Counters {
	counters := &Counters{
		cfg: cfg,

		apiReqs: promauto.NewCounter(prometheus.CounterOpts{Name: "api_reqs", Help: "The total number of received requests"}),
		apiLat:  promauto.NewSummary(prometheus.SummaryOpts{Name: "api_latency", Help: "Latency of the API calls"}),

		cacheHit:  promauto.NewCounter(prometheus.CounterOpts{Name: "cache_hit", Help: "Number of cache hits"}),
		cacheMiss: promauto.NewCounter(prometheus.CounterOpts{Name: "cache_miss", Help: "Number of cache misses"}),
		cacheLat:  promauto.NewSummary(prometheus.SummaryOpts{Name: "cache_latency", Help: "Cache access latency"}),

		dbHit:  promauto.NewCounter(prometheus.CounterOpts{Name: "db_hit", Help: "Number of succesful DB requests"}),
		dbMiss: promauto.NewCounter(prometheus.CounterOpts{Name: "db_miss", Help: "Number of unsuccesful (e.g. key not found) DB requests"}),
		dbLat:  promauto.NewSummary(prometheus.SummaryOpts{Name: "db_latency", Help: "DB access latency"}),
	}
	return counters
}

func (counters *Counters) StartHTTP() {
	log.Printf("Starting monitoring HTTP at %s%s", counters.cfg.Addr, counters.cfg.Path)
	go func() {
		http.Handle(counters.cfg.Path, promhttp.Handler())
		http.ListenAndServe(counters.cfg.Addr, nil)
	}()
}

func Latency(f func()) time.Duration {
	start := time.Now()
	f()
	elapsed := time.Since(start)
	return elapsed
}

func Duration(name string, start time.Time, f func(time.Duration)) {
	elapsed := time.Since(start)
	f(elapsed)
}

func Track(name string, f func(time.Duration)) (string, time.Time, func(time.Duration)) {
	return name, time.Now(), f
}

func (counters *Counters) IncApiReqs() {
	counters.apiReqs.Inc()
}

func (counters *Counters) ApiLatency(latency time.Duration) {
	counters.apiLat.Observe(float64(latency.Microseconds()))
}

func (counters *Counters) CacheLatency(hit bool, latency time.Duration) {
	counters.cacheLat.Observe(float64(latency.Microseconds()))
	if hit {
		counters.cacheHit.Inc()
	} else {
		counters.cacheMiss.Inc()
	}
}

func (counters *Counters) DbLatency(success bool, latency time.Duration) {
	counters.dbLat.Observe(float64(latency.Microseconds()))
	if success {
		counters.dbHit.Inc()
	} else {
		counters.dbMiss.Inc()
	}
}
