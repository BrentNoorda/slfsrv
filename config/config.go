package config

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/BrentNoorda/slfsrv/ssutil"
)

type PathDescription struct {
	FsPath      []string
	WebPath     []string
	Description string
	IsRelative	bool
}

type Settings struct {
	SecretKey string
	Port int
}
var webRoutes []PathDescription
var settings Settings

// call once to initialize
func ParseConfigFile() {

	type PathDescriptionJson struct {
		FsPath      string `json:"fsPath"`
		WebPath     string `json:"webPath"`
		Description string `json:"description"`
		Security    int    `json:"security"`

	}

	type ConfigJson struct {
		Routing  []PathDescriptionJson `json:"routing"`
		Foo				string	`json:"foo"`
		SecretKey string `json:"secretKey"`
		Port int `json:"port"`
	}

	// todo: let this be entered on command line
	content, err := ioutil.ReadFile("./config.json")
	if err == nil {
		var tmpConfig ConfigJson
		err := json.Unmarshal(content, &tmpConfig)
		if err != nil {
			fmt.Println(err)
		}
		settings.SecretKey = tmpConfig.SecretKey
		settings.Port = tmpConfig.Port

		webRoutes = make([]PathDescription, len(tmpConfig.Routing))

		for i := 0; i < len(tmpConfig.Routing); i++ {
			item := tmpConfig.Routing[i]
			webRoutes[i].FsPath,webRoutes[i].IsRelative = ssutil.PathStringToArray(item.FsPath)
			webRoutes[i].WebPath,_ = ssutil.PathStringToArray(item.WebPath)
			webRoutes[i].Description = item.Description
		}
		fmt.Printf("webRoutes: %+v\n\n", webRoutes)
		fmt.Printf("settings: %+v\n\n", settings)
	} else {
		panic(err)
	}
}

func GetSettings() *Settings {
	return &settings
}

/*
func GetRoutingTable() []PathDescription {
	return webRoutes
}
*/

// call when we have a web request.  Gets an array of all possible
// file system paths that could satisfy the request. This allows multiple
// file system directories to handle a single web directory....it will
// just serve the first file (ordered based on the config file) that
// is valid for that url
func GetPossibleFilePathsFromUrl(url string) []string {
	out := make([]string, 0, 3)

	urlPath,_ := ssutil.PathStringToArray(url)

	fmt.Printf("urlPath: %+v\n\n", urlPath)
	for i := 0; i < len(webRoutes); i++ {
		success := true
		webPath := webRoutes[i].WebPath
			fmt.Printf("webPath %d: %+v\n\n", i, webPath)

		lenWebPath := len(webPath)
		for j := 0; j < lenWebPath; j++ {
			if j >= len(urlPath) {
				break
			}
			if webPath[j] != urlPath[j] {
				fmt.Printf("fail %d: %s %s\n\n", i, webPath[j], urlPath[j])
				success = false
				break
			}
		}
		if success {
			fmt.Printf("success: %s\n\n", "/" +
					strings.Join(webRoutes[i].FsPath, "/") + "/" +
					strings.Join(urlPath[lenWebPath:], "/"))
			out = append(out, "/" +
					strings.Join(webRoutes[i].FsPath, "/") + "/" +
					strings.Join(urlPath[lenWebPath:], "/"))
		}
	}
	/*
	trimmedOut := make([]string, len(out))
	copy(trimmedOut, out)
	return trimmedOut
	*/
	return out
}



