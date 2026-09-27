// Copyright 2026 Ting. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
)

func TestUnitProxyChecksMgmtServer(t *testing.T) {
	var upstreamHits atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamHits.Add(1)
		if r.URL.Path != "/api/v1/chat/completions" && r.URL.Path != "/api/v1/messages" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("authorization = %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(upstream.Close)

	var mgmtHits atomic.Int32
	var mgmtBody string
	mgmt := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mgmtHits.Add(1)
		body, _ := io.ReadAll(r.Body)
		mgmtBody = string(body)
		if r.URL.Path == "/deny" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(mgmt.Close)

	base, err := url.Parse(upstream.URL)
	if err != nil {
		t.Fatal(err)
	}

	allow := httptest.NewServer(newHandler(base, "test-key", mgmt.URL))
	t.Cleanup(allow.Close)
	deny := httptest.NewServer(newHandler(base, "test-key", mgmt.URL+"/deny"))
	t.Cleanup(deny.Close)

	okReq, err := http.NewRequest(http.MethodPost, allow.URL+"/api/v1/chat/completions", strings.NewReader(`{"model":"x"}`))
	if err != nil {
		t.Fatal(err)
	}
	okReq.Header.Set("User-Agent", "pi (linux 7.0.12-linuxkit; arm64)")
	okReq.Header.Set("X-Api-Key", "run-123")
	okReq.Header.Set("Content-Type", "application/json")
	okRes, err := http.DefaultClient.Do(okReq)
	if err != nil {
		t.Fatal(err)
	}
	okRes.Body.Close()
	if okRes.StatusCode != http.StatusOK {
		t.Fatalf("allowed status = %d", okRes.StatusCode)
	}
	if mgmtBody != `{"runId":"run-123"}` {
		t.Fatalf("mgmt body = %s", mgmtBody)
	}

	msgReq, err := http.NewRequest(http.MethodPost, allow.URL+"/api/v1/messages?beta=true", strings.NewReader(`{"model":"x"}`))
	if err != nil {
		t.Fatal(err)
	}
	msgReq.Header.Set("User-Agent", "pi (linux 7.0.12-linuxkit; arm64)")
	msgReq.Header.Set("X-Api-Key", "run-123")
	msgRes, err := http.DefaultClient.Do(msgReq)
	if err != nil {
		t.Fatal(err)
	}
	msgRes.Body.Close()
	if msgRes.StatusCode != http.StatusOK {
		t.Fatalf("messages status = %d", msgRes.StatusCode)
	}

	denyReq, err := http.NewRequest(http.MethodPost, deny.URL+"/api/v1/chat/completions", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	denyReq.Header.Set("User-Agent", "pi (linux 7.0.12-linuxkit; arm64)")
	denyReq.Header.Set("X-Api-Key", "run-123")
	denied, err := http.DefaultClient.Do(denyReq)
	if err != nil {
		t.Fatal(err)
	}
	denied.Body.Close()
	if denied.StatusCode != http.StatusForbidden {
		t.Fatalf("denied status = %d", denied.StatusCode)
	}

	other, err := http.Post(allow.URL+"/api/v1/key", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	other.Body.Close()
	if other.StatusCode != http.StatusForbidden {
		t.Fatalf("other path status = %d", other.StatusCode)
	}

	if upstreamHits.Load() != 2 {
		t.Fatalf("upstream hits = %d", upstreamHits.Load())
	}
	if mgmtHits.Load() != 3 {
		t.Fatalf("mgmt hits = %d", mgmtHits.Load())
	}
}
