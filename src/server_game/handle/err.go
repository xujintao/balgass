package handle

import (
	"fmt"
	"log"
)

func init() {
	mapConfigErrors = make(map[int]*ConfigError)
	for _, v := range configErrors {
		s := v.service
		if dup, ok := mapConfigErrors[s]; ok {
			log.Panicf("duplicated config error [service]%d [description]%s",
				dup.service, dup.Description)
		}
		mapConfigErrors[s] = v
	}
}

func MakeError(service int, err error) *ConfigError {
	ce, ok := mapConfigErrors[service]
	if !ok {
		log.Printf("cannot find config error [service]%d\n", service)
		return MakeError(Unknown, nil)
	}
	ce.err = err
	return ce
}

var mapConfigErrors map[int]*ConfigError

type ConfigError struct {
	service     int
	Code        int
	Description string
	err         error
}

func (ce *ConfigError) Error() string {
	return fmt.Sprintf("[description]%s [err]%v", ce.Description, ce.err)
}

const (
	Unknown int = iota
	CreateAccountBind
	CreateAccountValidate
	CreateAccountDB
	GetAccountListDB
	DeleteAccountMissingParam
	DeleteAccountDB
)

var configErrors = [...]*ConfigError{
	{Unknown, 500, "Internal server error", nil},
	{CreateAccountBind, 500, "create account bind body failed", nil},
	{CreateAccountValidate, 500, "create account validate body failed", nil},
	{CreateAccountDB, 500, "create account db failed", nil},
	{GetAccountListDB, 500, "get account list db failed", nil},
	{DeleteAccountMissingParam, 500, "delete account missing param", nil},
	{DeleteAccountDB, 500, "delete account db failed", nil},
}
