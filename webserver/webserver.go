package webserver

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/BrentNoorda/slfsrv/bundle"
	"github.com/BrentNoorda/slfsrv/ssutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type jsonResponse map[string]interface{} // from http://nesv.blogspot.com/2012/09/super-easy-json-http-responses-in-go.html

type webserverData struct {
	keepAliveChan chan int
	verbose       bool
	fullRootPath  string
	initFile      string
	storeFilespec string
	zipReader     *zip.Reader // will be nil if reading from raw files
}

func (r jsonResponse) String() (s string) {
	b, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

var mimeTypes map[string]string = map[string]string{
	"js":   "application/javascript;charset=UTF-8",
	"html": "text/html;charset=UTF-8",
	"htm":  "text/html;charset=UTF-8",
	"css":  "text/css",
	"png":  "image/png",
	"gif":  "image/gif",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
}

func getMimeType(filespec string) (mimeType string) { // return "" if not known
	var lastDot int = strings.LastIndex(filespec, ".")
	if -1 != lastDot {
		var fileExtension = strings.ToLower(filespec[lastDot+1:])
		mimeType = mimeTypes[fileExtension]
	}
	return
}

func web_server_forever(secretKey string, wsData *webserverData,
	port int, rootPath string, myselfExecutable string) {
	var lenSecretKey int = len(secretKey)
	if wsData.verbose {
		fmt.Printf("Listen forever on port %d...\n", port)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, status int) {
		// from http://stackoverflow.com/questions/9996767/showing-custom-404-error-page-with-standard-http-package
		w.WriteHeader(status)
		if status == http.StatusNotFound {
			fmt.Fprint(w, "custom 404")
		}
	}

	self_serving_core_js_handler := func(w http.ResponseWriter, r *http.Request) {
		var s string = selfServingJsSrc()
		s = strings.Replace(s, "<<<SECRET_KEY>>>", secretKey, 1)
		s = strings.Replace(s, "<<<OS>>>", runtime.GOOS, 1)
		s = strings.Replace(s, "0/*<<<PORT>>>*/", strconv.Itoa(port), 1)
		s = strings.Replace(s, "<<<ROOTPATH>>>", strings.Replace(filepath.FromSlash(wsData.fullRootPath), "\\", "\\\\", -1), 1)
		s = strings.Replace(s, "<<<INITFILE>>>", wsData.initFile, 1)
		s = strings.Replace(s, "<<<SELF>>>", strings.Replace(filepath.FromSlash(myselfExecutable), "\\", "\\\\", -1), 1)

		s = strings.Replace(s, "/* INSERT JSON2.JS HERE */", json2_jsselfServingJsSrc(), 1)

		s = strings.Replace(s, "/* INSERT NATIVE_PROMISE_ONLY.JS HERE */", native_promise_only_jsselfServingJsSrc(), 1)

		w.Header().Set("Content-Type", "application/javascript;charset=UTF-8")
		fmt.Fprintf(w, "%s", s)
	}

	favicon := func(w http.ResponseWriter, r *http.Request) {

		var s string = faviconData()
		w.Header().Set("Content-Type", "image/x-icon")
		fmt.Fprintf(w, "%s", s)
	}

	noSecretKeyHandler := func(w http.ResponseWriter, r *http.Request) {
		if wsData.verbose {
			fmt.Fprintf(os.Stderr, "    No secret key for URL request \"%s\"\n", r.URL)
		}
		errorHandler(w, r, http.StatusNotFound)
	}

	callHandler := func(w http.ResponseWriter, r *http.Request) {
		var jResp jsonResponse
		var err error
		var paramsUrl = r.URL.Path[lenSecretKey+7:]
		var parts []string = strings.Split(paramsUrl, "/")
		var functionName string = parts[0]
		var timeout int
		if timeout, err = strconv.Atoi(parts[1]); err != nil {
			errorHandler(w, r, http.StatusBadRequest)
			return
		}
		if wsData.verbose && (functionName != "keepalive") && (functionName != "check_wait_status") {
			fmt.Printf("functionName \"%s\", timeout %d\n", functionName, timeout)
		}

		jResp, err = clientCallHandler(functionName, timeout, r.Body, wsData)

		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			errorHandler(w, r, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, jResp)
	}

	rawFileHandler := func(w http.ResponseWriter, r *http.Request) {
		var err error
		var urlPath = r.URL.Path[lenSecretKey+2:]
		if wsData.verbose {
			fmt.Printf(" server return response for %s\n", urlPath)
		}

		var mimeType string = getMimeType(urlPath)
		if len(mimeType) != 0 {
			w.Header().Set("Content-Type", mimeType)
		}

		if wsData.zipReader != nil {
			// return data from zip file
			var data []byte

			var filespec string = urlPath

			data, err = bundle.ReadFile(wsData.zipReader, filespec)
			if err != nil {
				fmt.Fprintf(os.Stderr, "    Unable to read file %s\n", filepath.FromSlash(filespec))
				errorHandler(w, r, http.StatusNotFound)
				return
			} else {
				if wsData.verbose {
					fmt.Printf("    Return file %s\n", filepath.FromSlash(filespec))
				}
				fmt.Fprintf(w, "%s", data)
			}

		} else {
			// return RAW data from filesystem

			var fullFilename string = path.Join(wsData.fullRootPath, urlPath)

			var isdir bool
			if isdir, err = ssutil.IsDir(fullFilename); err != nil {
				fmt.Fprintf(os.Stderr, "    Invalid URL request for file %s - %s\n", filepath.FromSlash(path.Join(rootPath, urlPath)), err.Error())
				errorHandler(w, r, http.StatusNotFound)
				return
			}

			if isdir {
				urlPath = path.Join(urlPath, wsData.initFile)
				fullFilename = path.Join(wsData.fullRootPath, urlPath)
			}

			filepathFromSlash := filepath.FromSlash(fullFilename)
			file, err := os.Open(filepathFromSlash)
			if err != nil {
				fmt.Fprintf(os.Stderr, "    Unable to read file %s\n", filepath.FromSlash(path.Join(rootPath, urlPath)))
				errorHandler(w, r, http.StatusNotFound)
				return
			}
			defer file.Close()

			fileTime := time.Time{}

			rangeStr := r.Header.Get("Range")
			if rangeStr == "" {
				fileStat, err := os.Stat(filepathFromSlash)
				if err == nil {
					fileTime = fileStat.ModTime()
				}
			}

			_, filenameOnly := path.Split(filepathFromSlash)
			if wsData.verbose {
				fmt.Printf("    Return file %s\n", filepath.FromSlash(path.Join(rootPath, urlPath)))
			}
			http.ServeContent(w, r, filenameOnly, fileTime, file)
		}
	}

	http.HandleFunc("/call/"+secretKey+"/", callHandler)
	http.HandleFunc("/"+secretKey+"/slfsrv-core.js", self_serving_core_js_handler)
	http.HandleFunc("/"+secretKey+"/", rawFileHandler)
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/", noSecretKeyHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}

func keep_aliver(keepAliveChan chan int, exitChan chan int, verbose bool, keepAliveSeconds int64) {
	// quit if not tickled within a few seconds

	maxwait := time.Second * time.Duration(keepAliveSeconds)
	for {
		select {
		case <-keepAliveChan:
		case <-time.After(maxwait):
			if verbose {
				fmt.Println("CONNECTION WITH CLIENT LOST!")
			}
			exitChan <- 1
		}
	}
}

func ListenAndServe(port int, secretKey string, keepAliveSeconds int64, zipReader *zip.Reader,
	rootPath string, fullRootPath string, initFile string,
	verbose bool, exitChan chan int, myselfExecutable string,
	storeFilespec string) int {
	if port == 0 {
		// search around for an available port
		for {
			port = 8000 + rand.Intn(1000) // 8000 -> 8999
			if verbose {
				fmt.Printf("Try port %d ", port)
			}
			psock, err := net.Listen("tcp", ":"+strconv.Itoa(port))
			if err == nil {
				psock.Close()
				if verbose {
					fmt.Printf(" OK\n")
				}
				break
			}
			fmt.Printf(" - not available\n")
		}
	}

	var keepAliveChan chan int = make(chan int)
	var wsData webserverData = webserverData{keepAliveChan: keepAliveChan, verbose: verbose,
		fullRootPath: fullRootPath, initFile: initFile, storeFilespec: storeFilespec,
		zipReader: zipReader}
	go web_server_forever(secretKey, &wsData, port, rootPath, myselfExecutable)

	go keep_aliver(keepAliveChan, exitChan, verbose, keepAliveSeconds)

	return port
}
