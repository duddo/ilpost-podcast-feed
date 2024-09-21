package endpoint

import (
	"net/http"
)

type CookieCache map[string][]*http.Cookie
