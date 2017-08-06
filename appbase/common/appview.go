package common

import (
	// b64 "encoding/base64"
	"fmt"
	"github.com/appbaseio/abc/log"
	"net/url"
	"strings"
)

// MakeDejavuURL ...
func MakeDejavuURL(appURL string) (string, error) {
	idx := strings.LastIndex(appURL, "/")
	hostURL := appURL[:idx]
	appName := appURL[idx+1:]
	jsonStr := fmt.Sprintf("{\"appname\":\"%s\",\"url\":\"%s\",\"selectedType\":[]}", appName, hostURL)
	// jsonStr = b64.StdEncoding.EncodeToString([]byte(jsonStr))
	jsonStr = url.QueryEscape(jsonStr)
	url := "https://opensource.appbase.io/dejavu/live/#?app=" + jsonStr
	log.Debugln(url)
	// base64 encode
	// https://base64-redirect.glitch.me/
	// encURI := "https://base64-redirect.glitch.me/redirect?to=" + b64.StdEncoding.EncodeToString([]byte(url))

	return url, nil
}

// MakeMirageURL ...
func MakeMirageURL(appURL string) (string, error) {
	idx := strings.LastIndex(appURL, "/")
	hostURL := appURL[:idx]
	appName := appURL[idx+1:]
	jsonStr := fmt.Sprintf("{\"appname\":\"%s\",\"url\":\"%s\",\"selectedType\":[]}", appName, hostURL)
	jsonStr = url.QueryEscape(jsonStr)
	url := "https://opensource.appbase.io/mirage/#?app=" + jsonStr
	log.Debugln(url)

	return url, nil
}
