package voximplant

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	client     http.Client
	url        string
	token      string
	domainName string
}

func New(url, token, domainName string) *Client {
	return &Client{
		client:     http.Client{},
		url:        url,
		token:      token,
		domainName: domainName,
	}
}

func (c *Client) SendNotification(clientID, channelID, templateID, headerParam, buttonURLParam string) error {
	endpoint := fmt.Sprintf("%s/api/v3/botService/sendTemplateMessage", c.url)

	params := url.Values{}
	params.Set("domain", c.domainName)
	params.Set("client_id", clientID)
	params.Set("channel_id", channelID)
	params.Set("message_template_id", templateID)
	params.Set("access_token", c.token)
	params.Set("header_param_value", headerParam)
	params.Set("button_url_param_value", buttonURLParam)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, response: %s", res.StatusCode, string(body))
	}

	fmt.Printf("Notification sent successfully: %s\n", string(body))
	return nil
}
