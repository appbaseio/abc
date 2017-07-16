package logout

import (
	"fmt"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"os"
)

// UserLogout log outs a user
func UserLogout(doAll bool) error {
	err := session.DeleteUserSession()
	if err != nil {
		return err
	}
	// remove env var
	os.Unsetenv("ABC_TOKEN")
	fmt.Println("Logged out successfully")
	// logout from browser
	if doAll {
		url := fmt.Sprintf("%s/logout", common.AccAPIURL)
		fmt.Printf("Opening %s in the browser to log you out.\n", url)
		err := common.OpenURL(url)
		if err != nil {
			fmt.Println("Failed to open browser. Please logout manually using the link.")
		}
	}
	return nil
}
