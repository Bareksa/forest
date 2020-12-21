package forest

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// CreateTokenInstance struct to create token
type CreateTokenInstance struct {
	ID              string            `json:"id,omitempty"`
	RoleName        string            `json:"role_name,omitempty"`
	Policies        []string          `json:"policies,omitempty"`
	Meta            map[string]string `json:"meta,omitempty"`
	NoParent        bool              `json:"no_parent,omitempty"`
	NoDefaultPolicy bool              `json:"no_default_policy,omitempty"`
	Renewable       bool              `json:"renewable,omitempty"`
	TTL             string            `json:"ttl,omitempty"`
	Type            string            `json:"type,omitempty"`
	ExplicitMaxTTL  string            `json:"explicit_max_ttl,omitempty"`
	DisplayName     string            `json:"display_name,omitempty"`
	NumUses         int               `json:"num_uses,omitempty"`
	EntityAlias     string            `json:"entity_alias,omitempty"`
	Period          string            `json:"period,omitempty"`
	authToken       string
	baseURL         string
	httpClient      *http.Client
}

type tokenResponse struct {
	RequestID     string      `json:"request_id"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int64       `json:"lease_duration"`
	Data          interface{} `json:"data"`
	WrapInfo      interface{} `json:"wrap_info"`
	Warnings      []string    `json:"warnings"`
	Auth          struct {
		ClientToken   string            `json:"client_token"`
		Accessor      string            `json:"accessor"`
		Policies      []string          `json:"policies"`
		TokenPolicies []string          `json:"token_policies"`
		Metadata      map[string]string `json:"metadata"`
		LeaseDuration int64             `json:"lease_duration"`
		Renewable     bool              `json:"renewable"`
		EntityID      string            `json:"entity_id"`
		TokenType     string            `json:"token_type"`
		Orphan        bool              `json:"orphan"`
	} `json:"auth"`
}

// CreateNewToken creates an instance of Token override. Call ``.Do(ctx)`` on the instance to actually create new token.
// Default value for instance follows vault token documentation,
// on here: https://www.vaultproject.io/api/auth/token#parameters
//
// Which means the new token will be renewable by default and has display name of 'token'
func (v *Vault) CreateNewToken() *CreateTokenInstance {
	return &CreateTokenInstance{
		Renewable:   true,
		DisplayName: "token",
		authToken:   v.Config.Token,
		baseURL:     v.BaseURL,
		httpClient:  v.Config.HTTPClient,
	}
}

// WithID replaces instance ID. Make sure there's no character '.' in the argument string
func (c *CreateTokenInstance) WithID(id string) *CreateTokenInstance {
	c.ID = id
	return c
}

// WithRoleName replaces instance rolename.
func (c *CreateTokenInstance) WithRoleName(rolename string) *CreateTokenInstance {
	c.RoleName = rolename
	return c
}

// WithPolicies replaces instance policies
func (c *CreateTokenInstance) WithPolicies(policies ...string) *CreateTokenInstance {
	c.Policies = policies
	return c
}

// WithMeta replaces instance metadata
func (c *CreateTokenInstance) WithMeta(metadata map[string]string) *CreateTokenInstance {
	c.Meta = metadata
	return c
}

// WithCurrentTokenAsParent Uses instance token as parent. Default true.
// If using DoOverride(token), the override token will be set as parent instead.
func (c *CreateTokenInstance) WithCurrentTokenAsParent(b bool) *CreateTokenInstance {
	c.NoParent = !b
	return c
}

// WithDefaultPolicy replaces instance default policy. by default true
func (c *CreateTokenInstance) WithDefaultPolicy(b bool) *CreateTokenInstance {
	c.NoDefaultPolicy = !b
	return c
}

// WithRenewableStatus replaces instance renewable. default true.
func (c *CreateTokenInstance) WithRenewableStatus(b bool) *CreateTokenInstance {
	c.Renewable = b
	return c
}

// WithTimeToLive replaces instance time to live, which by default will depends on Vault's default lease TTL
// If this method is called, duration will be rounded down to the nearest hour argument passed with minimum value of 1 hour
// Only hourly is supported in this package.
func (c *CreateTokenInstance) WithTimeToLive(t time.Duration) *CreateTokenInstance {
	h := t / time.Hour
	if h == 0 {
		c.TTL = "1h"
	} else {
		str := strconv.FormatInt(h.Nanoseconds(), 10)
		c.TTL = str + "h"
	}
	return c
}

// WithPeriod replaces token period. Token that is not renewed in this set of time cannot be renewed again.
// By default, if unset will follow's Vault's default lease TTL.
// Has minimum value of 1 hour. Only hourly is supported in this package.
func (c *CreateTokenInstance) WithPeriod(t time.Duration) *CreateTokenInstance {
	h := t / time.Hour
	if h == 0 {
		c.Period = "1h"
	} else {
		str := strconv.FormatInt(h.Nanoseconds(), 10)
		c.Period = str + "h"
	}
	return c
}

// WithSetAsBatchToken replaces instance token type
func (c *CreateTokenInstance) WithSetAsBatchToken() *CreateTokenInstance {
	c.Type = "batch"
	return c
}

// WithSetAsServiceToken replaces instance token type
func (c *CreateTokenInstance) WithSetAsServiceToken() *CreateTokenInstance {
	c.Type = "service"
	return c
}

// WithExplicitMaxTTL replaces instance explicit time to live, which by default will depends on Vault's default lease TTL
// If this method is called, duration will be rounded down to the nearest hour argument passed with minimum value of 1 hour
func (c *CreateTokenInstance) WithExplicitMaxTTL(t time.Duration) *CreateTokenInstance {
	h := t / time.Hour
	if h == 0 {
		c.ExplicitMaxTTL = "1h"
	} else {
		str := strconv.FormatInt(h.Nanoseconds(), 10)
		c.ExplicitMaxTTL = str + "h"
	}
	return c
}

// WithDisplayName replaces token display name
func (c *CreateTokenInstance) WithDisplayName(s string) *CreateTokenInstance {
	c.DisplayName = s
	return c
}

// WithNumberOfUses replaces token's allowed number of usage. Signing in to vault using UI with the token is considered used one time. By default 0, which is infinite.
func (c *CreateTokenInstance) WithNumberOfUses(i int) *CreateTokenInstance {
	c.NumUses = i
	return c
}

// WithEntityAlias replaces instance entity alias. MUST BE USED alongside WithRoleName and Role name must exist within vault.
func (c *CreateTokenInstance) WithEntityAlias(s string) *CreateTokenInstance {
	c.EntityAlias = s
	return c
}

// Do creates the token.
func (c *CreateTokenInstance) Do(ctx context.Context) (token string, err error) {
	body, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/auth/token/create", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Vault-Token", c.authToken)
	req.Header.Set("X-Vault-Request", "true")
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if err := checkErrorResponse(res); err != nil {
		return token, err
	}
	response := tokenResponse{}
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}
	token = response.Auth.ClientToken
	return
}

// DoOverride creates the token, with passed token as auth and parent if orphan status is false (or no parent is true)
func (c *CreateTokenInstance) DoOverride(ctx context.Context, authToken string) (token string, err error) {
	body, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/auth/token/create", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Vault-Token", authToken)
	req.Header.Set("X-Vault-Request", "true")
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if err := checkErrorResponse(res); err != nil {
		return token, err
	}
	response := tokenResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}
	token = response.Auth.ClientToken
	return
}
