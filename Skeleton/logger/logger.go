package logger

import (
	"encoding/json"
	"sync"

	"LibraryProject20201/Skeleton/model"

	log "github.com/sirupsen/logrus"
)

var serviceName = "Test service"

var locks sync.Mutex

func MustGet(servName string) {
	locks.Lock()
	serviceName = servName
	locks.Unlock()
}

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
