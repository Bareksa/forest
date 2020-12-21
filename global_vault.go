package begundal

import (
	"context"
	"errors"
	"fmt"
	"io"
)

var gVault *Vault

// Init enables global top level functions
func Init(token string, opts ...OptionFunc) (err error) {
	if token == "" {
		return errors.New("Empty token")
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

	gVault = &Vault{
		BaseURL: fmt.Sprintf("%s/%s", config.Host, config.VaultAPIVersion),
		Config:  *config,
	}

	return
}

func checkNil() (err error) {
	if gVault == nil {
		return errors.New("Vault needs to be initialized first. Please call Init(token, opts) at least once before using other methods")
	}
	return
}

// GetConfig returns config from vault. The path format is '/{secret engine name}/{secret name}'
//
// Example: `data, err := vault.GetConfig(context.Background(), "/kv/foo")`
func GetConfig(ctx context.Context, path string) (data []byte, err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.GetConfig(ctx, path)
}

// GetConfigLoad returns config and loaded into a variable
// The path format is '/{secret engine name}/{secret name}'
func GetConfigLoad(ctx context.Context, path string, model interface{}) (err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.GetConfigLoad(ctx, path, model)
}

// GetKeyValue returns key value store from vault. 'Key' is the key name. Like for example 'ms-order-conf'
func GetKeyValue(ctx context.Context, key string) (data []byte, err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.GetKeyValue(ctx, key)
}

// GetKeyValueLoad returns config and loaded into a variable
// 'Key' is the key name. Like for example 'ms-order-conf'
func GetKeyValueLoad(ctx context.Context, key string, model interface{}) (err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.GetKeyValueLoad(ctx, key, model)
}

// TransitEncrypt will encrypt payload for sending somewhere else. 'key' is the encryptor name.
func TransitEncrypt(ctx context.Context, key string, payload []byte) (cipherText string, err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.TransitEncrypt(ctx, key, payload)
}

// TransitDecrypt decrypts a transit encrypted payload
func TransitDecrypt(ctx context.Context, key, cipherText string) (data []byte, err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.TransitDecrypt(ctx, key, cipherText)
}

// RenewTokenSelf attempts to renew token registered to self
func RenewTokenSelf(ctx context.Context) (err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.RenewTokenSelf(ctx)
}

// RenewTokenOther attempts to renew passed token using self's token
func RenewTokenOther(ctx context.Context, token string) (err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.RenewTokenOther(ctx, token)
}

// RenewTokenOverride attempts to renew passed token with passed token as auth
func RenewTokenOverride(ctx context.Context, token string) (err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.RenewTokenOverride(ctx, token)
}

// CreateNewToken creates an instance of Token override. Call ``.Do(ctx)`` on the instance to actually create new token.
// Default value for instance follows vault token documentation,
// on here: https://www.vaultproject.io/api/auth/token#parameters
//
// Which means the new token will be renewable by default and has display name of 'token'
func CreateNewToken() *CreateTokenInstance {
	return gVault.CreateNewToken()
}

// LookupSelf lookup information on current token used in the instance
func LookupSelf(ctx context.Context) (token LookupToken, err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.LookupSelf(ctx)
}

// LookupOther lookup information on current token used in the instance
func LookupOther(ctx context.Context, token string) (lookup LookupToken, err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.LookupOther(ctx, token)
}

// TransitEncryptStream will encrypt payload in stream manner to prevent memory overload on huge number of operation.
// Use this function for big files. The returned io Reader is a stream of pure encoded vault data.
func TransitEncryptStream(ctx context.Context, key string, payload io.Reader) (cipher io.Reader, err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.TransitEncryptStream(ctx, key, payload)
}

// TransitDecryptStream decrypts a transit encrypted payload in streaming manner.
// Best usage is if you expect a big ciphertext from whatever your source is.
func TransitDecryptStream(ctx context.Context, key string, cipher io.Reader) (payload io.Reader, err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.TransitDecryptStream(ctx, key, cipher)
}

// CheckPolicy checks if policy exists and gets it's data
func CheckPolicy(ctx context.Context, policy string) (d PolicyData, err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.CheckPolicy(ctx, policy)
}

// UpsertPolicy creates/updates a policy.
// Token used in the instance must have the permission to even update policy itself.
// Root token have all permissions
func UpsertPolicy(ctx context.Context, policy string, permissions map[string][]Capability) (err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.UpsertPolicy(ctx, policy, permissions)
}

// RevokeToken revokes given token string.
// Cannot revoke root token if instance token is not root
// If getting permission denied error, then it's most likely that reason.
func RevokeToken(ctx context.Context, token string) (err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.RevokeToken(ctx, token)
}

// GetConfigInstance returns the config used by this instance
func GetConfigInstance() Config {
	return gVault.Config
}

// GetKVEngine returns the engine name used for the instance
func GetKVEngine() string {
	return gVault.Config.KeyValueEngine
}

// GetTransitEngine returns the engine name used for the instance
func GetTransitEngine() string {
	return gVault.Config.TransitEngine
}

// GetEngineKeys returns list of keys for given engine
func GetEngineKeys(ctx context.Context, engine string) (configs []string, err error) {
	if err = checkNil(); err != nil {
		return
	}
	return gVault.GetEngineKeys(ctx, engine)
}
