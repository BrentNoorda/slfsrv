package ssutil

import(
    "os"
    "path"
    "errors"
    "strings"
    "runtime"
	"path/filepath"
)

func FileOrDirExists(path string) (bool,error) {
    // from http://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-denoted-by-a-path-exists-in-golang
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func IsAbs(filespec string) bool {  // because the default one seems screwed up
    var nativePath string = filepath.FromSlash(filespec)
    if path.IsAbs(nativePath) {
        return true
    }
    if runtime.GOOS == "windows" {
        // bug in windows
        if nativePath[1:3] == ":\\" {
            return true;
        }
    }
    return false
}

func IsDir(name string) (bool,error) {  // else assume is directory
    // from http://stackoverflow.com/questions/8824571/golang-determining-whether-file-points-to-file-or-directory
    var err error = nil
    f, err := os.Open(filepath.FromSlash(name))
    if err != nil {
        //return false,err
        return false,nil
    }
    defer f.Close()
    fi, err := f.Stat()
    if err != nil {
        return false,err
    }
    switch mode := fi.Mode(); {
    case mode.IsDir():
        // do directory stuffdd
        return true,err
    case mode.IsRegular():
        // do file stuff
        return false,err
    }
    return false,errors.New("Unrecognized directory")
}

func CleanPathRelativeToCwd(pathspec string,cwd string) string {
    // if path can be written better relative to CWD, then do so
    if 0 == strings.Index(pathspec,cwd) {
        if cwd == pathspec {
            pathspec = ""
        } else if 0 == strings.Index(pathspec,cwd+"/") {
            pathspec = pathspec[len(cwd)+1:]
        }
    }
    return pathspec
}