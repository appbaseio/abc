package login

import (
	"fmt"
	"github.com/aviaryan/abc/appbase/common"
	"github.com/aviaryan/abc/appbase/session"
	"github.com/aviaryan/abc/appbase/user"
	"os/exec"
	"runtime"
)

// IsUserAuthenticated checks if user is logged in or not
func IsUserAuthenticated() bool {
	data, err := session.LoadUserSessionAsString()
	return err == nil && data != ""
}

// StartUserLogin starts user login process
func StartUserLogin(host string) error {
	url := fmt.Sprintf("%s/login/%s?next=%s/user/token", common.AccAPIURL, host, common.AccAPIURL)
	fmt.Printf("Opening %s in the browser.\n", url)
	fmt.Println("Once authenticated, copy the token from there and paste it into terminal.")
	// open in browser
	err := open(url)
	if err != nil {
		return err
	}
	// read input
	fmt.Print("> ")
	var token string
	fmt.Scanf("%s", &token)
	// save to file
	err = session.SaveUserSession(token)
	if err != nil {
		return err
	}
	// show email
	email, err := user.GetUserEmail()
	if err == nil {
		fmt.Printf("\nLogged in as %s\n", email)
	} else {
		fmt.Println("\nFailed to get user info. Please try again.")
	}
	return err
}

// https://stackoverflow.com/a/39324149/2295672
// open opens the specified URL in the default browser of the user.
func open(url string) error {
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
