package tiktokloginchrome

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type Browser struct {
	profilesDir string
	profileName string
	cookieDir   string
	context     context.Context
	cancelFunc  context.CancelFunc
}

var sellerInfoUrl = "https://seller-id.tiktok.com/profile/seller-profile?tab=seller_information_old"
var sellerUsernameSel = "div[data-tid=\"profile.view.seller_information_old.shop_name\"] > div"
var chatFloatSel = "div[data-tid=\"im.action.im_row_bottom_entrance\"]"
var avatarSel = "div[data-tid=\"m4b_avatar\"] img"

func (b *Browser) LoginTiktok(a *Account) (err error) {
	err = loginWithCookie(b, a)
	if err != nil {
		loginUrl := "https://seller-id.tiktok.com/account/login"
		ctx := b.context

		err = chromedp.Run(
			ctx,
			chromedp.Navigate(loginUrl),
		)

		if err != nil {
			return
		}

		select {
		case err = <-isLogin(b):
			if err != nil {
				return
			}
		case err = <-loginSteps(b, a):
			if err != nil {
				return
			}
		case <-time.After(time.Minute * 3):
			return errors.New(fmt.Sprintf("Proses login timeout: timeout %v menit", 3))
		}

		err = chromedp.Run(
			ctx,
			chromedp.Navigate(sellerInfoUrl),
			chromedp.WaitVisible(sellerUsernameSel),
			chromedp.TextContent(sellerUsernameSel, &a.Username),
			chromedp.AttributeValue(avatarSel, "src", &a.AvatarURL, nil),
			chromedp.ActionFunc(func(ctx context.Context) error {
				cookies, err := network.GetCookies().Do(ctx)
				if err != nil {
					return err
				}

				for _, c := range cookies {
					a.Cookies = append(a.Cookies, &http.Cookie{
						Name:     c.Name,
						Value:    c.Value,
						Path:     c.Path,
						Domain:   c.Domain,
						Expires:  time.Now().Add(2628000 * time.Second),
						HttpOnly: c.HTTPOnly,
					})
				}

				return nil
			}),
		)
	}

	jsonB, err := json.Marshal(a)
	if err != nil {
		return
	}

	if a.Email != "" {
		cookieName := a.Email
		f, err := os.Create(filepath.Join(b.cookieDir, cookieName+".json"))
		defer f.Close()

		if err != nil {
			return err
		}

		_, err = f.Write(jsonB)
		if err != nil {
			return err
		}
	}

	return
}

func urlParser(a *Account) (uri *url.URL, err error) {
	targetUrl, err := url.Parse("https://seller-id.tiktok.com/api/v1/seller/get")
	msToken := ""

	for _, c := range a.Cookies {
		if c.Name == "msToken" {
			msToken = c.Value
		}
	}

	q := targetUrl.Query()
	q.Set("locale", "id-ID")
	q.Set("language", "id")
	q.Set("aid", "4068")
	q.Set("app_name", "i18n_ecom_shop")
	q.Set("device_id", "0")
	q.Set("device_platform", "web")
	q.Set("cookie_enabled", "true")
	q.Set("screen_width", "1920")
	q.Set("screen_height", "1080")
	q.Set("browser_language", "en-US")
	q.Set("browser_platform", "Win32")
	q.Set("browser_name", "Mozilla")
	q.Set("browser_version", "5.0%20%28Windows%20NT%2010.0%3B%20Win64%3B%20x64%29%20AppleWebKit%2F537.36%20%28KHTML%2C%20like%20Gecko%29%20Chrome%2F115.0.0.0%20Safari%2F537.36")
	q.Set("browser_online", "true")
	q.Set("timezone_name", "Asia%2FJakarta")
	q.Set("need_newest_logo", "true")
	q.Set("msToken", msToken)

	targetUrl.RawQuery = q.Encode()

	return targetUrl, err
}

func loginWithCookie(b *Browser, a *Account) (err error) {
	cookiePath := filepath.Join(b.cookieDir, a.Email+".json")
	f, err := os.ReadFile(cookiePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(f, &a)
	if err != nil {
		return
	}

	uri, err := urlParser(a)
	cookies := a.Cookies

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return
	}

	jar.SetCookies(uri, cookies)

	client := &http.Client{Jar: jar}
	req, err := http.NewRequest("GET", uri.String(), nil)

	if err != nil {
		return
	}

	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	response, err := client.Do(req)

	if err != nil {
		return
	}

	defer response.Body.Close()
	reader, err := gzip.NewReader(response.Body)
	defer reader.Close()

	if err != nil {
		return
	}

	var accountResponse AccountResponse
	body, err := io.ReadAll(reader)
	err = json.Unmarshal(body, &accountResponse)

	if err != nil {
		return
	}

	a.AccountID = accountResponse.Data.Seller.SellerID

	return
}

func isLogin(b *Browser) (errChan chan error) {
	ctx := b.context

	errChan = make(chan error)

	go func() {
		err := chromedp.Run(
			ctx,
			chromedp.WaitVisible(chatFloatSel),
		)

		errChan <- err
	}()

	return
}

func loginSteps(b *Browser, a *Account) (errChan chan error) {
	ctx := b.context
	logWithEmailSel := "#TikTok_Ads_SSO_Login_Email_Panel_Button"
	emailInputSel := "#TikTok_Ads_SSO_Login_Email_Input"
	pwdInputSel := "#TikTok_Ads_SSO_Login_Pwd_Input"
	loginBtnSel := "#TikTok_Ads_SSO_Login_Btn"

	errChan = make(chan error)

	go func() {
		err := chromedp.Run(
			ctx,
			chromedp.WaitVisible(logWithEmailSel),
			chromedp.Click(logWithEmailSel),
			chromedp.Clear(emailInputSel),
			chromedp.SendKeys(emailInputSel, a.Email),
			chromedp.Clear(pwdInputSel),
			chromedp.SendKeys(pwdInputSel, a.Password),
			chromedp.Click(loginBtnSel),
			chromedp.WaitVisible(chatFloatSel),
		)

		errChan <- err
	}()

	return
}
