package tempdir

import (
    "os"
    "io"
    "fmt"
    "sync"
    "strings"
    "io/ioutil"
    "archive/zip"
    "path/filepath"
    "github.com/BrentNoorda/slfsrv/ssutil"
)


var tempdirMutex = &sync.Mutex{}  // lock calls to create new tempdir to make them synchronous

var tempdirs map[string]string = make(map[string]string)
    // all known tempdirs created, so we can return existing
    // ones, and so we can cleanup at the end


func writeToTempDir(src io.Reader,tempDir string,relFileName string,perm os.FileMode) (err error) {
    // write data to new file, creating and directories as needed
    var newDirname string
    newDirname,_ = filepath.Split(relFileName)
    if ( len(newDirname) != 0 ) {
        if err = os.MkdirAll(filepath.Join(tempDir,newDirname),0777); err != nil { return }
    }

    var dstF *os.File

    dstF,err = os.OpenFile(filepath.Join(tempDir,relFileName),os.O_CREATE | os.O_WRONLY | os.O_TRUNC,perm)
    if err == nil {
        defer dstF.Close()

        _,err = io.Copy(dstF,src)
    }

    return
}

func Tempdir(prefix string,zipReader *zip.Reader/*null if reading from raw files*/,
             rootPath string/*used if zipRead is nul*/,verbose bool) (tempdir string,err error) {

    tempdirMutex.Lock()
    defer tempdirMutex.Unlock()

    var foundInMap bool

    if tempdir,foundInMap = tempdirs[prefix]; foundInMap {
        // this prefix is already mapped, so nothing to do but return the existing mapping
        return;
    }

    if tempdir,err = ioutil.TempDir( "", "slfsrv-tmp-" + prefix ); err != nil { return }
    if verbose {
        fmt.Println("Tempdir \"" + tempdir + "\" created for prefix \"" + prefix + "\"")
    }
    tempdirs[prefix] = tempdir

    if len(prefix) != 0 {
        // copy all the local files in our prefix directory to the temp dir
        if zipReader != nil {
            // read from the zip file

            var startsWith string = prefix + "/"  // if if windows we store this slash type
            var startsWithLength int = len(startsWith)

            var f *zip.File
            for _, f = range zipReader.File {

                if strings.HasPrefix(f.Name,startsWith) {
                    var rc io.ReadCloser

                    if rc,err = f.Open(); err != nil { return }
                    err = writeToTempDir(rc,tempdir,filepath.FromSlash(f.Name[startsWithLength:]),0777/*f.Mode()*/)
                    if err != nil {
                        rc.Close()
                    } else {
                        err = rc.Close()
                    }
                }
            }

        } else {
            // read from the raw file system

            var isDir bool
            var copyFromDir string = filepath.Join(rootPath,prefix)

            if isDir,err = ssutil.IsDir(copyFromDir); (err == nil) && isDir {
                // prefix is a directory, so copy everything from it to tempdir

                var skipLen int = len(copyFromDir) + 1
                walkFunc := func(path string, info os.FileInfo, err error) error {
                    if err == nil {
                        if !info.IsDir() {
                            var relativePath string = path[skipLen:]

                            var f *os.File

                            if f,err = os.Open(path); err == nil {
                                err = writeToTempDir(f,tempdir,relativePath,info.Mode())
                                if err != nil {
                                    f.Close()
                                } else {
                                    err = f.Close()
                                }
                            }

                        }
                    }
                    return err
                }
                if err = filepath.Walk(copyFromDir,walkFunc); err != nil { return }
            }
        }
    }

    return
}







func Cleanup(verbose bool) {
    // remove all tempdirs
    var err error
    var dir,prefix string

    for prefix,dir = range tempdirs {
        if verbose {
            fmt.Println("delete tempdir \"" + dir + "\" for prefix \"" + prefix + "\"");
        }
        if err = os.RemoveAll(dir); err != nil {
            fmt.Fprintf(os.Stderr,"ERROR: Unable to remove tempdir directory \"" + dir + "\"")
        }
    }
}