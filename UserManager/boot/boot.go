package boot

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"DBproject1/core"
	"DBproject1/model"
	webServer "DBproject1/webserver"

	mssqlxA "github.com/cicdata-io/mssqlx-adapter"

	"github.com/casbin/casbin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/linxGnu/mssqlx"
	log "github.com/sirupsen/logrus"
)

// ------------------------ DAEMONS MANAGEMENT ------------------------
var daemonChannels = []model.TermChan{}
var daemonChannelLock sync.Mutex

func registerDaemons(wg *sync.WaitGroup, testMode bool) {
	// Web Server
	registerDaemon(wg, testMode, webServer.WebServer)
}

func registerDaemon(wg *sync.WaitGroup, testMode bool, halfDaemon model.HalfDaemonFunc) {
	if halfDaemon != nil {
		dm := model.NewTermChan()
		daemon, err := halfDaemon(dm, testMode)
		if err != nil {
			panic(err)
		}

		_registerDaemon(wg, dm, testMode, daemon)
	}
}

func _registerDaemon(wg *sync.WaitGroup, dm model.TermChan, testMode bool, daemon model.DaemonFunc) {
	if daemon != nil {
		daemonChannelLock.Lock()
		defer daemonChannelLock.Unlock() // no need to prevent defer()

		daemonChannels = append(daemonChannels, dm)

		wg.Add(1)
		go daemon(wg, dm, testMode)
	}
}

func killDaemons(sig os.Signal) {
	daemonChannelLock.Lock()
	defer daemonChannelLock.Unlock() // no need to prevent defer()

	for i := len(daemonChannels) - 1; i >= 0; i-- {
		daemonChannels[i] <- sig
	}
}

// ------------------------ Bootstrap ------------------------
// Boot app with mode (test or production). Return error if something wrong.
// If not, run all daemons and wait until os(signal kill) or external termination signal.
// Daemons would be notified on killed for graceful shutdown
func Boot(externalTermination model.TermChan, testMode bool) (err error) {
	// check termination channel
	if cap(externalTermination) == 0 {
		err = core.ErrExtTermChanCapInvalid
		return
	}

	if !testMode {
		// Log as JSON instead of the default ASCII formatter
		// for easy parsing by logstash or Splunk.
		log.SetFormatter(&log.JSONFormatter{})
	}

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// change working directory
	currentDir, err := os.Getwd()
	_ = os.Chdir(currentDir)

	var dbs *mssqlx.DBs
	// log mode
	if testMode {
		core.LogInfo(&model.LogFormat{Action: "Boot", Data: "Test mode"})
		dbs, err = core.InitDBForTestUnsafe()
	} else {
		core.LogInfo(&model.LogFormat{Action: "Boot", Data: "Production mode"})
		// try database construct
		if err = core.GetConfig().Database.Construct(); err != nil {
			return
		}

		// Init db
		dbs, err = core.InitDBForProduction(core.GetConfig().Database)
	}

	if err != nil {
		panic(err)
	}

	dbs.SetMaxIdleConns(core.GetConfig().Database.MaxIdleConn)
	dbs.SetMaxOpenConns(core.GetConfig().Database.MaxOpenConn)
	core.SetDB(dbs)

	// create casbin enforcer
	a := mssqlxA.NewAdapter(core.GetConfig().Database.Type, core.GetConfig().Database.Masters, core.GetConfig().Database.Slaves, true)
	e := casbin.NewEnforcer("model.conf", a)
	core.SetCasbinEnforcer(e)

	return bootstrapDaemons(externalTermination, testMode)
}

func bootstrapDaemons(extTerm <-chan os.Signal, testMode bool) (err error) {
	defer core.LogWarning(&model.LogFormat{Action: "Shutdown", Success: "Gracefully shutdown daemons"})

	// os signal handling
	sigs := make(chan os.Signal, 10)
	done := make(chan os.Signal, 10)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGTSTP)

	// waiter all daemons
	var wg sync.WaitGroup

	// register daemons
	registerDaemons(&wg, testMode)

	// catch signal
	go func() {
		select {
		case sig := <-sigs:
			core.LogWarning(&model.LogFormat{Action: "HandleSignal", Data: sig, Source: "Internal"})
			done <- sig
		case sig := <-extTerm:
			core.LogWarning(&model.LogFormat{Action: "HandleSignal", Data: sig, Source: "External"})
			done <- sig
		}
	}()

	// kill daemons
	killDaemons(<-done)

	// wait all daemons done
	wg.Wait()

	return
}
