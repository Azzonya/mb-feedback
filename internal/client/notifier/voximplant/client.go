package voximplant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mb-feedback/internal/errs"
	"net/http"
	"net/url"
)

type Client struct {
	client     *http.Client
	baseURL    string
	token      string
	domainName string
	templateID string
	channelID  string
}

type TextParamValues struct {
	Name2 string `json:"name2"`
}

func New(baseURL, token, domainName, templateID, channelID string) *Client {
	return &Client{
		client:     &http.Client{},
		baseURL:    baseURL,
		token:      token,
		domainName: domainName,
		templateID: templateID,
		channelID:  channelID,
	}
}

func (c *Client) SendNotification(ctx context.Context, orderID, userPhone, userName, productCode string) error {
	endpoint := fmt.Sprintf("%s/api/v3/botService/sendTemplateMessage", c.baseURL)

	buttonUrlParam := fmt.Sprintf("orderCode=%s&productCode=%s&rating=5", orderID, productCode)

	data, err := json.MarshalIndent(&TextParamValues{
		Name2: userName,
	}, "", "  ")
	if err != nil {
		slog.Error("Marshal text_param values error:", "error", err)
		return err
	}

	params := url.Values{
		"domain":                 {c.domainName},
		"client_id":              {userPhone},
		"message_template_id":    {c.templateID},
		"channel_id":             {c.channelID},
		"access_token":           {c.token},
		"header_param_value":     {userName},
		"button_url_param_value": {buttonUrlParam},
		"text_param_values":      {string(data)},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBufferString(params.Encode()))
	if err != nil {
		slog.Error("NewRequestWithContext error:", "error", err)
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.client.Do(req)
	if err != nil {
		slog.Error("Do request error:", "error", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Read response body error:", "error", err)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		slog.Error("unexpected status code: %d, response: %s", "statusCode", res.StatusCode, "respBody", string(body))
		return errs.BadStatusCode
	}

	fmt.Printf("Notification sent successfully: %s\n", string(body))
	return nil
}
