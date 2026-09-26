// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	token := os.Getenv("OPENROUTER_API_KEY")
	if token == "" {
		log.Fatal("OPENROUTER_API_KEY is required")
	}

	target, err := url.Parse("https://openrouter.ai")
	if err != nil {
		log.Fatal(err)
	}

	addr := listenAddr()
	log.Printf("proxying %s -> %s", addr, target)
	log.Fatal(http.ListenAndServe(addr, logRequests(newHandler(target, token))))
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func listenAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	}
	if port[0] == ':' {
		return port
	}
	return ":" + port
}

func newHandler(target *url.URL, token string) http.Handler {
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			r.Out.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("%s %s: %v", r.Method, r.URL.RequestURI(), err)
			http.Error(w, "bad gateway", http.StatusBadGateway)
		},
	}
}
