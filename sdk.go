package wax

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	getTableRowsEndpoint = "/v1/chain/get_table_rows"
	getInfoEndpoint      = "/v1/chain/get_info"
)

var ErrNoTableRows = errors.New("No table rows found")

type Sdk struct {
	nodeURL     string
	httpWrapper httpWrapper
}

func NewSdk(nodeURL string, httpOpts ...HttpOption) *Sdk {
	httpClient := http.DefaultClient
	for _, opt := range httpOpts {
		opt(httpClient)
	}

	return &Sdk{nodeURL: strings.TrimSuffix(nodeURL, "/"), httpWrapper: httpWrapper{httpClient: httpClient}}
}

func (s *Sdk) GetTableRowsContext(ctx context.Context, payload GetTableRowsPayload, dest interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "marshalling payload")
	}

	var res struct {
		Rows []interface{} `json:"rows"`
	}

	url := s.nodeURL + getTableRowsEndpoint
	err = s.httpWrapper.Post(ctx, url, b, &res)
	if err != nil {
		return errors.Wrap(err, "calling post")
	}

	if len(res.Rows) == 0 {
		return ErrNoTableRows
	}

	b, err = json.Marshal(res.Rows)
	if err != nil {
		return errors.Wrap(err, "marshalling rows")
	}

	err = json.Unmarshal(b, dest)
	if err != nil {
		return errors.Wrap(err, "unmarshalling rows")
	}

	return nil
}

func (s *Sdk) GetTableRows(payload GetTableRowsPayload, dest interface{}) error {
	return s.GetTableRowsContext(context.Background(), payload, dest)
}

func (s *Sdk) GetInfoContext(ctx context.Context) (GetInfoResponse, error) {
	var res GetInfoResponse

	url := s.nodeURL + getInfoEndpoint
	err := s.httpWrapper.Get(ctx, url, &res)
	if err != nil {
		return GetInfoResponse{}, errors.Wrap(err, "calling post")
	}

	return res, nil
}

func (s *Sdk) GetInfo() (GetInfoResponse, error) {
	return s.GetInfoContext(context.Background())
}
