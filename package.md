# begundal
--
    import "."


## Usage

```go
const (
	// PathTokenLookupSelf Used for Policy Creation
	PathTokenLookupSelf = "auth/token/lookup-self"
	// PathTokenRevokeSelf Used for Policy Creation
	PathTokenRevokeSelf = "auth/token/revoke-self"
	// PathSysCapabilitiesSelf Used for Policy Creation
	PathSysCapabilitiesSelf = "sys/capabilities-self"
)
```

#### func  GetConfig

```go
func GetConfig(ctx context.Context, path string) (data []byte, err error)
```
GetConfig returns config from vault. The path format is '/{secret engine
name}/{secret name}'

Example: `data, err := vault.GetConfig(context.Background(), "/kv/foo")`

#### func  GetConfigLoad

```go
func GetConfigLoad(ctx context.Context, path string, model interface{}) (err error)
```
GetConfigLoad returns config and loaded into a variable The path format is
'/{secret engine name}/{secret name}'

#### func  GetKeyValue

```go
func GetKeyValue(ctx context.Context, key string) (data []byte, err error)
```
GetKeyValue returns key value store from vault. 'Key' is the key name. Like for
example 'ms-order-conf'

#### func  GetKeyValueLoad

```go
func GetKeyValueLoad(ctx context.Context, key string, model interface{}) (err error)
```
GetKeyValueLoad returns config and loaded into a variable 'Key' is the key name.
Like for example 'ms-order-conf'

#### func  Init

```go
func Init(token string, opts ...OptionFunc) (err error)
```
Init enables global top level functions

#### func  RenewTokenOther

```go
func RenewTokenOther(ctx context.Context, token string) (err error)
```
RenewTokenOther attempts to renew passed token using self's token

#### func  RenewTokenOverride

```go
func RenewTokenOverride(ctx context.Context, token string) (err error)
```
RenewTokenOverride attempts to renew passed token with passed token as auth

#### func  RenewTokenSelf

```go
func RenewTokenSelf(ctx context.Context) (err error)
```
RenewTokenSelf attempts to renew token registered to self

#### func  TransitDecrypt

```go
func TransitDecrypt(ctx context.Context, key, cipherText string) (data []byte, err error)
```
TransitDecrypt decrypts a transit encrypted payload

#### func  TransitDecryptStream

```go
func TransitDecryptStream(ctx context.Context, key string, cipher io.Reader) (payload io.Reader, err error)
```
TransitDecryptStream decrypts a transit encrypted payload in streaming manner.
Best usage is if you expect a big ciphertext from whatever your source is.

#### func  TransitEncrypt

```go
func TransitEncrypt(ctx context.Context, key string, payload []byte) (cipherText string, err error)
```
TransitEncrypt will encrypt payload for sending somewhere else. 'key' is the
encryptor name.

#### func  TransitEncryptStream

```go
func TransitEncryptStream(ctx context.Context, key string, payload io.Reader) (cipher io.Reader, err error)
```
TransitEncryptStream will encrypt payload in stream manner to prevent memory
overload on huge number of operation. Use this function for big files. The
returned io Reader is a stream of pure encoded vault data.

#### func  UpsertPolicy

```go
func UpsertPolicy(ctx context.Context, policy string, permissions map[string][]Capability) (err error)
```
UpsertPolicy creates/updates a policy. Token used in the instance must have the
permission to even update policy itself. Root token have all permissions

#### type APIVersion

```go
type APIVersion string
```

APIVersion the base type is string

```go
const (
	// V2 Version is not supported until Hashicorp says otherwise
	V2 APIVersion = "v2"
	// V1 Default version used
	V1 APIVersion = "v1"
)
```

#### type Capability

```go
type Capability string
```

Capability s

```go
const (
	// CapabilityRead Determines if the token will be able to access the resource. Used for Policy Creation
	CapabilityRead Capability = "read"
	// CapabilityUpdate Determines if the token will be able to update the resource. Used for Policy Creation
	CapabilityUpdate Capability = "update"
	// CapabilityList Determines if the token will be able to LIST resources (not access them). Used for Policy Creation
	CapabilityList Capability = "list"
	// CapabilityDelete Determines if the token can delete a resource. Used for Policy Creation
	CapabilityDelete Capability = "delete"
)
```

