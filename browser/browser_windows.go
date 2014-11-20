package browser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func LaunchDefaultBrowser(port int, secretKey string, initPathUrl string, initQuery string, verbose bool) {
	var runSpec string = "http://127.0.0.1:" + strconv.Itoa(port) + "/" + secretKey + "/" + filepath.ToSlash(initPathUrl)
	if initQuery != "" {
		runSpec += "?" + initQuery
	}
	if verbose {
		fmt.Println("cmd.exe /C start", runSpec)
	}
	cmd := exec.Command("cmd.exe", "/C", "start", runSpec)
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
