package browser

import (
    "os"
    "fmt"
    "strconv"
    "os/exec"
)

func LaunchDefaultBrowser(port int,secretKey string,initPathUrl string,initQuery string,verbose bool) {
    var runSpec string = "http://127.0.0.1:" + strconv.Itoa(port) + "/" + secretKey + "/" + initPathUrl
    if initQuery != "" {
        runSpec += "?" + initQuery
    }
    if verbose {
        fmt.Println("open",runSpec)
    }
    cmd := exec.Command("open",runSpec)
    err := cmd.Start()
    if err != nil {
        fmt.Fprintf(os.Stderr,"%s\n",err)
        os.Exit(1)
    }
}
