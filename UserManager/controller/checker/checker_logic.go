package checker

import (
	"DBproject1/core"
	"DBproject1/dao"
	"DBproject1/model"
	"fmt"
	"reflect"

	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func checkLogin(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*reqLoginCheck), c.Request().Context()

	// Log login info
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Login", Data: req.Username}

	// Validate request input
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))

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
		statusCode, err = http.StatusUnauthorized, fmt.Errorf("Password is incorrect. You only have %d times left to retry", user.PwdCounter)
		return
	}

	//TODO: get role and permission from casbin
	e := core.GetCasbinEnforcer()
	roles := e.GetRolesForUser(req.Username)
	permission := e.GetPermissionsForUser(req.Username)

	// Create token with claims
	expire := time.Now().Add(time.Minute * time.Duration(conf.WebServer.Secure.JWT.ExpireInMinute))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.Claim{
		Username:   req.Username,
		Group:      roles,
		Permission: permission,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	})

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(conf.WebServer.Secure.JWT.SecretKey))
	if err != nil {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Can not create secure token")
		return
	}

	return http.StatusOK, t, lg, false, nil
}

func checkPermission(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req := request.(*reqPermissionCheck)
	// Log login info
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Check Permission", Data: req.Token}

	e, conf := core.GetCasbinEnforcer(), core.GetConfig()
	if conf == nil {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Config not initialized")
		return
	}

	claim := model.Claim{}
	token, err := jwt.ParseWithClaims(req.Token, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.WebServer.Secure.JWT.SecretKey), nil
	})

	if !token.Valid {
		return http.StatusOK, false, lg, false, fmt.Errorf("Invalid or expired token")
	}

	pers := e.GetPermissionsForUser(claim.Username)
	roles := e.GetRolesForUser(claim.Username)
	/*
		if !reflect.DeepEqual(pers, claim.Permission) || !reflect.DeepEqual(roles, claim.Group) {
			return http.StatusOK, false, lg, false, fmt.Errorf("Invalid or expired token")
		}
	*/
	for i := 0; i < len(claim.Permission); i++ {
		ok := false
		for j := 0; j < len(pers); j++ {
			if reflect.DeepEqual(pers[j], claim.Permission[i]) {
				ok = true
				break
			}
		}
		if !ok {
			return http.StatusOK, false, lg, false, fmt.Errorf("Invalid or expired token")
		}
	}
	for i := 0; i < len(claim.Group); i++ {
		ok := false
		for j := 0; j < len(roles); j++ {
			if reflect.DeepEqual(roles[j], claim.Group[i]) {
				ok = true
				break
			}
		}
		if !ok {
			return http.StatusOK, false, lg, false, fmt.Errorf("Invalid or expired token")
		}
	}
	fmt.Println(req.Des)
	fmt.Println(req.Action)
	fmt.Println(claim.Username)
	
	if e.Enforce(claim.Username, req.Des, req.Action) == true {
		return http.StatusOK, true, lg, false, nil
	} else {
		return http.StatusOK, false, lg, false, fmt.Errorf("User is not allowed to execute this action")
	}
}
