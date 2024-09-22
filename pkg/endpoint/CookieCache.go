package endpoint

import (
	"net/http"
	"time"
)

type CookieCache map[string][]*http.Cookie

func (storage CookieCache) TryGetValidCookie(username string) (bool, []*http.Cookie) {
	var cookies = storage[username]

	if cookies == nil {
		return false, nil
	}

	for _, cookie := range cookies {
		if !cookie.Expires.IsZero() && cookie.Expires.Before(time.Now()) {
			delete(storage, username)
			return false, nil
		}
	}

	return true, cookies
}

func (storage CookieCache) Add(username string, cookies []*http.Cookie) {
	storage[username] = cookies
}
