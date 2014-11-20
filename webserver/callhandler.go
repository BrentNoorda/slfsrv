package webserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BrentNoorda/slfsrv/ssutil"
	"github.com/BrentNoorda/slfsrv/tempdir"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

////////////////// keep alive //////////////////

type keepaliveIn struct {
	Animal         string
	FavoriteNumber int
}

func (ka *keepaliveIn) goServe(lc *longCall, wsData *webserverData) {
	wsData.keepAliveChan <- 1
	lc.Response <- jsonResponse{"state": "done", "message": "Hello!"}
}

////////////////// tempdir //////////////////

type tempdirIn struct {
	Dir string
}

func (td *tempdirIn) goServe(lc *longCall, wsData *webserverData) {
	var dir string
	var err error

	dir, err = tempdir.Tempdir(td.Dir, wsData.zipReader, wsData.fullRootPath, wsData.verbose)
	if err != nil {
		lc.Response <- jsonResponseError("error: " + err.Error())
	} else {
		lc.Response <- jsonResponse{"state": "done", "dir": dir}
	}
}

////////////////// cwd //////////////////

type cwdGetIn struct{}

func (cwd *cwdGetIn) goServe(lc *longCall, wsData *webserverData) {
	dir, err := os.Getwd()
	if err != nil {
		lc.Response <- jsonResponseError("error: " + err.Error())
	} else {
		lc.Response <- jsonResponse{"state": "done", "cwd": dir}
	}
}

type cwdSetIn struct {
	Dir string
}

func (cwd *cwdSetIn) goServe(lc *longCall, wsData *webserverData) {
	var err error = os.Chdir(cwd.Dir)
	if err != nil {
		lc.Response <- jsonResponseError("error: " + err.Error())
	} else {
		lc.Response <- jsonResponse{"state": "done"}
	}
}

////////////////// dir //////////////////

type dirListIn struct {
	Dir string
}

func (dirlist *dirListIn) goServe(lc *longCall, wsData *webserverData) {

	var fis []os.FileInfo
	var err error

	if dirlist.Dir == "" {
		dirlist.Dir = "."
	}

	if fis, err = ioutil.ReadDir(dirlist.Dir); err != nil {
		lc.Response <- jsonResponseError("error: " + err.Error())
	} else {
		var fi os.FileInfo
		var list []map[string]interface{} = make([]map[string]interface{}, 0)
		for _, fi = range fis {
			list = append(list, jsonResponse{
				"name": fi.Name(),
				"size": fi.Size(),
				"dir":  fi.IsDir(),
				"perm": fi.Mode() & 0777,
				"mod":  fi.ModTime().UnixNano() / 1000000,
			})
		}
		lc.Response <- jsonResponse{"state": "done", "list": list}
	}
}

type dirExistsIn struct {
	Dir string
}

func (de *dirExistsIn) goServe(lc *longCall, wsData *webserverData) {
	var err error
	var exists, isdir bool

	if exists, err = ssutil.FileOrDirExists(de.Dir); err == nil {
		isdir, err = ssutil.IsDir(de.Dir)
	}
	if err != nil {
		lc.Response <- jsonResponseError("error: " + err.Error())
	} else {
		lc.Response <- jsonResponse{"state": "done", "exists": exists && isdir}
	}
}

////////////////// env //////////////////

type envGetIn struct {
	Key string
}

func (ge *envGetIn) goServe(lc *longCall, wsData *webserverData) {
	var value string = os.Getenv(ge.Key)
	lc.Response <- jsonResponse{"state": "done", "value": value}
}

type envSetIn struct {
	Key   string
	Value string
}

func (se *envSetIn) goServe(lc *longCall, wsData *webserverData) {
	var err error = os.Setenv(se.Key, se.Value)
	if err != nil {
		lc.Response <- jsonResponseError("error: " + err.Error())
	} else {
		lc.Response <- jsonResponse{"state": "done"}
	}
}

////////////////// file //////////////////

type fileReadIn struct {
	File string
	Mode string
}

