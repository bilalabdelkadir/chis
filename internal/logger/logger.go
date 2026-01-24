package logger

import (
	"log/slog"
	"os"
)

func Setup() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(handler))
}
