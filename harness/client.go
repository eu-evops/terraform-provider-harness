package harness

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client struct {
	apiKey   string
	endpoint string
}

func NewClient(apiKey string, endpoint string) *Client {
	return &Client{
		apiKey:   apiKey,
		endpoint: endpoint,
	}
}

type GraphQLQuery struct {
	OperationName string      `json:"operationName,omitempty"`
	Query         string      `json:"query"`
	Variables     interface{} `json:"variables"`
}

type Error struct {
	Message string   `json:"message"`
	Path    []string `json:"path`
}

func (h *Client) query(q *GraphQLQuery, response interface{}) error {
	queryBytes, err := json.Marshal(q)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", h.endpoint, bytes.NewBuffer(queryBytes))
	req.Header.Set("x-api-key", h.apiKey)
	req.Header.Set("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return err
	}

	return nil
}
