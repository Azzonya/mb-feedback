package mb_broker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	orderModel "mb-feedback/internal/domain/order/model"
	"mb-feedback/internal/errs"
	"net/http"
	"net/url"
)

type Client struct {
	client  http.Client
	baseURL string
	token   string
}

func New(baseURL, token string) *Client {
	return &Client{
		client:  http.Client{},
		baseURL: baseURL,
		token:   token,
	}
}

func (c *Client) FetchCompletedOrders(ctx context.Context) ([]*orderModel.Order, error) {
	endpoint := fmt.Sprintf("%s/ord", c.baseURL)

	repObj := &FetchCompletedOrdersRepSt{}

	//creationTsGte := time.Now().Add(-time.Hour).Format(time.RFC3339)

	statusOk, respBody, err := c.sendRequest(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
		url.Values{
			"prv_id":    {"kaspi"},
			"page_size": {"100"},
			//"creation_ts_gte": {creationTsGte},
			"status": {"COMPLETED"},
		},
		nil,
		&repObj,
		nil)
	if err != nil {
		slog.Error("FetchCompletedOrders", "error", fmt.Errorf("failed to send request: %w", err))
		return nil, err
	}
	if !statusOk {
		slog.Error("FetchCompletedOrders", "statusOk", statusOk, "body", string(respBody))
		return nil, errs.BadStatusCode
	}

	result := make([]*orderModel.Order, 0, len(repObj.Results))
	for _, v := range repObj.Results {
		result = append(result, &orderModel.Order{
			ExternalOrderID: v.PrvCode,
			UserPhone:       v.Customer.CellPhone,
			UserName:        v.Customer.FirstName,
		})
	}

	return result, nil
}

func (c *Client) FetchProductCodes(ctx context.Context, orderID string) ([]string, error) {
	endpoint := fmt.Sprintf("%s/ord/product_codes", c.baseURL)

	var repObj []string

	statusOk, respBody, err := c.sendRequest(
		ctx,
		http.MethodPost,
		endpoint,
		nil,
		nil,
		&FetchProductCodesReqSt{
			PrvCode: orderID,
		},
		&repObj,
		nil)
	if err != nil {
		slog.Error("FetchProductCodesByOrder", "error", fmt.Errorf("failed to send request: %w", err))
		return nil, err
	}
	if !statusOk {
		slog.Error("FetchProductCodesByOrder", "statusOk", statusOk, "body", string(respBody))
		return nil, errs.BadStatusCode
	}

	return repObj, nil
}

func (c *Client) sendRequest(
	ctx context.Context,
	method string,
	url string,
	header http.Header,
	params url.Values,
	reqObj any,
	repObj any,
	errRepObj any) (bool, []byte, error) {

	var reqBody io.Reader
	if reqObj != nil {
		reqJson, err := json.Marshal(reqObj)
		if err != nil {
			return false, nil, fmt.Errorf("json.Marshal: %w", err)
		}
		reqBody = bytes.NewReader(reqJson)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return false, nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	if header != nil {
		req.Header = header
	}

	if reqObj != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("Authorization", "Bearer "+c.token)

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, nil, fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	repBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, nil, fmt.Errorf("resp.Body.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if errRepObj != nil && len(repBody) > 0 {
			_ = json.Unmarshal(repBody, errRepObj)
		}
		return false, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if repObj != nil {
		if err = json.Unmarshal(repBody, repObj); err != nil {
			return false, nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
	}

	return true, repBody, nil
}
