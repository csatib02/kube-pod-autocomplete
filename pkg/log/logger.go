package log

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"slices"

	slogmulti "github.com/samber/slog-multi"
	slogsyslog "github.com/samber/slog-syslog"

	"github.com/csatib02/kube-pod-autocomplete/internal/config"
)

func InitLogger(config *config.Config) {
	var level slog.Level

	err := level.UnmarshalText([]byte(config.LogLevel))
	if err != nil { // Silently fall back to info level
		level = slog.LevelInfo
	}

	levelFilter := func(levels ...slog.Level) func(ctx context.Context, r slog.Record) bool {
		return func(_ context.Context, r slog.Record) bool {
			return slices.Contains(levels, r.Level)
		}
	}

	router := slogmulti.Router()
	if config.JSONLog {
		// Send logs with level higher than warning to stderr
		router = router.Add(
			slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level}),
			levelFilter(slog.LevelWarn, slog.LevelError),
		)

		// Send info and debug logs to stdout
		router = router.Add(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
			levelFilter(slog.LevelDebug, slog.LevelInfo),
		)
	} else {
		// Send logs with level higher than warning to stderr
		router = router.Add(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}),
			levelFilter(slog.LevelWarn, slog.LevelError),
		)

		// Send info and debug logs to stdout
		router = router.Add(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
			levelFilter(slog.LevelDebug, slog.LevelInfo),
		)
	}

	if config.LogServerAddress != "" {
		writer, err := net.Dial("udp", config.LogServerAddress)
		if err != nil {
			slog.Error(fmt.Errorf("failed to connect to log server: %w", err).Error())
			os.Exit(1)
		}

		router = router.Add(slogsyslog.Option{Level: slog.LevelInfo, Writer: writer}.NewSyslogHandler())
	}

	logger := slog.New(router.Handler()).With(slog.String("app", "kube-pod-autocomplete"))

	// Set the default logger to the configured logger,
	// enabling direct usage of the slog package for logging.
	slog.SetDefault(logger)
}
