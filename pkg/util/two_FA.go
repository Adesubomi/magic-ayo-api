package util

import (
	"fmt"
	"github.com/xlzd/gotp"
	"net/url"
	"time"
)

type Google2faType string

const (
	Google2faTypeTimeBased    Google2faType = "totp"
	Google2faTypeCounterBased Google2faType = "hotp"
)

type Google2faConfig struct {
	Type         Google2faType
	Secret       string
	InitialCount int
}

func Generate2faSecret() *Google2faConfig {
	secret := gotp.RandomSecret(16)
	return &Google2faConfig{
		Type:   Google2faTypeTimeBased,
		Secret: secret,
	}
}

func Google2faConfigFromSecret(secret string, baseType Google2faType) *Google2faConfig {
	return &Google2faConfig{
		Type:   baseType,
		Secret: secret,
	}
}

func (g *Google2faConfig) SetInitialCount(count int) {
	g.InitialCount = count
}

func (g *Google2faConfig) GetProvisioningUrl(issuer, username string) string {
	urlSample := fmt.Sprintf("otpauth://totp/%v?issuer=%v&secret=%v",
		url.QueryEscape(username),
		url.QueryEscape(issuer),
		g.Secret)
	if g.Type == Google2faTypeCounterBased {
		urlSample = fmt.Sprintf("otpauth://totp/%v?issuer=%v&secret=%v&counter=%v",
			url.QueryEscape(username),
			url.QueryEscape(issuer),
			g.Secret,
			g.InitialCount)
	}

	return urlSample
}

func (g *Google2faConfig) Authenticate(token string) (bool, error) {
	valid := gotp.NewDefaultTOTP(g.Secret).Verify(token, time.Now().Unix())
	return valid, nil
}
