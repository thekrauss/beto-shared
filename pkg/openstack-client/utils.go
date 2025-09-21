package openstack

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/thekrauss/beto-shared/pkg/errors"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// exécute une requête HTTP vers OpenStack
func doRequest(ctx context.Context, method, url, token string, body interface{}, out interface{}) error {
	var buf io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return errors.Wrap(err, errors.CodeInternal, "failed to marshal request body")
		}
		buf = bytes.NewBuffer(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to build request")
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("X-Auth-Token", token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "http request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return errors.Newf(errors.CodeInternal, "openstack API error %d: %s", resp.StatusCode, string(b))
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}
