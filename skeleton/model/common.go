package model

import (
	"context"
	"database/sql"
)

// Daemon abstract a daemon.
type Daemon func()

// DaemonGenerator generates a Daemon
type DaemonGenerator func(ctx context.Context) (Daemon, error)

// RPCErrCode rpc error code
type RPCErrCode int32

const (
	// RPCSuccess request success without any error
	RPCSuccess RPCErrCode = iota
	// RPCCustomErr custom error
	RPCCustomErr
	// RPCInternalServerErr internal server error
	RPCInternalServerErr
)

// codes for lead status update service response
const (
	CustomRPCErr RPCErrCode = -1

	UpdateSuccess RPCErrCode = 0

	InvalidLeadIDErr RPCErrCode = 1

	NotClientLeadErr RPCErrCode = 2

	InvalidLeadInfoErr RPCErrCode = 3
)

// ExecTransaction template to execute transaction
func ExecTransaction(ctx context.Context, tx *sql.Tx, exec func(ctx context.Context, tx *sql.Tx) error) (err error) {
	if err = exec(ctx, tx); err != nil {
		_ = tx.Rollback()
		return
	}

	return tx.Commit()
}

// LogFormat log pattern for whole system
type LogFormat struct {
	ServiceName string      `json:"srv"`
	Source      string      `json:"src,omitempty"`
	Action      string      `json:"act,omitempty"`
	Data        interface{} `json:"dat,omitempty"`   // data
	Err         interface{} `json:"err,omitempty"`   // error
	Success     interface{} `json:"suc,omitempty"`   // success
	Stack       interface{} `json:"stack,omitempty"` // stack trace
	Message     string      `json:"msg,omitempty"`
}

// ToMapStringItf ...
func (lg *LogFormat) ToMapStringItf() map[string]interface{} {
	return map[string]interface{}{
		"srv":   lg.ServiceName,
		"src":   lg.Source,
		"act":   lg.Action,
		"dat":   lg.Data,
		"err":   lg.Err,
		"suc":   lg.Success,
		"stack": lg.Stack,
	}
}

// ShardAddress address of shard, point to lead service
type ShardAddress struct {
	ID      uint64 `json:"id"`
	Address string `json:"address"`
}
