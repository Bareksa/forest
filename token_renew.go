package begundal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type renewToken struct {
	Token     string `json:"token"`
	Increment string `json:"increment,omitempty"`
}

// RenewTokenSelf attempts to renew token registered to self. Cannot renew root token with 0 time to live (never expire).
func (v *Vault) RenewTokenSelf(ctx context.Context) (err error) {
	body, err := json.Marshal(renewToken{
		Token: v.Config.Token,
	})
	if err != nil {
		return
	}
	req, err := v.requestGen(ctx, http.MethodPost, "/auth/token/renew", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if err := checkErrorResponse(res); err != nil {
		return err
	}
	return
}

// RenewTokenOther attempts to renew passed token using self's token
func (v *Vault) RenewTokenOther(ctx context.Context, token string) (err error) {
	body, err := json.Marshal(renewToken{
		Token: token,
	})
	if err != nil {
		return
	}
	req, err := v.requestGen(ctx, http.MethodPost, "/auth/token/renew", bytes.NewBuffer(body))
	if err != nil {
		return
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if err := checkErrorResponse(res); err != nil {
		return err
	}
	return
}

// RenewTokenOverride attempts to renew passed token with passed token as auth
func (v *Vault) RenewTokenOverride(ctx context.Context, token string) (err error) {
	body, err := json.Marshal(renewToken{
		Token: token,
	})
	if err != nil {
		return
	}
	req, err := v.requestGenOverride(ctx, http.MethodPost, "/auth/token/renew", token, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if err := checkErrorResponse(res); err != nil {
		return err
	}
	return
}
