package core

import (
	"fmt"
	"time"

	"DBproject1/lib"
	"DBproject1/model"

	jwt "github.com/dgrijalva/jwt-go"
)

// WebServer hold configurations for WebServer
type WebServer struct {
	// BodyLimit The body limit is determined based on both Content-Length request header and actual content read, which makes it super secure.
	// Limit can be specified as 4x or 4xB, where x is one of the multiple from K, M, G, T or P. Example: 2M = 2 Megabyte
	BodyLimit string
	// Secure configuration
	Secure Secure
}

// MysqlConnConfig mysql database connection configuration
type MysqlConnConfig struct {
	Type        string
	DB          string
	Username    string
	Password    string
	Masters     []string
	Slaves      []string
	Args        string
	MaxIdleConn int
	MaxOpenConn int
	IsWsrep     bool
}

func (c *MysqlConnConfig) Construct() (err error) {
	if c.Masters != nil && len(c.Masters) > 0 {
		for i := range c.Masters {
			c.Masters[i] = fmt.Sprintf("%s:%s@(%s)/%s?%s", c.Username, c.Password, c.Masters[i], c.DB, c.Args)
		}
	}

	if c.Slaves != nil && len(c.Slaves) > 0 {
		for i := range c.Slaves {
			c.Slaves[i] = fmt.Sprintf("%s:%s@(%s)/%s?%s", c.Username, c.Password, c.Slaves[i], c.DB, c.Args)
		}
	}

	return
}

// Config main configuration
type Config struct {
	// WebServer configuration
	WebServer *WebServer `json:"WebServer"`
	// Database configuration
	Database *MysqlConnConfig `json:"Database"`
}

// BindingConf binding configuration for webserver
type BindingConf struct {
	Port int
	Cert string
	Key  string
}

// JWTConfig configuration for jwt token within web server
type JWTConfig struct {
	// ContextKey to get JWT token from context
	ContextKey string
	// SecretKey to generate JWT Token
	SecretKey string
	// ExpireInMinute jwt token will expire after minutes
	ExpireInMinute int64
}

// SecureCookie secure cookie configuration
type SecureCookie struct {
	// CookieName name of secure cookie
	CookieName string
	// ContextKey to get SecureCookie from context
	ContextKey string
	// MaxAge of cookie
	MaxAge int
	// ExpireInMinute
	ExpireInMinute int64
	// HashKey 64 character
	HashKey string
	// BlockKey 32 character
	BlockKey string
}

// Secure config
type Secure struct {
	// JWT for web application/mobile
	JWT JWTConfig

	// SecureCookie secure cookie configuration
	SecureCookie SecureCookie

	// SipHashSum0
	SipHashSum0 uint64
	// SipHashSum1
	SipHashSum1 uint64
}

var (
	config = &Config{}
)

// Sign jwt token
func (c *JWTConfig) Sign(claim *model.Claim) (string, error) {
	if claim == nil {
		return "", fmt.Errorf("Claim is nil")
	}

	// modify claim for expire at
	claim.ExpiresAt = time.Now().Add(time.Minute * time.Duration(c.ExpireInMinute)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// generate encoded token and send it as response.
	return token.SignedString([]byte(c.SecretKey))
}

// SipHash do sip hash 4-8 sum
func (c *Secure) SipHash(payload []byte) uint64 {
	return uint64(lib.SipHash48(c.SipHashSum0, c.SipHashSum1, payload))
}
