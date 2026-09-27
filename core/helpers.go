// Copyright 2026 Ting. All rights reserved.
// License can be found in the LICENSE file.

package core

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 5 * time.Second}

func isAllowed(r *http.Request, mgmtURL string) bool {
	if !strings.Contains(strings.ToLower(r.Header.Get("User-Agent")), "pi") {
		return false
	}

	body, err := json.Marshal(map[string]string{
		"runId": r.Header.Get("X-Api-Key"),
	})
	if err != nil {
		return false
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, mgmtURL, bytes.NewReader(body))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	_, _ = io.Copy(io.Discard, res.Body)

	return res.StatusCode == http.StatusOK
}
