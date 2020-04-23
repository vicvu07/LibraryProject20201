package model

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/labstack/echo"
)

// SecureCookieConfig secure cookie middleware configuration
type SecureCookieConfig struct {
	HashKey    []byte
	BlockKey   []byte
	CookieName string
	ContextKey string
}

type cookieValidator struct {
	secureCookie *securecookie.SecureCookie
	config       *SecureCookieConfig
}

func (c cookieValidator) MakeSecureCookie(val string) (*http.Cookie, error) {
	if c.secureCookie == nil || c.config == nil {
		return nil, fmt.Errorf("CookieValidator not initialized")
	}

	encoded, err := c.secureCookie.Encode(c.config.CookieName, val)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:  c.config.CookieName,
		Value: encoded,
	}, nil
}

func (c cookieValidator) ExpireSecureCookie() (*http.Cookie, error) {
	if c.secureCookie == nil || c.config == nil {
		return nil, fmt.Errorf("CookieValidator not initialized")
	}

	return &http.Cookie{
		Name:   c.config.CookieName,
		MaxAge: -1,
	}, nil
}

// CookieValidator ...
var CookieValidator = cookieValidator{}

func readSecureCookie(secureCookie *securecookie.SecureCookie, c echo.Context, cookieName string) (value string, err error) {
	_cookie, err := c.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	cookie := _cookie.Value

	var val string
	err = secureCookie.Decode(cookieName, cookie, &val)
	value = val

	return
}

// NewSecureCookieMW new secure cookie middleware
func NewSecureCookieMW(config SecureCookieConfig) echo.MiddlewareFunc {
	CookieValidator.secureCookie = securecookie.New(config.HashKey, config.BlockKey)

	if len(config.ContextKey) == 0 {
		config.ContextKey = "USER"
	}

	if len(config.CookieName) == 0 {
		config.CookieName = "auth"
	}

	CookieValidator.config = &config

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cookie, err := readSecureCookie(CookieValidator.secureCookie, c, CookieValidator.config.CookieName); err == nil {
				c.Set(config.ContextKey, cookie)
			} else {
				return echo.ErrUnauthorized
			}

			// Continue the chain of middleware
			return next(c)
		}
	}
}
