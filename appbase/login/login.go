package login

import (
	"fmt"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
	"github.com/appbaseio/abc/appbase/user"
	"github.com/appbaseio/abc/log"
	"os"
)

// IsUserAuthenticated checks if user is logged in or not
func IsUserAuthenticated() bool {
	data, err := session.LoadUserSessionAsString()
	if err == nil && data != "" {
		return true
	}
	// env token
	token := os.Getenv("ABC_TOKEN")
	return len(token) > 0
}

// StartUserLogin starts user login process
func StartUserLogin(host string) error {
	url := fmt.Sprintf("%s/login/%s?next=%s/user/token", common.AccAPIURL, host, common.AccAPIURL)
	fmt.Printf("Opening %s in the browser.\n", url)
	fmt.Println("Once authenticated, copy the token from there and paste it into terminal.")
	// open in browser
	err := common.OpenURL(url)
	if err != nil {
		fmt.Println("Failed to open browser. Please get the token manually from the link.")
		// won't work in docker, so don't err here
		// return err
	}
	// read input
	fmt.Print("> ")
	var token string
	fmt.Scanf("%s", &token)
	log.Debugf("Token: %s", token)
	// save to file
	err = session.SaveUserSession(token)
	if err != nil {
		return err
	}
	// show email
	spinner.StartText("Checking token")
	email, err := user.GetUserEmail()
	if err == nil {
		fmt.Printf("\nLogged in as %s\n", email)
	} else {
		log.Errorln(err)
		fmt.Println("\nFailed to get user info. Please try again.")
	}
	spinner.Stop()
	return err
}