#### type Config

```go
type Config struct {
	Token               string       // Required
	Host                string       // Optional. Defaults to https://127.0.0.1:8200. Make sure to not add '/' in last character
	NoRenew             bool         // Optional. Defaults to false, so it will attempt to renew token every time Renew Timing passes.
	NoRenewOnInitialize bool         // Optional. Defaults to false, which will indeed renew on initiazlie. Will not renew if NoRenew is set to true
	RenewTiming         string       // Optional. Defaults to 0 0 * * * (Every midnight). Uses cron tab syntax.
	VaultAPIVersion     APIVersion   // Optional. Defaults to v1
	HTTPClient          *http.Client // Uses default http client if nil
	KeyValueEngine      string
	TransitEngine       string
}
```

Config struct

#### type ConfigResponse

```go
type ConfigResponse struct {
	RequestID     string          `json:"request_id"`
	LeaseID       string          `json:"lease_id"`
	Renewable     bool            `json:"renewable"`
	LeaseDuration int             `json:"lease_duration"`
	WrapInfo      interface{}     `json:"wrap_info"`
	Warnings      []string        `json:"warnings"`
	Data          json.RawMessage `json:"data"`
}
```

ConfigResponse struct

#### type CreateTokenInstance

```go
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
}
```

CreateTokenInstance struct to create token

#### func  CreateNewToken

```go
func CreateNewToken() *CreateTokenInstance
```
CreateNewToken creates an instance of Token override. Call ``.Do(ctx)`` on the
instance to actually create new token. Default value for instance follows vault
token documentation, on here:
https://www.vaultproject.io/api/auth/token#parameters

Which means the new token will be renewable by default and has display name of
'token'

#### func (*CreateTokenInstance) Do

```go
func (c *CreateTokenInstance) Do(ctx context.Context) (token string, err error)
```
Do creates the token.

#### func (*CreateTokenInstance) DoOverride

```go
func (c *CreateTokenInstance) DoOverride(ctx context.Context, authToken string) (token string, err error)
```
DoOverride creates the token, with passed token as auth and parent if orphan
status is false (or no parent is true)

#### func (*CreateTokenInstance) WithCurrentTokenAsParent

```go
func (c *CreateTokenInstance) WithCurrentTokenAsParent(b bool) *CreateTokenInstance
```
WithCurrentTokenAsParent Uses instance token as parent. Default true. If using
DoOverride(token), the override token will be set as parent instead.

#### func (*CreateTokenInstance) WithDefaultPolicy

```go
func (c *CreateTokenInstance) WithDefaultPolicy(b bool) *CreateTokenInstance
```
WithDefaultPolicy replaces instance default policy. by default true

#### func (*CreateTokenInstance) WithDisplayName

```go
func (c *CreateTokenInstance) WithDisplayName(s string) *CreateTokenInstance
```
WithDisplayName replaces token display name

#### func (*CreateTokenInstance) WithEntityAlias

```go
func (c *CreateTokenInstance) WithEntityAlias(s string) *CreateTokenInstance
```
WithEntityAlias replaces instance entity alias. MUST BE USED alongside
WithRoleName and Role name must exist within vault.

#### func (*CreateTokenInstance) WithExplicitMaxTTL

```go
func (c *CreateTokenInstance) WithExplicitMaxTTL(t time.Duration) *CreateTokenInstance
```
WithExplicitMaxTTL replaces instance explicit time to live, which by default
will depends on Vault's default lease TTL If this method is called, duration
will be rounded down to the nearest hour argument passed with minimum value of 1
hour

#### func (*CreateTokenInstance) WithID

```go
func (c *CreateTokenInstance) WithID(id string) *CreateTokenInstance
```
WithID replaces instance ID. Make sure there's no character '.' in the argument
string

#### func (*CreateTokenInstance) WithMeta

```go
func (c *CreateTokenInstance) WithMeta(metadata map[string]string) *CreateTokenInstance
```
WithMeta replaces instance metadata

#### func (*CreateTokenInstance) WithNumberOfUses

```go
func (c *CreateTokenInstance) WithNumberOfUses(i int) *CreateTokenInstance
```
WithNumberOfUses replaces token's allowed number of usage. Signing in to vault
using UI with the token is considered used one time. By default 0, which is
infinite.

