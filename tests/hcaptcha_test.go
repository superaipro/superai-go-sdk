package tests

import (
	"context"
	"fmt"
	"log"
	"superaipro"
	"testing"
	"time"
)

func Test_hcaptcha(t *testing.T) {
	// apiKey
	apiKey := "aa37d476-c06a-4a8f-bb94-xxxxx" // your apikey
	client := superaipro.NewClient(apiKey)

	// userInfo
	userInfo, err := client.GetUser()
	if err != nil {
		log.Fatalln("GetUser error", err)
		return
	}
	//wallet
	fmt.Println(userInfo.Data)

	//hcaptcha
	capData := superaipro.HCaptcha{
		SiteKey: "51829642-2cda-4b09-896c-594f89d700cc",
		Url:     "http://democaptcha.com/demo-form-eng/hcaptcha.html",
		Type:    "HCaptchaV1",
		Timeout: 120,
	}

	req := capData.ToRequest()
	req.SetProxy("HTTPS", "xxx.xx.xxx.xxx:8080:username:password") // your proxy
	req.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36")
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)
	defer cancel()
	taskStatus, err := client.Solve(ctx, req)
	if err != nil {
		return
	}
	fmt.Println(taskStatus.Token, err)

}
func Test_hcaptcha_identity(t *testing.T) {
	// apiKey
	apiKey := "aa37d476-c06a-4a8f-bb94-xxxxx" // your apikey
	client := superaipro.NewClient(apiKey)
	//identify
	identifyData := superaipro.Identify{
		Type: "HCaptchaV1",
		Images: []string{
			"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",
		},
	}
	identifyResult, err := client.Identify(identifyData.ToRequest())
	if err != nil {
		log.Fatalln("Identify error", err)
		return
	}
	fmt.Println(identifyResult.Data)

}
