package superaipro

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	BaseURL          = "https://api.superai.pro"
	taskStateSuccess = "success"
	taskStateIdle    = "idle"
	taskStateFailed  = "failed"
)

var (
	ErrNetwork          = errors.New("superai: Network failure")
	ErrApi              = errors.New("superai: API error")
	ErrNotAuthorization = errors.New("superai: request unauthorized")
	ResErr              = "api error StatusCode:"
)

func NewClient(apiKey string) IClient {
	base, _ := url.Parse(BaseURL)
	return &Client{
		BaseURL:          base,
		ApiKey:           apiKey,
		DefaultTimeout:   120,
		PollingInterval:  10,
		RecaptchaTimeout: 600,
		httpClient:       &http.Client{},
	}
}

func (c *Client) res(path string, req Request) (string, error) {

	rel := &url.URL{Path: path}
	uri := c.BaseURL.ResolveReference(rel)

	c.httpClient.Timeout = time.Duration(c.DefaultTimeout) * time.Second

	var resp *http.Response = nil

	values := url.Values{}
	for key, val := range req.Params {
		values.Add(key, val.(string))
	}
	uri.RawQuery = values.Encode()

	var err error = nil
	resp, err = http.Get(uri.String())
	if err != nil {
		return "", ErrNetwork
	}

	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	data := body.String()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("%v%v", ResErr, resp.StatusCode))
	}

	return data, nil
}

func (c *Client) req(ctx context.Context, path string, req Request) (string, error) {
	rel := &url.URL{Path: path}
	uri := c.BaseURL.ResolveReference(rel)

	c.httpClient.Timeout = time.Duration(c.DefaultTimeout) * time.Second

	dataByte, _ := json.Marshal(req.Params)
	request, err := http.NewRequest("POST", uri.String(), bytes.NewReader(dataByte))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("authorization", c.ApiKey)
	if ctx != nil {
		request = request.WithContext(ctx)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return "", ErrNetwork
	}

	defer resp.Body.Close()
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	data := body.String()

	if resp.StatusCode == 401 {
		return "", ErrNotAuthorization
	}
	if resp.StatusCode != http.StatusOK {
		return "", ErrApi
	}

	return data, nil
}

func (c *Client) Send(ctx context.Context, req Request) (result Result, err error) {
	data, err := c.req(ctx, "/api/v1/task", req)
	if err != nil {
		return result, err
	}
	if err = json.Unmarshal([]byte(data), &result); err != nil {
		return result, err
	}
	return result, err
}

func (c *Client) Solve(ctx context.Context, req Request) (result Result, err error) {
	result, err = c.Send(ctx, req)
	if err != nil {
		return result, err
	}
	retryNum := 0
	for {
		select {
		case <-time.After(time.Second * 3):
			break
		case <-ctx.Done():
			return result, ctx.Err()
		}
		result, err = c.GetResult(ctx, result.ID) // get Task
		if err != nil {
			httpCode := c.getHttpStatusCode(err.Error())
			if httpCode != 0 && httpCode > 500 && retryNum <= 3 {
				retryNum += 1
				continue
			}
			return result, err
		}
		if result.Status == taskStateSuccess || result.Status == taskStateFailed {
			break
		}
	}
	return result, nil
}

func (c *Client) GetResult(ctx context.Context, taskId string) (result Result, err error) {
	req := Request{
		Params: map[string]interface{}{
			"id":     taskId,
			"apiKey": c.ApiKey,
		},
	}

	data, err := c.res("/api/v1/task", req)
	if err != nil {
		return result, err
	}

	if err = json.Unmarshal([]byte(data), &result); err != nil {
		return result, err
	}
	return result, nil
}

func (c *Client) GetUser() (user User, err error) {
	req := Request{
		Params: map[string]interface{}{
			"apiKey": c.ApiKey,
		},
	}
	data, err := c.res("/api/v1/user", req)
	if err != nil {
		return user, err
	}

	if err = json.Unmarshal([]byte(data), &user); err != nil {
		return user, err
	}

	return user, nil
}

func (c *Client) Identify(req Request) (identifyResult IdentifyResult, err error) {
	data, err := c.req(context.TODO(), "/api/v1/identify", req)
	if err != nil {
		return identifyResult, err
	}
	if err = json.Unmarshal([]byte(data), &identifyResult); err != nil {
		return identifyResult, err
	}
	return identifyResult, err
}

func (c *Client) GetWallet() (wallet Wallet, err error) {
	req := Request{
		Params: map[string]interface{}{
			"apiKey": c.ApiKey,
		},
	}
	data, err := c.res("/api/v1/wallet", req)
	if err != nil {
		return wallet, err
	}
	if err = json.Unmarshal([]byte(data), &wallet); err != nil {
		return wallet, err
	}
	return wallet, nil
}

func (c *Client) SetTimeout(timeout int) {
	c.DefaultTimeout = timeout
}

func (c *Client) getHttpStatusCode(errString string) int {
	if !strings.Contains(errString, ResErr) {
		return 0
	}
	code, _ := strconv.Atoi(strings.Replace(errString, ResErr, "", 1))
	return code
}
