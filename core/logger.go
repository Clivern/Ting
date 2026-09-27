// Copyright 2026 Ting. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// SetupLogging configures the logging system based on viper configuration.
func SetupLogging() error {
	output := viper.GetString("app.log.output")
	writer, err := logWriter(output)
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to read hostname")
	}

	if viper.GetString("app.log.format") == "json" {
		log.Logger = zerolog.New(writer).With().Timestamp().Str("hostname", hostname).Logger()
	} else {
		log.Logger = zerolog.New(
			zerolog.ConsoleWriter{Out: writer},
		).With().Timestamp().Str("hostname", hostname).Logger()
	}

	level := strings.ToLower(viper.GetString("app.log.level"))

	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return nil
}

func logWriter(output string) (io.Writer, error) {
	if output == "" || output == "stdout" {
		return os.Stdout, nil
	}

	dir := filepath.Dir(output)
	if err := os.MkdirAll(dir, 0o775); err != nil {
		return nil, fmt.Errorf("directory [%s] creation failed: %w", dir, err)
	}

	f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o664)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %w", err)
	}

	return f, nil
}
