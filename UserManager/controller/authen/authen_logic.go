package authen

import (
	"DBproject1/core"
	"DBproject1/dao"
	"DBproject1/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func login(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*reqLogin), c.Request().Context()

	// Log login info
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Login", Data: req.Username}

	// Validate request input
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	if req.Username != "root" {
		statusCode, err = http.StatusUnauthorized, fmt.Errorf("Not root User")
		return
	}
	db, conf := core.GetDB(), core.GetConfig()
	if db == nil || conf == nil {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Database connection / Config not initialized")
		return
	}

	// Select user by his username
	userDAO := dao.GetUserDAO()

	user, er := userDAO.SelectByUsername(ctx, db, req.Username)
	if er != nil {
		statusCode, err = http.StatusUnauthorized, er
		return
	} else if user == nil {
		statusCode, err = http.StatusUnauthorized, fmt.Errorf("User not found")
		return
	}
	if !user.ValidateChecksum(conf.WebServer.Secure.SipHashSum0, conf.WebServer.Secure.SipHashSum1) {
		statusCode, err = http.StatusUnauthorized, fmt.Errorf("Invalid user checksum")
		return
	}

	// TODO: using constants later
	if user.Status != 1 {
		statusCode, err = http.StatusUnauthorized, fmt.Errorf("This account has been blocked. Please contact admin of this site")
		return
	}

	// Select user_sec by his username
	userSec, er := userDAO.SelectSecByUsername(ctx, db, req.Username)
	if err != nil {
		statusCode, err = http.StatusUnauthorized, fmt.Errorf("User not found")
		return
	}
	if !userSec.ValidateChecksum(conf.WebServer.Secure.SipHashSum0, conf.WebServer.Secure.SipHashSum1) {
		statusCode, err = http.StatusUnauthorized, fmt.Errorf("Invalid user sec checksum")
		return
	}

	if err = bcrypt.CompareHashAndPassword(userSec.Password, []byte(req.Password)); err != nil {
		user.PwdCounter--
		userDAO.UpdatePasswordCounter(ctx, db, user.ID, user.PwdCounter)

		if user.PwdCounter <= 0 {
			// TODO: using const later
			userDAO.UpdateStatus(ctx, db, user.ID, 0)
		}

		statusCode, err = http.StatusUnauthorized, fmt.Errorf("Password is incorrect. You only have %d times left to retry", user.PwdCounter)
		return
	}
	userDAO.UpdatePasswordCounter(ctx, db, user.ID, 5)

	//TODO: get role and group from casbin

	// Create token with claims
	expire := time.Now().Add(time.Minute * time.Duration(conf.WebServer.Secure.SecureCookie.ExpireInMinute))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.Claim{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	})

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(conf.WebServer.Secure.JWT.SecretKey))
	if err != nil {
		log.Println("im here")
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Can not create secure token")
		return
	}

	// Set secure cookie
	cookie, err := model.CookieValidator.MakeSecureCookie(strconv.FormatUint(user.ID, 10))
	if err != nil {
		log.Println("im here 1", err)
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Can not create secure token")
		return
	}
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(time.Duration(conf.WebServer.Secure.SecureCookie.ExpireInMinute) * time.Minute)
	cookie.MaxAge = conf.WebServer.Secure.SecureCookie.MaxAge
	c.SetCookie(cookie)

	return http.StatusOK, t, lg, false, nil
}
