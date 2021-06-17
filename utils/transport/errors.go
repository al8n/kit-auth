package transport

import (
	"errors"
	"fmt"
)

type Methods = int8

type Type = int8

const (
	Request Methods = iota
	Response
)

const (
	GRPC Type =  iota
	Thrift
	HTTP
)

func errorCodecCasting(name string, method Methods) string {
	switch method {
	case Request:
		return fmt.Sprintf("error when casting request in %s", name)
	case Response:
		return fmt.Sprintf("error when casting response in %s", name)
	default:
		return fmt.Sprintf("error when casting in %s", name)
	}
}

func ErrorCodecCasting(name string, method Methods, typ Type) error  {
	switch typ {
	case GRPC:
		return errors.New(fmt.Sprintf("GRPC: %s", errorCodecCasting(name, method)))
	case Thrift:
		return errors.New(fmt.Sprintf("Thrift: %s", errorCodecCasting(name, method)))
	case HTTP:
		return errors.New(fmt.Sprintf("HTTP: %s", errorCodecCasting(name, method)))
	}

	return errors.New(fmt.Sprintf("Unkown: %s", errorCodecCasting(name, method)))
}

