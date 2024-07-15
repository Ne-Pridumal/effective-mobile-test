package slog

import "log/slog"

func Debug(msg ...any) slog.Attr {
	return slog.Attr{
		Key:   "debug",
		Value: slog.AnyValue(msg),
	}
}
