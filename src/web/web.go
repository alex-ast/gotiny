package web

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alex-ast/gotiny/apisrv"
	"github.com/alex-ast/gotiny/utils"
)

var baseDir = "."
var apiURL = "http://127.0.0.1:81/api"

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	// XXX: CORS security off to allow local requests from the browser.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	filePath := path.Clean(path.Join(baseDir, r.URL.Path))
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) && r.URL.Path != "/" {
		ServeRedirect(w, r)
	} else {
		http.ServeFile(w, r, filePath)
	}
}

func ServeRedirect(w http.ResponseWriter, r *http.Request) {

	shortID := strings.Trim(r.URL.Path, "/")
	if len(shortID) == 0 {
		log.Printf("[http] Bad request: %s", shortID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := http.Get(apiURL + "/url/" + shortID)
	if resp == nil {
		log.Printf("[http] API call failed: %s, %s", shortID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, apiResp := apisrv.NewGetStructures()
	utils.Unmarshal(resp.Body, &apiResp)

	if apiResp.Status.Success {
		log.Printf("[http] Redirecting (from %s): %s => %s", apiResp.Source, shortID, apiResp.URLInfo.LongURL)
		http.Redirect(w, r, apiResp.URLInfo.LongURL, http.StatusSeeOther)
		return
	}
	log.Printf("[http] Unable to redirect. HTTP status=%s, ErrMsg=%s", resp.Status, apiResp.Status.ErrorMsg)
	w.WriteHeader(http.StatusNotFound)
}

func StartWeb() {

	var port string

	flag.StringVar(&apiURL, "a", apiURL, "API URL")
	flag.StringVar(&baseDir, "b", baseDir, "Base directory for the static files")
	flag.StringVar(&port, "p", ":80", "Port to bind to")

	println("Usage: ")
	flag.PrintDefaults()

	flag.Parse()

	baseDir = path.Clean(baseDir)
	if !filepath.IsAbs(baseDir) {
		baseDir, _ = filepath.Abs(baseDir)
	}

	log.Printf("Starting server on %s, baseDir: %s", port, baseDir)
	log.Printf("API URL: %s", apiURL)

	http.HandleFunc("/", ServeHTTP)
	go func() {
		log.Fatal(http.ListenAndServe(port, nil))
	}()
}
