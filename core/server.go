// Copyright 2026 Ting. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// RunServer starts the reverse proxy and blocks until shutdown.
func RunServer() error {
	target, err := url.Parse(viper.GetString("upstream.url"))
	if err != nil {
		return err
	}

	token := viper.GetString("upstream.api_key")
	port := viper.GetInt("app.port")
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", strconv.Itoa(port)),
		Handler:           logRequests(newHandler(target, token)),
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	serr := make(chan error, 1)

	go func() {
		log.Info().
			Int("port", port).
			Str("upstream", target.String()).
			Msg("Starting proxy server")

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			serr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serr:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Info().
			Str("signal", sig.String()).
			Msg("Received shutdown signal")

		shutdownTimeout := 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		log.Info().
			Dur("timeout", shutdownTimeout).
			Msg("Gracefully shutting down server")

		err := srv.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}

		log.Info().Msg("Server shutdown complete")
	}

	return nil
}
