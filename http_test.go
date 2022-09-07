package wax

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	Success    = "\033[32m✓"
	Failed     = "\033[31m✗"
	ColorReset = "\033[0m"
)

func TestCreateRequestErrorMessage(t *testing.T) {
	tcs := []struct {
		name string
		req  *http.Request
		res  *http.Response
		exp  RequestErrorMessage
	}{
		{
			name: "response nil",
			req:  &http.Request{URL: &url.URL{Host: "test.localhost", Scheme: "http"}},
			res:  nil,
			exp:  RequestErrorMessage{RequestURL: "http://test.localhost"},
		},
		{
			name: "response empty body",
			req:  &http.Request{URL: &url.URL{Host: "test.localhost", Scheme: "http"}},
			res:  &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBuffer([]byte{}))},
			exp:  RequestErrorMessage{RequestURL: "http://test.localhost", StatusCode: http.StatusOK},
		},
		{
			name: "response nil body",
			req:  &http.Request{URL: &url.URL{Host: "test.localhost", Scheme: "http"}},
			res:  &http.Response{StatusCode: http.StatusOK, Body: nil},
			exp:  RequestErrorMessage{RequestURL: "http://test.localhost", StatusCode: http.StatusOK},
		},
		{
			name: "response with body",
			req:  &http.Request{URL: &url.URL{Host: "test.localhost", Scheme: "http"}},
			res:  &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(bytes.NewBuffer([]byte("test response")))},
			exp:  RequestErrorMessage{RequestURL: "http://test.localhost", StatusCode: http.StatusCreated, Body: "test response"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := createRequestErrorMessage(tc.req, tc.res)
			if diff := cmp.Diff(got, tc.exp); diff != "" {
				t.Fatalf("%s\t Messages should be identical. Found diff: %s. %s\n", Failed, diff, ColorReset)
			}
			t.Logf("%s\t Messages should be identical. %s\n", Success, ColorReset)
		})
	}
}

func TestProcess(t *testing.T) {
	type fakeDest struct {
		Name string
		Age  int
	}

	server := httptest.NewServer(nil)
	defer server.Close()

	tcs := []struct {
		name       string
		req        *http.Request
		statusCode int
		exp        fakeDest
		expErr     error
	}{
		{
			name:       "no error status ok",
			req:        &http.Request{Method: http.MethodGet, Body: nil},
			statusCode: http.StatusOK,
			exp:        fakeDest{Age: 30, Name: "Alex"},
			expErr:     nil,
		},
		{
			name:       "err status bad request",
			req:        &http.Request{Method: http.MethodGet, Body: nil},
			statusCode: http.StatusBadRequest,
			exp:        fakeDest{},
			expErr: errors.Errorf(
				"Got unexpected status code. Details: %+v",
				RequestErrorMessage{
					RequestURL: server.URL,
					Body:       `{"Name":"","Age":0}`,
					StatusCode: http.StatusBadRequest,
				},
			),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(tc.exp)
			if err != nil {
				t.Fatalf(
					"%s\t Should be able to marshall the request body. Got err: %s. %s\n",
					Failed,
					err.Error(),
					ColorReset,
				)
			}
			handlerFn := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(tc.statusCode)
				res.Write(body)
			})
			server.Config.Handler = handlerFn

			u, err := url.Parse(server.URL)
			if err != nil {
				t.Fatalf(
					"%s\t Should be able to parse the httptest server url. Err: %s. %s\n",
					Failed,
					err.Error(),
					ColorReset,
				)
			}

			tc.req.URL = u
			hw := httpWrapper{httpClient: server.Client()}

			var dest fakeDest
			err = hw.process(tc.req, &dest)
			if !equalError(err, tc.expErr) {
				t.Fatalf(
					"%s\t Should get the expected error. Found: %s, Expected: %s. %s\n",
					Failed,
					err.Error(),
					tc.expErr,
					ColorReset,
				)
			}
			t.Logf("%s\t Should get the expected error.%s\b", Success, ColorReset)

			if diff := cmp.Diff(dest, tc.exp); diff != "" {
				t.Fatalf("%s\t Response dest should be identical. Found diff: %s. %s\n", Failed, diff, ColorReset)
			}
			t.Logf("%s\t Response dest should be identical. %s\n", Success, ColorReset)
		})
	}
}

// equalError reports whether errors a and b are considered equal.
// They're equal if both are nil, or both are not nil and a.Error() == b.Error().
func equalError(a, b error) bool {
	return a == nil && b == nil || a != nil && b != nil && a.Error() == b.Error()
}
