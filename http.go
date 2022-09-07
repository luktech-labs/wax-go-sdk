package wax

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type httpWrapper struct {
	httpClient *http.Client
}

func (hw *httpWrapper) Get(ctx context.Context, url string, dest interface{}) error {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}
	r.Close = true

	res, err := hw.httpClient.Do(r)
	if err != nil {
		return errors.Wrap(err, "performing request")
	}
	defer res.Body.Close()
	defer io.Copy(io.Discard, res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "reading response body")
	}

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("Got unexpected status code: %d. Body: %hw", res.StatusCode, string(b))
	}

	err = json.Unmarshal(b, dest)
	if err != nil {
		return errors.Wrap(err, "unmarshalling response body")
	}

	return nil
}

func (hw *httpWrapper) Post(ctx context.Context, url string, body []byte, dest interface{}) error {
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	r.Header.Add("Content-type", "application/json")

	res, err := hw.httpClient.Do(r)
	if err != nil {
		return errors.Wrap(err, "performing request")
	}
	defer res.Body.Close()
	defer io.Copy(io.Discard, res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "reading response body")
	}

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("Got unexpected status code: %d. Body: %hw", res.StatusCode, string(b))
	}

	err = json.Unmarshal(b, dest)
	if err != nil {
		return errors.Wrap(err, "unmarshalling response body")
	}

	return nil
}
