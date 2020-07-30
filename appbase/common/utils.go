package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/appbaseio/abc/log"
)

// GetKeyForValue returns key for the given value
func GetKeyForValue(data map[string]string, val string) string {
	for k, v := range data {
		if v == val {
			return k
		}
	}
	return ""
}

// JSONNumberToString converts a json.Number to string, properly
// i.e. no decimal points for a integer
// json.Number is required instead of normal types in map[..].. based decoding
func JSONNumberToString(number json.Number) string {
	str := number.String()
	if strings.HasSuffix(str, ".0") {
		return str[0 : len(str)-2]
	}
	return str
}

// JSONNumberToInt ...
func JSONNumberToInt(number json.Number) int64 {
	f, err := number.Float64()
	if err != nil {
		log.Errorln(err)
		return 0
	}
	return int64(f)
}

// StringInSlice checks if string is in list or not
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ColonPad pads spaces after colon
func ColonPad(text string, length int) string {
	// remove brackets in names, they are long
	bracket := strings.Index(text, "(")
	if bracket > -1 {
		text = text[:bracket]
	}
	textLen := len(text)
	text += ":"
	for i := 0; i < (length - textLen - 1); i++ {
		text += " "
	}
	return text
}

// OpenURL opens the specified URL in the default browser of the user.
// https://stackoverflow.com/a/39324149/2295672
func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// SizeInKB shows size in KB
func SizeInKB(size int) int {
	return size / 1024 // original size in bytes
}

// IsFileValid check if the file is valid
func IsFileValid(file string) error {

	if _, err := url.ParseRequestURI(file); err == nil { //do not check remote file validity here
		log.Infoln("Importing data from remote file", file)
		return nil
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return err
	}
	return nil
}

// RemoveDuplicates removes duplicate values in a slice
// https://groups.google.com/forum/#!topic/golang-nuts/-pqkICuokio
func RemoveDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

// Max function
func Max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

// DefaultDownloadDirectory to download remote files
var DefaultDownloadDirectory = "./temp"

//DownloadFile create a local copy of the remote json file
func DownloadFile(filepath string, url string) error {

	log.Infoln("Downloading the file from remote URL:", url)
	_, err := os.Stat(DefaultDownloadDirectory) //check if  directory exist
	if os.IsNotExist(err) {

		err = os.Mkdir(DefaultDownloadDirectory, 0755) //if not create, directory
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// RemoveFile delete remote file
func RemoveFile(fileName string) error {
	log.Infoln("Deleting the temporary file:", fileName)
	if err := os.Remove(fileName); err != nil {
		log.Debugf(fmt.Sprintf("Unable to delete temporary file: %s", fileName), err)
		return err
	}
	log.Infoln("file successfully deleted at path:", fileName)
	return nil
}
