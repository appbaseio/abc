package login

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
	"github.com/appbaseio/abc/appbase/user"
	"github.com/appbaseio/abc/log"
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
	url := fmt.Sprintf(
		"%s/logout?next=%s/login/%s?next=%s/user/token",
		common.AccAPIURL, common.AccAPIURL, host, common.AccAPIURL)
	fmt.Printf("Opening the following url in the browser.\n")
	fmt.Println(url)
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

// StartUserLoginBasicAuth starts user login process
func StartUserLoginBasicAuth(creds string) error {
	req, err := http.NewRequest("GET", common.AccAPIURL, nil)
	if err != nil {
		fmt.Println("Failed to initialize request")
	}

	credSlice := strings.Split(creds, ":")
	req.SetBasicAuth(credSlice[0], credSlice[1])

	spinner.Start()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to perform request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Bad request returned %d, wasn't able to login user", resp.StatusCode)
	}

	err = user.ShowUserEmail()
	if err != nil {
		log.Errorln(err)
		fmt.Println("\nFailed to get user info. Please try again.")
	}

	spinner.Stop()
	return err
}
