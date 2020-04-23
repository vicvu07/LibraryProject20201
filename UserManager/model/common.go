package model

import (
	"context"
	"database/sql"
	"os"
	"sync"
)

// TermChan os.Signal channel for termination
type TermChan chan os.Signal

// NewTermChan create new termination channel with capacity 1
func NewTermChan() TermChan {
	return make(TermChan, 1)
}

// LogFormat log pattern for whole system
type LogFormat struct {
	ServiceName string      `json:"srv"`
	Source      string      `json:"src,omitempty"`
	Action      string      `json:"act,omitempty"`
	Data        interface{} `json:"dat,omitempty"` // data
	Err         interface{} `json:"err,omitempty"` // error
	Success     interface{} `json:"suc,omitempty"` // success
}

// DaemonFunc daemon function
type DaemonFunc func(wg *sync.WaitGroup, termination TermChan, testMode bool)

// HalfDaemonFunc return a daemon function for invoker to do next step
type HalfDaemonFunc func(termination TermChan, testMode bool) (DaemonFunc, error)

// ExecTransaction template to execute transaction
func ExecTransaction(ctx context.Context, tx *sql.Tx, exec func(ctx context.Context, tx *sql.Tx) error) (err error) {
	if err = exec(ctx, tx); err != nil {
		tx.Rollback()
		return
	}

	return tx.Commit()
}