#### func (*CreateTokenInstance) WithPeriod

```go
func (c *CreateTokenInstance) WithPeriod(t time.Duration) *CreateTokenInstance
```
WithPeriod replaces token period. Token that is not renewed in this set of time
cannot be renewed again. By default, if unset will follow's Vault's default
lease TTL. Has minimum value of 1 hour. Only hourly is supported in this
package.

#### func (*CreateTokenInstance) WithPolicies

```go
func (c *CreateTokenInstance) WithPolicies(policies ...string) *CreateTokenInstance
```
WithPolicies replaces instance policies

#### func (*CreateTokenInstance) WithRenewableStatus

```go
func (c *CreateTokenInstance) WithRenewableStatus(b bool) *CreateTokenInstance
```
WithRenewableStatus replaces instance renewable. default true.

#### func (*CreateTokenInstance) WithRoleName

```go
func (c *CreateTokenInstance) WithRoleName(rolename string) *CreateTokenInstance
```
WithRoleName replaces instance rolename.

#### func (*CreateTokenInstance) WithSetAsBatchToken

```go
func (c *CreateTokenInstance) WithSetAsBatchToken() *CreateTokenInstance
```
WithSetAsBatchToken replaces instance token type

#### func (*CreateTokenInstance) WithSetAsServiceToken

```go
func (c *CreateTokenInstance) WithSetAsServiceToken() *CreateTokenInstance
```
WithSetAsServiceToken replaces instance token type

#### func (*CreateTokenInstance) WithTimeToLive

```go
func (c *CreateTokenInstance) WithTimeToLive(t time.Duration) *CreateTokenInstance
```
WithTimeToLive replaces instance time to live, which by default will depends on
Vault's default lease TTL If this method is called, duration will be rounded
down to the nearest hour argument passed with minimum value of 1 hour Only
hourly is supported in this package.

#### type LookupToken

```go
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
```

LookupToken token lookup

#### func  LookupOther

```go
func LookupOther(ctx context.Context, token string) (lookup LookupToken, err error)
```
LookupOther lookup information on current token used in the instance

#### func  LookupSelf

```go
func LookupSelf(ctx context.Context) (token LookupToken, err error)
```
LookupSelf lookup information on current token used in the instance

#### type OptionFunc

```go
type OptionFunc func(options *Config)
```

OptionFunc opts

#### func  WithAPIVersion

```go
func WithAPIVersion(version APIVersion) OptionFunc
```

#### func  WithHost

```go
func WithHost(hostname string) OptionFunc
```

#### func  WithHttpClient

```go
func WithHttpClient(httpClient *http.Client) OptionFunc
```

#### func  WithKeyValueEngine

```go
func WithKeyValueEngine(engine string) OptionFunc
```

#### func  WithTransitEngine

```go
func WithTransitEngine(engine string) OptionFunc
```

#### type Options

```go
type Options struct {
	Host            string       // Optional. Defaults to https://127.0.0.1:8200. Make sure to not add '/' in last character
	VaultAPIVersion APIVersion   // Optional. Defaults to v1
	HTTPClient      *http.Client // Uses default http client if nil
}
```

Options struct

#### type PolicyData

```go
type PolicyData struct {
	Name   string `json:"name"`
	Policy string `json:"policy"`
}
```

PolicyData struct

#### func  CheckPolicy

```go
func CheckPolicy(ctx context.Context, policy string) (d PolicyData, err error)
```
CheckPolicy checks if policy exists and gets it's data

#### type PolicyResponse

