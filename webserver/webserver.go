package webserver

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/BrentNoorda/slfsrv/bundle"
	"github.com/BrentNoorda/slfsrv/config"
	"math/rand"
	"net"
	"path"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"container/list"
	"text/template"
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

func isDirectory(path string) bool {
    fileInfo, err := os.Stat(path)
    if err != nil {
    	return false
    }
    return fileInfo.IsDir()
}

func getMimeType(filespec string) (mimeType string) { // return "" if not known
	var lastDot int = strings.LastIndex(filespec, ".")
	if -1 != lastDot {
		var fileExtension = strings.ToLower(filespec[lastDot+1 : len(filespec)])
		mimeType = mimeTypes[fileExtension]
	}
	return
}

func web_server_forever(secretKey string, wsData *webserverData,
	port int, rootPath string, myselfExecutable string) {
	// rjb var lenSecretKey int = len(secretKey)

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
		//  rjb      s = strings.Replace(s,"<<<SECRET_KEY>>>",secretKey,1)
		s = strings.Replace(s, "<<<SECRET_KEY>>>", "", 1)
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
		var paramsUrl = r.URL.Path // lenSecretKey+7:] rjb
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

	serveFileIfExists := func(filePath string, webPath string, w http.ResponseWriter, r *http.Request) bool {
		name := filepath.FromSlash(filePath)
		fmt.Printf("name: %+v\n\n", name)
		dirNameOnly := ""
		if isDirectory(name) {
			fmt.Println("---------------- is dir [" + name + "]")
			if(strings.LastIndex(webPath, "/") != (len(webPath)-1)) {
				http.Redirect(w, r, webPath+"/", http.StatusFound)
        return true
			}
			dirNameOnly = name
			name = name + "index.html"
			fmt.Println("dir 1 " + dirNameOnly + "[" + name + "]")
		}

		file, err := os.Open(name)
		if err != nil {
			if(len(dirNameOnly)>0) {

				// Opening the file handle
				dir, err := os.Open(dirNameOnly)
				if err != nil {
					http.Error(w, "404 Not Found : Error while opening the file.", 404)
					return false
				}

				defer dir.Close()

				handleDirectory(dir, w, r)
				fmt.Println("dir 2 " + dirNameOnly + "[" + name + "]")
				return true
			}
			return false
		}
		defer file.Close()
		_, filenameOnly := path.Split(name)
		http.ServeContent(w, r, filenameOnly, time.Time{}, file)
		return true
	}

	rawFileHandler := func(w http.ResponseWriter, r *http.Request) {
		var err error
		var urlPath = r.URL.Path // rjb[lenSecretKey+2:]

		if wsData.verbose {
			fmt.Printf(" -- server return response for %s\n", urlPath)
		}
		/* rjb
		        // need the cookie?
		        cookie, err := r.Cookie("cookiename")
		        if(err == nil) {
						  fmt.Printf(" -- cookie cookiename: %+v\n", cookie.Value)
						}
		*/

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
					fmt.Printf("    Return zipped file %s\n", filepath.FromSlash(filespec))
				}
				fmt.Fprintf(w, "%s", data)
			}

		} else {
			// return RAW data from filesystem
			paths := config.GetPossibleFilePathsFromUrl(urlPath)
			success := false
			for i := 0; i < len(paths); i++ {
				if serveFileIfExists(paths[i], urlPath, w, r) == true {
					success = true
					break
				}
			}
			if !success {
				fmt.Fprintf(os.Stderr, "    Unable to read file %+v\n", paths)
				errorHandler(w, r, http.StatusNotFound)
			}
		}
	}

	// rjb   http.HandleFunc("/call/"+secretKey+"/", callHandler)
	//    http.HandleFunc("/"+secretKey+"/slfsrv-core.js", self_serving_core_js_handler)
	//   http.HandleFunc("/"+secretKey+"/", rawFileHandler)

	http.HandleFunc("/call/", callHandler)
	http.HandleFunc("/slfsrv-core.js", self_serving_core_js_handler)
	http.HandleFunc("/", rawFileHandler)

	http.HandleFunc("/favicon.ico", favicon)

	// rjb http.HandleFunc("/", noSecretKeyHandler)
	_ = noSecretKeyHandler
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}

