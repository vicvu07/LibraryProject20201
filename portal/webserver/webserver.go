package webserver

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/pinezapple/LibraryProject20201/portal/core"
	"github.com/pinezapple/LibraryProject20201/skeleton/logger"
	"github.com/pinezapple/LibraryProject20201/skeleton/model"
)

// WebServer booting web server by configuration
func WebServer(ctx context.Context) (fn model.Daemon, err error) {
	// get configs
	mainConf, httpServerConf := core.GetConfig(), core.GetHTTPServerConf()
	lg := core.GetLogger()

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

	// init router
	//initRouter(e)

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
