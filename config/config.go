package config

import (
	"encoding/json"
	"fmt"
	"github.com/BrentNoorda/slfsrv/ssutil"
	"io/ioutil"
	"strings"
)

type PathDescription struct {
	FsPath      []string
	WebPath     []string
	Description string
	IsRelative  bool
}

type Settings struct {
	SecretKey        string
	Port             int
	KeepAliveSeconds int64
	SecretKeyInPath  bool
}

var webRoutes []PathDescription
var settings Settings

// call once to initialize
func ParseConfigFile(filename string) {

	type PathDescriptionJson struct {
		FsPath      string `json:"fsPath"`
		WebPath     string `json:"webPath"`
		Description string `json:"description"`
		Security    int    `json:"security"`
	}

	type ConfigJson struct {
		Routing          []PathDescriptionJson `json:"routing"`
		Port             int                   `json:"port"`
		KeepAliveSeconds int64                 `json:"keepAliveSeconds"`
		SecretKey        string                `json:"secretKey"`
		SecretKeyInPath  bool                  `json:"secretKeyInPath"`
	}

	if filename == "" {
		filename = "./config.json"
	}

	content, err := ioutil.ReadFile(filename)
	if err == nil {
		var tmpConfig ConfigJson
		err := json.Unmarshal(content, &tmpConfig)
		if err != nil {
			fmt.Println(err)
		}
		settings.Port = tmpConfig.Port
		settings.KeepAliveSeconds = tmpConfig.KeepAliveSeconds
		settings.SecretKey = tmpConfig.SecretKey
		settings.SecretKeyInPath = tmpConfig.SecretKeyInPath

		webRoutes = make([]PathDescription, len(tmpConfig.Routing))

		for i := 0; i < len(tmpConfig.Routing); i++ {
			item := tmpConfig.Routing[i]
			webRoutes[i].FsPath, webRoutes[i].IsRelative = ssutil.PathStringToArray(item.FsPath)
			webRoutes[i].WebPath, _ = ssutil.PathStringToArray(item.WebPath)
			webRoutes[i].Description = item.Description
		}
		// fmt.Printf("webRoutes: %+v\n\n", webRoutes)
		// fmt.Printf("settings: %+v\n\n", settings)
	}
}

func GetSettings() *Settings {
	return &settings
}

// call when we have a web request.  Gets an array of all possible
// file system paths that could satisfy the request. This allows multiple
// file system directories to handle a single web directory....it will
// just serve the first file (ordered based on the config file) that
// is valid for that url
func GetPossibleFilePathsFromUrl(url string) []string {
	out := make([]string, 0, 3)

	urlPath, _ := ssutil.PathStringToArray(url)

	for i := 0; i < len(webRoutes); i++ {
		success := true
		webPath := webRoutes[i].WebPath
		lenWebPath := len(webPath)
		for j := 0; j < lenWebPath; j++ {
			if j >= len(urlPath) {
				break
			}
			if webPath[j] != urlPath[j] {
				success = false
				break
			}
		}
		if success {
			if lenWebPath <= len(urlPath) {
				startPath := "/"
				if webRoutes[i].IsRelative {
					startPath = ""
				}
				out = append(out, startPath+
					strings.Join(webRoutes[i].FsPath, "/")+"/"+
					strings.Join(urlPath[lenWebPath:], "/"))
			}
		}
	}
	return out
}
