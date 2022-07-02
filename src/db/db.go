package db

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	dbmodels "github.com/alex-ast/gotiny/db/models"
	"github.com/alex-ast/gotiny/metrics"
	"github.com/alex-ast/gotiny/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbCfg struct {
	// Type of database to use:
	//	inproc - mock database, keeps data in memory
	//	MongoDB - use MongoDB
	Type string `mapstructure:"Type"`
	// Connection string for MongoDB
	Connection string `mapstructure:"Connection"`
	// Database name
	Database string `mapstructure:"Database"`
}

type Db interface {
	Connect() error
	Store(key string, urlInfo dbmodels.Url) error
	Load(key string, urlInfo *dbmodels.Url) error
	Delete(key string) error
}

type DbCommon struct {
	Db
	cfg      DbCfg
	log      *log.Logger
	counters *metrics.Counters
}

type InprocDb struct {
	DbCommon
	storage map[string][]byte
}

type MongoDb struct {
	DbCommon
	ctx      context.Context
	database *mongo.Database
}

func New(cfg DbCfg, counters *metrics.Counters) Db {
	dbLogger := log.New(os.Stderr, "[db] ", log.LstdFlags|log.Lmsgprefix)

	if strings.EqualFold(cfg.Type, "inproc") {
		storage := make(map[string][]byte)
		db := &InprocDb{storage: storage, DbCommon: DbCommon{cfg: cfg, log: dbLogger, counters: counters}}
		return db
	} else if strings.EqualFold(cfg.Type, "MongoDB") {
		db := &MongoDb{DbCommon: DbCommon{cfg: cfg, log: dbLogger, counters: counters}}
		return db
	} else {
		dbLogger.Fatal("Unknown db type\n", cfg.Type)
	}

	return nil
}

func (db *InprocDb) Connect() error {
	return nil
}

func (db *InprocDb) Store(key string, urlInfo dbmodels.Url) error {
	db.storage[key] = utils.MarshalToBytes(urlInfo)
	return nil
}

func (db *InprocDb) Load(key string, urlInfo *dbmodels.Url) error {
	bytes, found := db.storage[key]
	if !found {
		return errors.New("Key not found: " + key)
	}
	return utils.UnmarshalFromBytes(bytes, urlInfo)
}

func (db *InprocDb) Delete(key string) error {
	delete(db.storage, key)
	return nil
}

func (db *MongoDb) Connect() error {
	db.ctx = context.Background()

	db.log.Printf("Connecting to MongoDB: %s", db.cfg.Connection)
	clientOptions := options.Client().ApplyURI(db.cfg.Connection)
	client, err := mongo.Connect(db.ctx, clientOptions)
	if err != nil {
		db.log.Fatal("Can't connect: " + err.Error())
		return err
	}

	err = client.Ping(db.ctx, nil)
	if err != nil {
		db.log.Fatal(err)
	}

	db.database = client.Database(db.cfg.Database)

	return nil
}

func (db *MongoDb) Store(key string, urlInfo dbmodels.Url) error {
	var err error

	latency := metrics.Latency(func() {
		urls := db.database.Collection("urls")
		_, err = urls.InsertOne(db.ctx, urlInfo)
	})

	success := (err == nil)
	db.counters.DbLatency(success, latency)
	if !success {
		db.log.Print("WARN: Error inserting URL", err)
	}
	return err
}

func (db *MongoDb) Load(key string, urlInfo *dbmodels.Url) error {
	filter := bson.D{primitive.E{Key: "shortId", Value: key}}

	urls := db.database.Collection("urls")
	res := urls.FindOne(db.ctx, filter)
	res.Decode(&urlInfo)

	return res.Err()
}

func (db *MongoDb) Delete(key string) error {
	var result *mongo.SingleResult
	latency := metrics.Latency(func() {
		filter := bson.D{primitive.E{Key: "shortId", Value: key}}

		urls := db.database.Collection("urls")
		result = urls.FindOneAndDelete(db.ctx, filter)
	})
	success := (result.Err() == nil)
	db.counters.DbLatency(success, latency)
	return result.Err()
}
