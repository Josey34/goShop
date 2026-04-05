package logger

import (
	"log/slog"
	"os"
)

func Setup(env string) {
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	slog.SetDefault(slog.New(handler))
}
