package browser

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func LaunchDefaultBrowser(port int, secretKey string, initPathUrl string, initQuery string, verbose bool) {
	var runSpec string = "http://127.0.0.1:" + strconv.Itoa(port) + "/" + secretKey + "/" + initPathUrl
	if initQuery != "" {
		runSpec += "?" + initQuery
	}
	if verbose {
		fmt.Println("xdg-open", runSpec)
	}
	cmd := exec.Command("xdg-open", runSpec)
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
