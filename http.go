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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	req.Header.Add("Content-type", "application/json")

	err = hw.process(req, dest)
	if err != nil {
		return errors.Wrap(err, "processing request")
	}

	return nil
}

func (hw *httpWrapper) Post(ctx context.Context, url string, payload []byte, dest interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	req.Header.Add("Content-type", "application/json")

	err = hw.process(req, dest)
	if err != nil {
		return errors.Wrap(err, "processing request")
	}

	return nil
}

// process will use the stored http.Client to send the request, checks if the status code
// is an unexpected one and unmarshall the response body to the destination.
func (hw *httpWrapper) process(req *http.Request, dest interface{}) error {
	res, err := hw.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "sending request")
	}
	defer res.Body.Close()
	defer io.Copy(io.Discard, res.Body)

	if res.StatusCode <= 199 || res.StatusCode >= 300 {
		return errors.Errorf("Got unexpected status code. Details: %+v", createRequestErrorMessage(req, res))
	}

	err = unmarshallResponseBody(res.Body, dest)
	if err != nil {
		return errors.Wrap(err, "unmarshalling response")
	}

	return nil
}

type RequestErrorMessage struct {
	RequestURL string
	StatusCode int
	Body       string
}

func createRequestErrorMessage(req *http.Request, res *http.Response) RequestErrorMessage {
	message := RequestErrorMessage{
		RequestURL: req.URL.String(),
	}

	if res != nil {
		if res.Body != nil {
			body, _ := io.ReadAll(res.Body)
			message.Body = string(body)
		}

		message.StatusCode = res.StatusCode
	}

	return message
}

func unmarshallResponseBody(body io.Reader, dest interface{}) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return errors.Wrap(err, "reading response body")
	}

	err = json.Unmarshal(b, dest)
	if err != nil {
		return errors.Wrap(err, "unmarshalling response payload")
	}

	return nil
}