```go
type PolicyResponse struct {
	RequestID     string      `json:"request_id"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int         `json:"lease_duration"`
	WrapInfo      interface{} `json:"wrap_info"`
	Warnings      []string    `json:"warnings"`
	Data          PolicyData  `json:"data"`
}
```

PolicyResponse struct

#### type Vault

```go
type Vault struct {
	BaseURL string
	Config  Config
}
```

Vault struct

#### func  NewClient

```go
func NewClient(token string, opts ...OptionFunc) (*Vault, error)
```
NewClient creates a new vault instance

#### func (*Vault) CheckPolicy

```go
func (v *Vault) CheckPolicy(ctx context.Context, policy string) (d PolicyData, err error)
```
CheckPolicy checks if policy exists and gets it's data

#### func (*Vault) CreateNewToken

```go
func (v *Vault) CreateNewToken() *CreateTokenInstance
```
CreateNewToken creates an instance of Token override. Call ``.Do(ctx)`` on the
instance to actually create new token. Default value for instance follows vault
token documentation, on here:
https://www.vaultproject.io/api/auth/token#parameters

Which means the new token will be renewable by default and has display name of
'token'

#### func (*Vault) GetConfig

```go
func (v *Vault) GetConfig(ctx context.Context, path string) (data []byte, err error)
```
GetConfig returns config from vault. The path format is '/{secret engine
name}/{secret name}'

Example: `data, err := vault.GetConfig(context.Background(), "/kv/foo")`

#### func (*Vault) GetConfigLoad

```go
func (v *Vault) GetConfigLoad(ctx context.Context, path string, model interface{}) (err error)
```
GetConfigLoad returns config and loaded into a variable The path format is
'/{secret engine name}/{secret name}'

#### func (*Vault) GetKeyValue

```go
func (v *Vault) GetKeyValue(ctx context.Context, key string) (data []byte, err error)
```
GetKeyValue returns key value store from vault. 'Key' is the key name. Like for
example 'ms-order-conf'

#### func (*Vault) GetKeyValueLoad

```go
func (v *Vault) GetKeyValueLoad(ctx context.Context, key string, model interface{}) (err error)
```
GetKeyValueLoad returns config and loaded into a variable 'Key' is the key name.
Like for example 'ms-order-conf'

#### func (*Vault) LookupOther

```go
func (v *Vault) LookupOther(ctx context.Context, token string) (lookup LookupToken, err error)
```
LookupOther lookup information on passed token

#### func (*Vault) LookupSelf

```go
func (v *Vault) LookupSelf(ctx context.Context) (lookup LookupToken, err error)
```
LookupSelf lookup information on current token used in the instance

#### func (*Vault) RenewTokenOther

```go
func (v *Vault) RenewTokenOther(ctx context.Context, token string) (err error)
```
RenewTokenOther attempts to renew passed token using self's token

#### func (*Vault) RenewTokenOverride

```go
func (v *Vault) RenewTokenOverride(ctx context.Context, token string) (err error)
```
RenewTokenOverride attempts to renew passed token with passed token as auth

#### func (*Vault) RenewTokenSelf

```go
func (v *Vault) RenewTokenSelf(ctx context.Context) (err error)
```
RenewTokenSelf attempts to renew token registered to self. Cannot renew root
token with 0 time to live (never expire).

#### func (*Vault) TransitDecrypt

```go
func (v *Vault) TransitDecrypt(ctx context.Context, key, cipherText string) (data []byte, err error)
```
TransitDecrypt decrypts a transit encrypted payload. If decyprting a big
ciphertext like if decrypted it's actually an image, please use
TransitDecryptStream.

#### func (*Vault) TransitDecryptStream

```go
func (v *Vault) TransitDecryptStream(ctx context.Context, key string, cipher io.Reader) (payload io.Reader, err error)
```
TransitDecryptStream decrypts a transit encrypted payload in streaming manner.
Best usage is if you expect a big ciphertext from whatever your source is.

#### func (*Vault) TransitEncrypt

```go
func (v *Vault) TransitEncrypt(ctx context.Context, key string, payload []byte) (cipherText string, err error)
```
TransitEncrypt will encrypt payload for sending somewhere else. 'key' is the
encryptor name.

#### func (*Vault) TransitEncryptStream

```go
func (v *Vault) TransitEncryptStream(ctx context.Context, key string, payload io.Reader) (io.Reader, error)
```
TransitEncryptStream will encrypt payload in stream manner to prevent memory
overload on huge number of operation. Use this function for big files. The
returned io Reader is a stream of pure encoded vault data without bells and
whistles of JSON.

#### func (*Vault) UpsertPolicy

```go
func (v *Vault) UpsertPolicy(ctx context.Context, policy string, permissions map[string][]Capability) (err error)
```
UpsertPolicy creates/updates a policy. Token used in the instance must have the
permission to even update policy itself. Root token have all permissions
