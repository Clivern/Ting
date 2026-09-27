// Copyright 2026 Ting. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rs/zerolog/log"
)

func newHandler(target *url.URL, token string) http.Handler {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.Out.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Error().
				Err(err).
				Str("method", r.Method).
				Str("uri", r.URL.RequestURI()).
				Msg("upstream request failed")
			http.Error(w, "bad gateway", http.StatusBadGateway)
		},
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("method", r.Method).
			Str("uri", r.URL.RequestURI()).
			Msg("request")
		next.ServeHTTP(w, r)
	})
}
