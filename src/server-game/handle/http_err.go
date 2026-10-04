package handle

import (
	"log/slog"
	"os"
)

func init() {
	mapConfigErrors = make(map[int]*ConfigError)
	for _, v := range configErrors {
		s := v.service
		if dup, ok := mapConfigErrors[s]; ok {
			slog.Error("duplicated config error",
				"service", dup.service, "description", dup.Description)
			os.Exit(1)
		}
		mapConfigErrors[s] = v
	}
}

func MakeError(service int, err error) *ConfigError {
	ce, ok := mapConfigErrors[service]
	if !ok {
		slog.Error("cannot find config error", "service", service)
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
	return ce.err.Error()
}

const (
	Unknown int = iota
	CommandBind
	CommandValidate
	CommandInInvalid
	CommandExec
)

var configErrors = [...]*ConfigError{
	{Unknown, 500, "Internal server error", nil},
	{CommandBind, 400, "command body failed", nil},
	{CommandValidate, 400, "command validate failed", nil},
	{CommandInInvalid, 400, "command in field is invalid", nil},
	{CommandExec, 500, "command exec failed", nil},
}
