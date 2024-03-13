package notion

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) PostRequest(url string, payload io.Reader, response any) error {
	ctx, cancelContext := context.WithTimeout(context.Background(), requestTimeout)
	defer cancelContext()

	req, err := http.NewRequestWithContext(ctx, "POST", url, payload)
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", "token_v2="+c.options.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, response)
	return err
}
