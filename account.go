package tiktokloginchrome

import "net/http"

type Account struct {
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"-"`
	AvatarURL string         `json:"avatarURL"`
	AccountID string         `json:"accountID"`
	Cookies   []*http.Cookie `json:"cookies"`
}
