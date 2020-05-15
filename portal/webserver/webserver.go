package webserver

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/pinezapple/LibraryProject20201/portal/controller/authen"
	"github.com/pinezapple/LibraryProject20201/portal/core"
	dao "github.com/pinezapple/LibraryProject20201/portal/dao/database"
	"github.com/pinezapple/LibraryProject20201/skeleton/logger"
	"github.com/pinezapple/LibraryProject20201/skeleton/model"
)

// WebServer booting web server by configuration
func WebServer(ctx context.Context) (fn model.Daemon, err error) {
	// get configs
	mainConf, httpServerConf := core.GetConfig(), core.GetHTTPServerConf()
	lg := core.GetLogger()

	// create admin by default
	conf := core.GetConfig()
	dao.GetUserDAO().Create(context.Background(), core.GetDB(), conf.WebServer.Secure.SipHashSum0, conf.WebServer.Secure.SipHashSum1, "root", "20wPk29bnP2kc93nb92bzEobm", "", "", "", "", "")
	/*
		// try to load TLS Config
		var tlsConf *tls.Config
		tlsConf, err = xhttp.GenServerTLSConfig(httpServerConf)
		if err != nil {
			return
		}
	*/
	// Initialize router and http server
	e, server := echo.New(), &http.Server{
		Addr: fmt.Sprintf(":%d", httpServerConf.Port),
	}

	// Disable echo logging
	e.Logger.SetOutput(ioutil.Discard)

	// Recover middleware recovers from panics anywhere in the chain, prints stack trace and handles the control to the centralized HTTPErrorHandler.
	// Default stack size for trace is 4KB. For more example, please refer to https://echo.labstack.com/
	e.Use(mw.Recover())

	// Remove trailing slash middleware removes a trailing slash from the request URI.
	e.Pre(mw.RemoveTrailingSlash())

	// Set BodyLimit Middleware. It will panic if fail. For more example, please refer to https://echo.labstack.com/
	e.Use(mw.BodyLimit(mainConf.WebServer.BodyLimit))

	// Secure middleware provides protection against cross-site scripting (XSS) attack, content type sniffing, clickjacking, insecure connection and other code injection attacks.
	// For more example, please refer to https://echo.labstack.com/
	e.Use(mw.Secure())

	// Restricted group of URIs for
	r := e.Group("/r")
	r.Use(mw.JWTWithConfig(mw.JWTConfig{
		Claims:      &model.Claim{},
		ContextKey:  mainConf.WebServer.Secure.JWT.ContextKey,
		SigningKey:  []byte(mainConf.WebServer.Secure.JWT.SecretKey),
		TokenLookup: "header:authorization",
	}))

	// init router
	initRouter(e, r)

	//	core.Logger().Infof("HTTP Server is starting on :%d", httpServerConf.Port)
	logger.LogInfo(lg, "HTTP Server is starting on "+strconv.Itoa(httpServerConf.Port))
	// Start server
	go func() {
		if err := e.StartServer(server); err != nil {
			logger.LogErr(lg, err)
			//core.Logger().Error(err)
		}
	}()

	fn = func() {
		<-ctx.Done()

		// try to shutdown server
		if err := e.Shutdown(context.Background()); err != nil {
			logger.LogErr(lg, err)
			//	core.Logger().Error(err)
		} else {
			//	core.Logger().Warn("gracefully shutdown webserver")
			logger.LogInfo(lg, "gracefully shutdown webserver")
		}
	}
	return
}

func initRouter(e *echo.Echo, r *echo.Group) {

	/*
		initUserRouter(e)
		initGroupRouter(e)
	*/

	e.Static("/www", "www")
	// index
	e.GET("/", func(c echo.Context) error { return c.File("www/dist/index.html") })
	// login
	e.POST("/p/login", authen.Login)
}
