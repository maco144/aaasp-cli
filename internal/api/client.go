package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.Status, e.Message)
}

func (c *Client) do(method, path string, body any, out any) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+"/v1"+path, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var errBody struct {
			Error   string `json:"error"`
			Message string `json:"message"`
			Errors  any    `json:"errors"`
		}
		_ = json.Unmarshal(respData, &errBody)
		msg := errBody.Error
		if msg == "" {
			msg = string(respData)
		}
		// Some errors pair a machine code with a tenant-facing explanation
		// (e.g. {"error":"not_runnable","message":"..."}); show both.
		if errBody.Message != "" {
			msg = msg + ": " + errBody.Message
		}
		return &APIError{Status: resp.StatusCode, Message: msg}
	}

	if out != nil && len(respData) > 0 {
		return json.Unmarshal(respData, out)
	}
	return nil
}

func (c *Client) Get(path string, out any) error {
	return c.do("GET", path, nil, out)
}

func (c *Client) Post(path string, body, out any) error {
	return c.do("POST", path, body, out)
}

func (c *Client) Delete(path string) error {
	return c.do("DELETE", path, nil, nil)
}

// Register is unauthenticated — hits /v1/register directly
func (c *Client) Register(email, password string) (map[string]any, error) {
	data, err := json.Marshal(map[string]string{"email": email, "password": password})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.baseURL+"/v1/register", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		var errBody struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(respData, &errBody)
		return nil, &APIError{Status: resp.StatusCode, Message: errBody.Error}
	}

	var out map[string]any
	if err := json.Unmarshal(respData, &out); err != nil {
		return nil, err
	}
	return out, nil
}
