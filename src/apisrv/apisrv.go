package apisrv

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alex-ast/gotiny/apisrv/models"
	"github.com/alex-ast/gotiny/cache"
	"github.com/alex-ast/gotiny/db"
	dbmodels "github.com/alex-ast/gotiny/db/models"
	"github.com/alex-ast/gotiny/metrics"
	"github.com/alex-ast/gotiny/utils"

	"github.com/gin-gonic/gin"
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

type ApiServer struct {
	log      *log.Logger
	cfg      ApiCfg
	cache    cache.Cache
	db       db.Db
	counters *metrics.Counters
}

func New(cfg ApiCfg, cache cache.Cache, db db.Db, counters *metrics.Counters) *ApiServer {
	srv := &ApiServer{
		cfg:      cfg,
		cache:    cache,
		db:       db,
		counters: counters,
	}
	srv.log = log.New(os.Stderr, "[api] ", log.LstdFlags|log.Lmsgprefix)
	return srv
}

func (srv *ApiServer) Start() {

	router := gin.Default()
	router.POST(srv.cfg.Path+"/url", srv.createUrl)
	router.GET(srv.cfg.Path+"/url/:url", srv.getUrl)
	router.DELETE(srv.cfg.Path+"/url/:url", srv.deleteUrl)

	go func() {
		srv.log.Printf("Starting API server on %s/%s", srv.cfg.Addr, srv.cfg.Path)
		srv.log.Fatal(router.Run(srv.cfg.Addr))
	}()
}

// TODO: Move outside of ApiSrv. E.g. to models?
func NewCreateStructs() (models.CreateURLRequest, models.CreateURLResponse) {
	req := models.CreateURLRequest{}
	resp := models.CreateURLResponse{Status: &models.Status{}, URLInfo: &models.URLInfo{}}
	return req, resp
}

func NewGetStructures() (models.GetURLRequest, models.GetURLResponse) {
	req := models.GetURLRequest{}
	resp := models.GetURLResponse{Status: &models.Status{}, URLInfo: &models.URLInfo{}}
	return req, resp
}

type DummyType struct{}

// Delete is handled via HTTP DELETE with the single param in the query string,
// thus no "DeleteURLRequest" type
func NewDeleteStructs() (DummyType, models.DeleteURLResponse) {
	dummy := DummyType{}
	resp := models.DeleteURLResponse{Status: &models.Status{}}
	return dummy, resp
}

func (srv *ApiServer) GenerateShortId() (string, error) {
	dummy := dbmodels.Url{}
	for i := 0; i < srv.cfg.ShortUrlRetries; i++ {
		id := utils.GenRandomString(srv.cfg.ShortUrlLen)
		err := srv.db.Load(id, &dummy)
		// TODO: Check for specific 'not found'
		if err != nil {
			return id, nil
		}
	}
	return "", errors.New("Failed to generate unique id")
}

func (srv *ApiServer) createUrl(c *gin.Context) {
	defer metrics.Duration(metrics.Track("createUrl", srv.counters.ApiLatency))
	srv.counters.IncApiReqs()

	// XXX: CORS security off to allow local requests from the browser.
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Read input structure
	apiRequest, apiResponse := NewCreateStructs()
	utils.Unmarshal(c.Request.Body, &apiRequest)

	if len(apiRequest.LongURL) == 0 {
		apiResponse.Status.Success = false
		apiResponse.Status.ErrorMsg = "No URL provided"
		c.IndentedJSON(http.StatusBadRequest, apiResponse)
		return
	}

	// Generate ID and store to DB
	dbRecord := dbmodels.Url{LongURL: apiRequest.LongURL}
	dbRecord.Created = time.Now().Format(time.RFC3339)
	for i := 0; i < srv.cfg.ShortUrlRetries; i++ {
		dbRecord.ShortID = utils.GenRandomString(srv.cfg.ShortUrlLen)
		err := srv.db.Store(dbRecord.ShortID, dbRecord)
		if err == nil {
			break
		}
		dbRecord.ShortID = ""
	}
	if len(dbRecord.ShortID) == 0 {
		apiResponse.Status.Success = false
		apiResponse.Status.ErrorMsg = "Failed to generate and store short URL"
		c.IndentedJSON(http.StatusBadRequest, apiResponse)
		return
	}

	// Report back
	apiResponse.Status.Success = true
	apiResponse.URLInfo.ShortID = dbRecord.ShortID
	apiResponse.URLInfo.LongURL = apiRequest.LongURL
	apiResponse.URLInfo.Created = dbRecord.Created

	c.IndentedJSON(http.StatusOK, apiResponse)
}

func (srv *ApiServer) deleteUrl(c *gin.Context) {

	defer metrics.Duration(metrics.Track("deleteUrl", srv.counters.ApiLatency))
	srv.counters.IncApiReqs()

	// XXX: CORS security off to allow local requests from the browser.
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	shortID := c.Param("url")

	err := srv.db.Delete(shortID)
	srv.cache.Delete(shortID)

	_, apiResponse := NewDeleteStructs()

	apiResponse.Status.Success = true
	apiResponse.ShortID = shortID
	if err != nil {
		apiResponse.Status.Success = false
		apiResponse.Status.ErrorMsg = err.Error()
	}

	c.IndentedJSON(http.StatusOK, apiResponse)
}

func (srv *ApiServer) getUrl(c *gin.Context) {

	defer metrics.Duration(metrics.Track("getUrl", srv.counters.ApiLatency))
	srv.counters.IncApiReqs()

	// XXX: CORS security off to allow local requests from the browser.
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	shortID := c.Param("url")

	dbUrlInfo := dbmodels.Url{}
	_, apiResponse := NewGetStructures()

	useCache := true
	if srv.cfg.Debug && (c.Query("cache") == "0" || c.Query("nocache") == "1") {
		srv.log.Printf("[dbg] Skipping cache due to cache=0 query option")
		useCache = false
	}

	if useCache {
		err := srv.cache.Get(shortID, &apiResponse)
		if err == nil {
			apiResponse.Status.Success = true
			apiResponse.Source = "cache"
			c.IndentedJSON(http.StatusOK, apiResponse)
			return
		}
	}

	err := srv.db.Load(shortID, &dbUrlInfo)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "URL not found"})
		return
	}

	apiResponse.Status.Success = true
	apiResponse.URLInfo.ShortID = dbUrlInfo.ShortID
	apiResponse.URLInfo.LongURL = dbUrlInfo.LongURL
	apiResponse.URLInfo.Created = dbUrlInfo.Created
	//apiResponse.URLInfo.Expires = TODO
	apiResponse.Source = "db"
	c.IndentedJSON(http.StatusOK, apiResponse)

	// Update cache off the hot path
	go func() {
		srv.log.Printf("Caching results for %s", shortID)
		srv.cache.Set(shortID, apiResponse)
	}()
}
