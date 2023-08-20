package tiktokloginchrome

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestWarmUp(t *testing.T) {
	options := BrowserOptions{
		profilesDirPath: "profiles",
		cookiesDirPath:  "cookies",
	}

	serv, err := NewBrowserService(options)
	if err != nil {
		panic(err)
	}

	account := Account{
		Username: "",
		Email:    "mantracode@yahoo.com",
		Password: "Muhammad123!",
	}

	browser := serv.NewBrowser(account.Email)
	defer browser.cancelFunc()

	err = browser.LoginTiktok(&account)
	if err != nil {
		panic(err)
	}

	jsonByte, err := json.Marshal(account)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#+v\n", string(jsonByte))
}