func (fr *fileReadIn) goServe(lc *longCall, wsData *webserverData) {
	var contents []byte
	var err error
	var binaryMode bool

	contents, err = ioutil.ReadFile(fr.File)
	if err != nil {
		lc.Response <- jsonResponseError("error: " + err.Error())
	} else {
		binaryMode = (fr.Mode == "binary")
		if binaryMode {
			var buf bytes.Buffer
			var c byte
			for _, c = range contents {
				if _, err = buf.WriteString(string(c)); err != nil {
					lc.Response <- jsonResponseError("error: " + err.Error())
					return
				}
			}
			lc.Response <- jsonResponse{"state": "done", "contents": buf.String()}
		} else {
			// text mode
			lc.Response <- jsonResponse{"state": "done", "contents": string(contents)}
		}
	}
}

type fileWriteIn struct {
	File     string
	Mode     string
	Contents string
	Bytes    []byte
}

func (fw *fileWriteIn) goServe(lc *longCall, wsData *webserverData) {
	var err error
	var binaryMode, appendMode bool
	var file *os.File
	var flag int

	binaryMode = strings.Contains(fw.Mode, "binary")
	appendMode = strings.Contains(fw.Mode, "append")

	if appendMode {
		flag = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	} else {
		flag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	if file, err = os.OpenFile(fw.File, flag, 0600); err == nil {

		defer func() {
			file.Close()
			// should file be removed? only on append?
		}()

		if binaryMode {

			var b [1]byte
			for i := 0; i < len(fw.Bytes); i++ {
				b[0] = fw.Bytes[i]
				_, err = file.Write(b[:])
				if err != nil {
					break
				}
			}
		} else {
			_, err = file.WriteString(fw.Contents)
		}
	}
	if err != nil {
		lc.Response <- jsonResponseError("error: " + err.Error())
	} else {
		lc.Response <- jsonResponse{"state": "done"}
	}
}

type fileExistsIn struct {
	File string
}

func (fe *fileExistsIn) goServe(lc *longCall, wsData *webserverData) {
	var err error
	var exists, isdir bool

	if exists, err = ssutil.FileOrDirExists(fe.File); err == nil {
		isdir, err = ssutil.IsDir(fe.File)
	}
	if err != nil {
		lc.Response <- jsonResponseError("error: " + err.Error())
	} else {
		lc.Response <- jsonResponse{"state": "done", "exists": exists && !isdir}
	}
}

type fileRemoveIn struct {
	File string
}

func (fr *fileRemoveIn) goServe(lc *longCall, wsData *webserverData) {
	var isdir bool
	var err error

	if isdir, err = ssutil.IsDir(fr.File); err != nil {
		lc.Response <- jsonResponseError("error: " + err.Error())
	}

	if isdir {
		lc.Response <- jsonResponseError("error: file.remove() cannot remove a directory")
	} else {
		err = os.Remove(fr.File)
		if err != nil {
			lc.Response <- jsonResponseError("error: " + err.Error())
		} else {
			lc.Response <- jsonResponse{"state": "done"}
		}
	}
}

////////////////// store //////////////////

var storeMutex = &sync.Mutex{}         // will be used so only one store call is working at a time
var store map[string]interface{} = nil // will maintain in-memory the full json file,
// maybe not great for memory but so what

func delete_store(wsData *webserverData) {
	// delete any existing store file and data
	store = make(map[string]interface{})

	var err error = os.Remove(wsData.storeFilespec)
	if err != nil {
		var pErr *os.PathError = err.(*os.PathError)
		if pErr.Err.Error() != "no such file or directory" {
			fmt.Fprintf(os.Stderr, "Unable to delete data file \"%s\". Error: %s\n", wsData.storeFilespec, err)
		}
	}
}

func store_read(wsData *webserverData) {
	// initialize store if not already initialized, if there was any corruption issue
	// then delete the whole file and start over (via delete_store())

	if store == nil {
		// following code from http://stackoverflow.com/questions/16681003/how-do-i-parse-a-json-file-into-a-struct-with-go
		var file *os.File
		var err error

		if file, err = os.Open(wsData.storeFilespec); err != nil {
			delete_store(wsData)
		} else {
			defer func() {
				if file != nil {
					file.Close()
				}
			}()
			// see http://blog.golang.org/json-and-go for arbitrary json
			var allStore interface{}
			var jsonParser *json.Decoder = json.NewDecoder(file)
			if err = jsonParser.Decode(&allStore); err != nil {
				file.Close()
				file = nil
				delete_store(wsData)
			} else {
				store = allStore.(map[string]interface{})
			}
		}
	}
}

func store_write(wsData *webserverData, maybeEmpty bool) (err error) {
	// if there is an error, will not call delete_store() (but might if map is empty)

	if maybeEmpty && (len(store) == 0) {
		delete_store(wsData)
	} else {
		var data []byte

		if data, err = json.MarshalIndent(&store, "", "\t"); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to encode store to file \"%s\". Error: %s\n", wsData.storeFilespec, err)
		} else {
			err = ioutil.WriteFile(wsData.storeFilespec, data, 0666)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to write store file \"%s\". Error: %s\n", wsData.storeFilespec, err)
			}
		}
	}
	return
}

