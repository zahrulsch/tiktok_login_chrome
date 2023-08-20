package tiktokloginchrome

import (
	"context"
	"os"
	"path/filepath"

	"github.com/chromedp/chromedp"
)

type BrowserService struct {
	profilesDir string
	cookiesDir  string
}

type BrowserOptions struct {
	profilesDirPath string
	cookiesDirPath  string
}

func NewBrowserService(opt BrowserOptions) (*BrowserService, error) {
	profilesDirPath, err := filepath.Abs(opt.profilesDirPath)
	cookiesDirPath, err := filepath.Abs(opt.cookiesDirPath)

	if err != nil {
		return nil, err
	}

	if _, err := os.ReadDir(profilesDirPath); err != nil {
		if err = os.MkdirAll(profilesDirPath, os.ModeDir); err != nil {
			return nil, err
		}
	}

	if _, err := os.ReadDir(cookiesDirPath); err != nil {
		if err = os.MkdirAll(cookiesDirPath, os.ModeDir); err != nil {
			return nil, err
		}
	}

	return &BrowserService{
		profilesDir: profilesDirPath,
		cookiesDir:  cookiesDirPath,
	}, nil
}

func (b *BrowserService) NewBrowser(profileName string) *Browser {
	options := []func(*chromedp.ExecAllocator){}
	options = append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("user-data-dir", filepath.Join(b.profilesDir, profileName)),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("detach", true),
		chromedp.Flag("start-maximized", true),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	ctx, cancel = chromedp.NewContext(ctx)

	return &Browser{
		profilesDir: b.profilesDir,
		profileName: profileName,
		cookieDir:   b.cookiesDir,
		context:     ctx,
		cancelFunc:  cancel,
	}
}
