package slog

import (
	"fmt"
	"log/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func ErrWrapper(err error, op string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
