package core

import (
	"encoding/json"
	"sync"

	"DBproject1/model"

	"github.com/spf13/viper"

	"github.com/casbin/casbin"
	mssqlx "github.com/linxGnu/mssqlx"
	log "github.com/sirupsen/logrus"
)

var (
	db                *mssqlx.DBs
	configLock        sync.RWMutex
	serverConf        *Config
	serverBindingConf *BindingConf
)

func init() {
	// if os.Getenv("APP_ENV") == "DEV" {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// } else {
	// 	viper.AutomaticEnv()
	// }

	//------------------------- Web Server -------------------------
	serverConf = &Config{
		WebServer: &WebServer{
			BodyLimit: viper.GetString("BODY_LIMIT"),
			Secure: Secure{
				SecureCookie: SecureCookie{
					CookieName:     "auth",
					ContextKey:     "lg_uid",
					MaxAge:         0,
					ExpireInMinute: 1440,
					HashKey:        "Wpdhgkwkpwngapoge93nx9sj2lsigwnx9529xn#px02naigm2-1$93*7nwlwsddf",
					BlockKey:       "DSgjo23058nxg84ns592*383nboiwkg+",
				},
				JWT: JWTConfig{
					ContextKey:     "lg_ctx_key",
					SecretKey:      "Wsjfi^sgjlkajskgwmsmvlk!utwEpc03q<qoP6[%4",
					ExpireInMinute: int64(viper.GetInt("JWT_EXPIRE_TIME")),
				},
				SipHashSum0: 947295729583939162,
				SipHashSum1: 323869573058327753,
			},
		},
		Database: &MysqlConnConfig{
			Type:        "mysql",
			DB:          viper.GetString("DB_NAME"),
			Username:    viper.GetString("DB_USER"),
			Password:    viper.GetString("DB_PASS"),
			Masters:     viper.GetStringSlice("DB_SLAVE"),
			Slaves:      viper.GetStringSlice("DB_MASTER"),
			Args:        viper.GetString("DB_ARGS"),
			MaxIdleConn: viper.GetInt("DB_MAX_IDLE"),
			MaxOpenConn: viper.GetInt("DB_MAX_OPEN"),
			IsWsrep:     false,
		},
	}

	serverBindingConf = &BindingConf{
		Port: viper.GetInt("PORT"),
	}
}

// GetConfig get current configuration for webserver
func GetConfig() (_serverConf *Config) {
	configLock.RLock()
	_serverConf = serverConf
	configLock.RUnlock()
	return
}

// SetConfig set current configuration for webserver
func SetConfig(_serverConf *Config) {
	configLock.RLock()
	serverConf = _serverConf
	configLock.RUnlock()
}

// GetDB get global database conn
func GetDB() (_db *mssqlx.DBs) {
	configLock.RLock()
	_db = db
	configLock.RUnlock()
	return
}

// SetDB set global db conn
func SetDB(_db *mssqlx.DBs) {
	configLock.RLock()
	db = _db
	configLock.RUnlock()
	return
}

var casbinEnforcerLock sync.RWMutex
var enforcer *casbin.Enforcer

// SetCasbinEnforcer set casbin enforcer
func SetCasbinEnforcer(e *casbin.Enforcer) {
	casbinEnforcerLock.RLock()
	enforcer = e
	casbinEnforcerLock.RUnlock()
}

// GetCasbinEnforcer get casbin enforcer
func GetCasbinEnforcer() (e *casbin.Enforcer) {
	casbinEnforcerLock.RLock()
	e = enforcer
	casbinEnforcerLock.RUnlock()
	return
}

// GetBindingConfUnsafe get webserver binding configuration in unsafe manner, no lock involve
func GetBindingConfUnsafe() *BindingConf {
	return serverBindingConf
}

// SetBindingConfUnsafe set webserver binding configuration in unsafe manner, no lock involve
func SetBindingConfUnsafe(v *BindingConf) {
	serverBindingConf = v
}

//------------------------- LOG UTILS -------------------------
const serviceName = "Horus"

// LogInfo information logging
func LogInfo(lg *model.LogFormat) {
	if lg == nil {
		return
	}
	lg.ServiceName = serviceName

	js, _ := json.Marshal(lg)
	log.Infof("%s", js)
}

// LogErr error logging
func LogErr(err error) {
	if err == nil {
		return
	}
	js, _ := json.Marshal(&model.LogFormat{ServiceName: serviceName, Err: err.Error()})
	log.Errorf("%s", js)
}

// LogWarning warning logging
func LogWarning(lg *model.LogFormat) {
	if lg == nil {
		return
	}
	lg.ServiceName = serviceName

	js, _ := json.Marshal(lg)
	log.Infof("%s", js)
}
