package mb_broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client http.Client
	url    string
	token  string
}

func New(uri, token string) *Client {
	return &Client{
		client: http.Client{},
		url:    uri,
		token:  token,
	}
}

type FetchProductCodesReqSt struct {
	PrvCode string `json:"prv_code"`
}

func (c *Client) FetchCompletedOrders() {

}

func (c *Client) FetchProductCodesByOrder(orderID string) ([]string, error) {
	endpoint := fmt.Sprintf("%s/ord/product_codes", c.url)

	payload := &FetchProductCodesReqSt{
		PrvCode: orderID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", res.StatusCode, string(body))
	}

	var productCodes []string
	if err := json.Unmarshal(body, &productCodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return productCodes, nil
}
