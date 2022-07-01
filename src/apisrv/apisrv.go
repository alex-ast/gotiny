package apisrv

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

type ApiCfg struct {
	// Address to bind and listen to, e.g. ":81"
	Addr string
	// Path part of the URL to handle, e.g. "/api/v1/"
	Path string
	// Length of short ID to generate
	ShortUrlLen int
	// Num retries in case of ID collision
	ShortUrlRetries int
	// Enable debug params to API
	Debug bool
}

