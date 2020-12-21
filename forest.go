package forest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Vault struct
type Vault struct {
	BaseURL string
	Config  Config
}

var (
	defaultHost           = "http://localhost:8200"
	defaultHTTPClient     = http.DefaultClient
	defaultTransitEngine  = "transit"
	defaultKeyValueEngine = "kv"
)

// NewClient creates a new vault instance
func NewClient(token string, opts ...OptionFunc) (*Vault, error) {
	if token == "" {
		return nil, errors.New("Empty token")
	}
	config := &Config{
		Token:           token,
		Host:            defaultHost,
		HTTPClient:      defaultHTTPClient,
		VaultAPIVersion: V1,
		KeyValueEngine:  defaultKeyValueEngine,
		TransitEngine:   defaultTransitEngine,
	}

	for _, opt := range opts {
		opt(config)
	}

	z := &Vault{
		BaseURL: fmt.Sprintf("%s/%s", config.Host, config.VaultAPIVersion),
		Config:  *config,
	}

	return z, nil
}

// GetKVEngine returns the engine name used for the instance
func (v *Vault) GetKVEngine() string {
	return v.Config.KeyValueEngine
}

// GetTransitEngine returns the engine name used for the instance
func (v *Vault) GetTransitEngine() string {
	return v.Config.TransitEngine
}

func (v *Vault) requestGen(ctx context.Context, method string, path string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, method, v.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", v.Config.Token)
	req.Header.Set("X-Vault-Request", "true")
	req.Header.Set("Content-Type", "application/json")
	return
}

func (v *Vault) requestGenOverride(ctx context.Context, method, path, token string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, method, v.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", token)
	req.Header.Set("X-Vault-Request", "true")
	req.Header.Set("Content-Type", "application/json")
	return
}