type storeGetIn struct {
	Key string
}

func (s *storeGetIn) goServe(lc *longCall, wsData *webserverData) {
	storeMutex.Lock()
	defer storeMutex.Unlock()
	store_read(wsData)

	var ok bool
	var value interface{}
	value, ok = store[s.Key]
	if ok {
		lc.Response <- jsonResponse{"state": "done", "value": value}
	} else {
		lc.Response <- jsonResponse{"state": "done"}
	}
}

type storeSetIn struct {
	Key   string
	Value *interface{}
}

func (s *storeSetIn) goServe(lc *longCall, wsData *webserverData) {
	storeMutex.Lock()
	defer storeMutex.Unlock()
	store_read(wsData)

	var err error
	if s.Value == nil {
		delete(store, s.Key)
		err = store_write(wsData, true)
	} else {
		store[s.Key] = *s.Value
		err = store_write(wsData, false)
	}

	if err != nil {
		lc.Response <- jsonResponseError("Unable to write data: " + err.Error())
	} else {
		lc.Response <- jsonResponse{"state": "done"}
	}
}

type storeListIn struct{}

func (s *storeListIn) goServe(lc *longCall, wsData *webserverData) {
	storeMutex.Lock()
	defer storeMutex.Unlock()
	store_read(wsData)

	var keys []string = make([]string, 0)
	var key string
	for key = range store {
		keys = append(keys, key)
	}

	lc.Response <- jsonResponse{"state": "done", "keys": keys}
}

////////////////// exec //////////////////

type execIn struct {
	Program string
	Args    []string
	Input   string
}

func (e *execIn) goServe(lc *longCall, wsData *webserverData) {

	var cmd *exec.Cmd
	var err error

	cmd = exec.Command(e.Program)
	cmd.Args = e.Args
	if len(e.Input) != 0 {
		cmd.Stdin = strings.NewReader(e.Input)
	}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err = cmd.Start(); err != nil {
		lc.Response <- jsonResponseError("Unable to start: " + e.Program)
		return
	}

	var exitStatus int = 0 // assume success unless error
	if err = cmd.Wait(); err != nil {

		// exit status from http://stackoverflow.com/questions/10385551/get-exit-code-go
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitStatus = status.ExitStatus()
			} else {
				lc.Response <- jsonResponseError("Unable to to determine exit status")
				return
			}
		} else {
			lc.Response <- jsonResponseError("Error with program: " + err.Error())
			return
		}
	}

	lc.Response <- jsonResponse{
		"state":      "done",
		"exitStatus": exitStatus,
		"stdout":     stdout.String(),
		"stderr":     stderr.String(),
	}
}

////////////////// infrastructure //////////////////

type longCall struct {
	Die      bool              // send to this channel if the coroutine should abort prematurely
	Response chan jsonResponse // if success, the goroutine puts response here
	Timeout  int               // not really time, but how many times to wait a second, which is close
}

var callsWaiting map[int]*longCall = make(map[int]*longCall)
var callsWaitingMutex = &sync.Mutex{}

