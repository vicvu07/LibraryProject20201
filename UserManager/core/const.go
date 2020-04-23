package core

import "fmt"

const (
	DEFAULTPORT = 8806
	testDSN     = "root:123@/test?charset=utf8&collation=utf8_general_ci&parseTime=true&loc=Asia%2FHo_Chi_Minh"
)

var (
	ErrBadRequest = fmt.Errorf("Bad request")

	ErrExtTermChanCapInvalid = fmt.Errorf("Term chan capacity is invalid")

	// ErrDBObjNull indicate DB Object is nil
	ErrDBObjNull = fmt.Errorf("DB Object is nil")
)
