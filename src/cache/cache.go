package cache

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
)

type CacheCfg struct {
	// Cache client to use:
	//	none - do not cache anything
	//	inpoc - in this process' memory, for debugging & development
	//	redis - use Redis
	Type string `maptructure:"Type"`
	// Redis connection string
	Connection string `maptructure:"Connection"`
}

type Cache interface {
	Connect() error
	Set(key string, obj any) error
	Get(key string, obj any) error
	Delete(key string) error
}
