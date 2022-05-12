package superaipro

import (
	"net/http"
	"net/url"
	"strconv"
)

type (
	Request struct {
		Params map[string]interface{}
	}

	Client struct {
		BaseURL          *url.URL
		ApiKey           string
		DefaultTimeout   int
		RecaptchaTimeout int
		PollingInterval  int
		httpClient       *http.Client
	}

	Identify struct {
		Type   string
		Images []string
	}

	GeeTest struct {
		GT        string
		Challenge string
		Url       string
		ApiServer string
	}

	HCaptcha struct {
		SiteKey string
		Url     string
		Type    string
		Timeout int
	}

	ReCaptcha struct {
		SiteKey   string
		Url       string
		Invisible bool
		Version   string
		Action    string
		Score     float64
	}

	IdentifyData struct {
		ProjectName string   `json:"projectName" validate:"required"`
		Image       []string `json:"image"`
	}

	Result struct {
		ID         string `json:"id"`
		Token      string `json:"token"`
		Success    bool   `json:"success"`
		Status     string `json:"status"`
		Message    string `json:"message"`
		ExpireTime int64  `json:"expireTime"`
	}

	IdentifyResult struct {
		Data    string `json:"data"`
		Success bool   `json:"success"`
	}

	User struct {
		Data struct {
			Email  string `json:"email"`
			Wallet struct {
				Balance      int `json:"balance"`
				NumberOfCall int `json:"numberOfCall"`
			} `json:"wallet"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	Wallet struct {
		Data struct {
			Balance      int `json:"balance"`
			NumberOfCall int `json:"numberOfCall"`
		} `json:"data"`
		Success bool `json:"success"`
	}
)

func (req *Request) SetProxy(proxyType string, uri string) {
	req.Params["proxytype"] = proxyType
	req.Params["proxy"] = uri
}

func (req *Request) SetUserAgent(userAgent string) {
	req.Params["userAgent"] = userAgent
}

func (req *Request) SetCookies(cookies string) {
	req.Params["cookies"] = cookies
}

func (c *GeeTest) ToRequest() Request {
	req := Request{
		Params: map[string]interface{}{"projectName": "geetest"},
	}
	if c.GT != "" {
		req.Params["gt"] = c.GT
	}
	if c.Challenge != "" {
		req.Params["challenge"] = c.Challenge
	}
	if c.Url != "" {
		req.Params["pageurl"] = c.Url
	}
	if c.ApiServer != "" {
		req.Params["api_server"] = c.ApiServer
	}

	return req
}

func (c *HCaptcha) ToRequest() Request {
	req := Request{
		Params: map[string]interface{}{"projectName": "hcaptcha"},
	}
	if c.SiteKey != "" {
		req.Params["siteKey"] = c.SiteKey
	}
	if c.Url != "" {
		req.Params["pageUrl"] = c.Url
	}
	if c.Type != "" {
		req.Params["type"] = c.Type
	}
	if c.Timeout != 0 {
		req.Params["timeout"] = c.Timeout
	}
	return req
}

func (c *ReCaptcha) ToRequest() Request {
	req := Request{
		Params: map[string]interface{}{"projectName": "userrecaptcha"},
	}
	if c.SiteKey != "" {
		req.Params["siteKey"] = c.SiteKey
	}
	if c.Url != "" {
		req.Params["siteUrl"] = c.Url
	}
	if c.Invisible {
		req.Params["invisible"] = "1"
	}
	if c.Version != "" {
		req.Params["version"] = c.Version
	}
	if c.Action != "" {
		req.Params["action"] = c.Action
	}
	if c.Score != 0 {
		req.Params["min_score"] = strconv.FormatFloat(c.Score, 'f', -1, 64)
	}

	return req
}

func (c *Identify) ToRequest() Request {
	req := Request{
		Params: map[string]interface{}{"projectName": "userrecaptcha"},
	}
	if c.Type != "" {
		req.Params["type"] = c.Type
	}
	if len(c.Images) > 0 {
		req.Params["image"] = c.Images
	}

	return req
}
