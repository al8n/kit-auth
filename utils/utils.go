package utils

import (
	"errors"

	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
)

// These annoying helper functions are required to translate Go error types to
// and from strings, which is the type we use in our IDLs to represent errors.
// There is special casing to treat empty strings as nil errors.
func Str2Err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func Err2Str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func SetTracerSpanError(span opentracing.Span) {
	span.SetTag("error", true)
}

func SetAndLogTracerSpanError(span opentracing.Span, err error)  {
	span.SetTag("error", true)
	span.LogFields(otlog.Error(err))
}

func LogTracerError(span opentracing.Span, errs ...error)  {
	var otErrs  = make([]otlog.Field, len(errs))
	for i, err := range errs {
		otErrs[i] = otlog.Error(err)
	}
	span.LogFields(otErrs...)
}
