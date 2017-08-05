package common

import (
	"fmt"
	"github.com/appbaseio/abc/log"
	"strings"
)

// MakeDejavuURL ...
func MakeDejavuURL(appURL string) (string, error) {
	idx := strings.LastIndex(appURL, "/")
	hostURL := appURL[:idx]
	appName := appURL[idx+1:]
	jsonStr := fmt.Sprintf("{\"appname\":\"%s\",\"url\":\"%s\",\"selectedType\":[]}", appName, hostURL)
	url := "https://opensource.appbase.io/dejavu/live/#?app=" + jsonStr
	log.Debugln(url)
	return url, nil
}
