package cache

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/alex-ast/gotiny/metrics"
	"github.com/alex-ast/gotiny/utils"
	"github.com/redis/go-redis/v9"
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

type CacheCommon struct {
	Cache
	log      *log.Logger
	counters *metrics.Counters
}

type NoCache struct {
	CacheCommon
}

type InprocCache struct {
	CacheCommon
	cache map[string][]byte
}

type RedisCache struct {
	CacheCommon
	cfg CacheCfg
	ctx context.Context
	rdb *redis.Client
}

var errKeyNotFound = errors.New("key not found")

func New(cfg CacheCfg, counters *metrics.Counters) Cache {
	cacheLogger := log.New(os.Stderr, "[cache] ", log.LstdFlags|log.Lmsgprefix)

	if strings.EqualFold(cfg.Type, "none") {
		return &NoCache{CacheCommon: CacheCommon{log: cacheLogger, counters: counters}}
	} else if strings.EqualFold(cfg.Type, "inproc") {
		storage := make(map[string][]byte)
		return &InprocCache{cache: storage, CacheCommon: CacheCommon{log: cacheLogger, counters: counters}}

	} else if strings.EqualFold(cfg.Type, "redis") {
		return &RedisCache{cfg: cfg, CacheCommon: CacheCommon{log: cacheLogger, counters: counters}}
	}
	cacheLogger.Fatal("Unknown cache type\n", cfg.Type)
	return nil
}

func (cache *NoCache) Connect() error                { return nil }
func (cache *NoCache) Set(key string, obj any) error { return nil }
func (cache *NoCache) Get(key string, obj any) error { return errKeyNotFound }
func (cache *NoCache) Delete(key string) error       { return nil }

func (cache *InprocCache) Connect() error {
	return nil
}

func (cache *InprocCache) Set(key string, obj any) error {
	cache.cache[key] = utils.MarshalToBytes(obj)
	return nil
}

func (cache *InprocCache) Get(key string, obj any) error {
	bytes, found := cache.cache[key]
	if !found {
		return errKeyNotFound
	}
	return utils.UnmarshalFromBytes(bytes, obj)
}

func (cache *InprocCache) Delete(key string) error {
	delete(cache.cache, key)
	return nil
}

func (cache *RedisCache) Connect() error {
	cache.log.Printf("Connecting to Redis %s\n", cache.cfg.Connection)
	cache.ctx = context.Background()
	cache.rdb = redis.NewClient(&redis.Options{
		Addr:     cache.cfg.Connection,
		Password: "",
		DB:       0,
	})

	return nil
}

func (cache *RedisCache) Set(key string, obj any) error {
	bytes := utils.MarshalToBytes(obj)
	log.Printf("Set(\"%s\", \"%s\")\n", key, bytes)
	err := cache.rdb.Set(cache.ctx, key, bytes, 0).Err()
	if err != nil {
		log.Println("WARN: redis cache Set failed", err)
	}
	return err
}

func (cache *RedisCache) Get(key string, obj any) error {
	var err error
	latency := metrics.Latency(func() {
		var val string
		val, err = cache.rdb.Get(cache.ctx, key).Result()
		bytes := []byte(val)
		if err == nil {
			err = utils.UnmarshalFromBytes(bytes, obj)
		}
	})
	hit := err == nil
	cache.counters.CacheLatency(hit, latency)
	return err
}

func (cache *RedisCache) Delete(key string) error {
	_, err := cache.rdb.Del(cache.ctx, key).Result()
	return err
}