type ajaxInputInterface interface {

	/* goServe has two purposes. Primarily (hopefully) it is to send a jsonResponse to lc.jsonResponse
	   (i.e. lc.Response <- jsonResponse{"state":"done", "message": "Hello!"}). Secondarily, if it is
	   long running, it should periodically check lc.Die and abort if it's true.
	*/
	goServe(lc *longCall, wsData *webserverData)
}

func jsonResponseError(msg string) jsonResponse {
	return jsonResponse{"state": "error", "message": msg}
}

func addLongCall(timeout int) (key int, lc *longCall) {
	var newLongCall longCall = longCall{Die: false, Response: make(chan jsonResponse), Timeout: timeout}
	lc = &newLongCall
	callsWaitingMutex.Lock()
	defer callsWaitingMutex.Unlock()
	for {
		key = 1 + rand.Intn(1000000) // 1 -> 1000000
		if _, ok := callsWaiting[key]; !ok {
			// this random key not already in the map, so use this slot
			callsWaiting[key] = lc
			return
		}
	}
}

func getLongCall(key int) (lc *longCall) {
	callsWaitingMutex.Lock()
	lc = callsWaiting[key]
	callsWaitingMutex.Unlock()
	return
}

func removeLongCall(key int, abort bool) {
	callsWaitingMutex.Lock()
	var lc *longCall = callsWaiting[key]
	delete(callsWaiting, key)
	callsWaitingMutex.Unlock()

	if abort {
		lc.Die = true
	}
}

type waitForStatusIn struct {
	Key int
}

func clientCallHandler(functionName string, timeout int, body io.ReadCloser, wsData *webserverData) (jResp jsonResponse, err error) {

	var key int
	var lc *longCall
	var in ajaxInputInterface
	var decoder *json.Decoder = json.NewDecoder(body)

	if functionName == "check_wait_status" {
		var waitIn waitForStatusIn
		if err = decoder.Decode(&waitIn); err != nil {
			return
		}
		key = waitIn.Key
		lc = getLongCall(key)
	} else {

		switch functionName {

		case "keepalive":
			var ka keepaliveIn
			in = &ka

		case "exec":
			var exec execIn
			in = &exec

		case "store.get":
			var storeget storeGetIn
			in = &storeget
		case "store.set":
			var storeset storeSetIn
			in = &storeset
		case "store.list":
			var storelist storeListIn
			in = &storelist

		case "file.read":
			var fileread fileReadIn
			in = &fileread
		case "file.write":
			var filewrite fileWriteIn
			in = &filewrite
		case "file.exists":
			var fileexists fileExistsIn
			in = &fileexists
		case "file.remove":
			var fileremove fileRemoveIn
			in = &fileremove

		case "cwd.get":
			var cwdget cwdGetIn
			in = &cwdget
		case "cwd.set":
			var cwdset cwdSetIn
			in = &cwdset
		case "dir.list":
			var dirlist dirListIn
			in = &dirlist
		case "dir.exists":
			var direxists dirExistsIn
			in = &direxists

		case "env.get":
			var envget envGetIn
			in = &envget
		case "env.set":
			var setenv envSetIn
			in = &setenv

		case "tempdir":
			var tempdir tempdirIn
			in = &tempdir

		default:
			err = errors.New("Unknown function name: " + functionName)
			removeLongCall(key, false)
			return
		}

		if err = decoder.Decode(&in); err != nil {
			return
		}

		key, lc = addLongCall(timeout)

		go in.goServe(lc, wsData)
	}

	// wait up to a second for the call to finish (i.e. for response in the channel)
	select {
	case jResp = <-lc.Response:
		removeLongCall(key, false)
	case <-time.After(time.Second * 1):
		if lc.Timeout < 1 {
			// tired of waiting for this call to ever finish
			removeLongCall(key, true)
			jResp = jsonResponse{"state": "timeout", "message": "call timed out"}
		} else {
			lc.Timeout--
			jResp = jsonResponse{"state": "wait", "wait_key": key, "message": "call not finished yet"}
		}
	}

	if wsData.verbose && (functionName != "keepalive") && (functionName != "check_wait_status") {
		fmt.Println(jResp)
	}

	return
}
