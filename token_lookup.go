package begundal

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// LookupToken token lookup
type LookupToken struct {
	Data struct {
		Accessor       string            `json:"accessor"`
		CreationTime   int64             `json:"creation_time"`
		CreationTTL    int64             `json:"creation_ttl"`
		DisplayName    string            `json:"display_name"`
		EntityID       string            `json:"entity_id"`
		ExpireTime     interface{}       `json:"expire_time"`
		ExplicitMaxTTL int64             `json:"explicit_max_ttl"`
		ID             string            `json:"id"` // This is the token auth
		Meta           map[string]string `json:"meta"`
		NumUses        int64             `json:"num_uses"`
		Orphan         bool              `json:"orphan"`
		Path           string            `json:"path"`
		Policies       []string          `json:"policies"`
		TTL            int64             `json:"ttl"`
		Type           string            `json:"type"`
	} `json:"data"`
	RequestID     string      `json:"request_id"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int64       `json:"lease_duration"`
	WrapInfo      interface{} `json:"wrap_info"`
	Warnings      []string    `json:"warnings"`
	Auth          interface{} `json:"auth"`
}

type sLookupOther struct {
	Token string `json:"token"`
}

// LookupSelf lookup information on current token used in the instance
func (v *Vault) LookupSelf(ctx context.Context) (lookup LookupToken, err error) {
	req, err := v.requestGen(ctx, http.MethodGet, "/auth/token/lookup-self", nil)
	if err != nil {
		return
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if err := checkErrorResponse(res); err != nil {
		return lookup, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return lookup, err
	}
	err = json.Unmarshal(body, &lookup)
	return lookup, nil
}

// LookupOther lookup information on passed token
func (v *Vault) LookupOther(ctx context.Context, token string) (lookup LookupToken, err error) {
	body := sLookupOther{
		Token: token,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return
	}
	req, err := v.requestGen(ctx, http.MethodPost, "/auth/token/lookup", bytes.NewBuffer(b))
	if err != nil {
		return
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if err := checkErrorResponse(res); err != nil {
		return lookup, err
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return lookup, err
	}
	err = json.Unmarshal(resBody, &lookup)
	return lookup, nil
}