func keep_aliver(keepAliveChan chan int, exitChan chan int) {
	// quit if not tickled within a few seconds

	// rjb:  temp  -- make sure it doesn't terminate
	var maxwait time.Duration = time.Second * 10000 // first call can take longer than others
	for {
		select {
		case <-keepAliveChan:
		case <-time.After(maxwait):
			exitChan <- 1
		}
		maxwait = time.Second * 2000
	}
}

func ListenAndServe(port int, secretKey string, zipReader *zip.Reader,
	rootPath string, fullRootPath string, initFile string,
	verbose bool, exitChan chan int, myselfExecutable string,
	storeFilespec string) int {


	settings := config.GetSettings()
	if settings.Port != 0 {
		port = settings.Port
	}

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

	go keep_aliver(keepAliveChan, exitChan)

	return port
}


// everything below is copied from https://gist.github.com/alexisrobert/982674
// and does a directory listing if there is no index.html file
type dirListing struct {
	Name           string
	childrenDir   []string
	childrenFiles []string
}

func copyToArray(src *list.List) []string {
	dst := make([]string, src.Len())

	i := 0
	for e := src.Front(); e != nil; e = e.Next() {
		dst[i] = e.Value.(string)
		i = i + 1
	}

	return dst
}

func handleDirectory(f *os.File, w http.ResponseWriter, req *http.Request) {
	names, _ := f.Readdir(-1)

	// Otherwise, generate folder content.
	childrenDirTmp := list.New()
	childrenFilesTmp := list.New()

	for _, val := range names {
		if val.Name()[0] == '.' {
			continue
		} // Remove hidden files from listing

		if val.IsDir() {
			childrenDirTmp.PushBack(val.Name())
		} else {
			childrenFilesTmp.PushBack(val.Name())
		}
	}

	// And transfer the content to the final array structure
	childrenDir := copyToArray(childrenDirTmp)
	childrenFiles := copyToArray(childrenFilesTmp)

	tpl, err := template.New("tpl").Parse(dirListingTemplate)
	if err != nil {
		http.Error(w, "500 Internal Error : Error while generating directory listing.", 500)
		fmt.Println(err)
		return
	}

	err = tpl.Execute(w, dirListing{Name: req.URL.Path,
		childrenDir: childrenDir, childrenFiles: childrenFiles})
	if err != nil {
		fmt.Println(err)
	}
}

const dirListingTemplate = `<?xml version="1.0" encoding="iso-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en">
<!-- Modified from lighttpd directory listing -->
<head>
<title>Index of {{.Name}}</title>
<style type="text/css">
a, a:active {text-decoration: none; color: blue;}
a:visited {color: #48468F;}
a:hover, a:focus {text-decoration: underline; color: red;}
body {background-color: #F5F5F5;}
h2 {margin-bottom: 12px;}
table {margin-left: 12px;}
th, td { font: 90% monospace; text-align: left;}
th { font-weight: bold; padding-right: 14px; padding-bottom: 3px;}
td {padding-right: 14px;}
td.s, th.s {text-align: right;}
div.list { background-color: white; border-top: 1px solid #646464; border-bottom: 1px solid #646464; padding-top: 10px; padding-bottom: 14px;}
div.foot { font: 90% monospace; color: #787878; padding-top: 4px;}
</style>
</head>
<body>
<h2>Index of {{.Name}}</h2>
<div class="list">
<table summary="Directory Listing" cellpadding="0" cellspacing="0">
<thead><tr><th class="n">Name</th><th class="t">Type</th><th class="dl">Options</th></tr></thead>
<tbody>
<tr><td class="n"><a href="../">Parent Directory</a>/</td><td class="t">Directory</td><td class="dl"></td></tr>
{{range .childrenDir}}
<tr><td class="n"><a href="{{.}}/">{{.}}/</a></td><td class="t">Directory</td><td class="dl"></td></tr>
{{end}}
{{range .childrenFiles}}
<tr><td class="n"><a href="{{.}}">{{.}}</a></td><td class="t">&nbsp;</td><td class="dl"><a href="{{.}}?dl">Download</a></td></tr>
{{end}}
</tbody>
</table>
</div>
</body>
</html>`

