package begundal

import "net/http"

// APIVersion the base type is string
type APIVersion string

// Config struct
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

const (
	// V2 Version is not supported until Hashicorp says otherwise
	V2 APIVersion = "v2"
	// V1 Default version used
	V1 APIVersion = "v1"
)

// OptionFunc opts
type OptionFunc func(options *Config)

// Options struct
type Options struct {
	Host            string       // Optional. Defaults to https://127.0.0.1:8200. Make sure to not add '/' in last character
	VaultAPIVersion APIVersion   // Optional. Defaults to v1
	HTTPClient      *http.Client // Uses default http client if nil
}

func WithHost(hostname string) OptionFunc {
	return func(o *Config) {
		o.Host = hostname
	}
}

func WithAPIVersion(version APIVersion) OptionFunc {
	return func(o *Config) {
		o.VaultAPIVersion = version
	}
}

func WithHttpClient(httpClient *http.Client) OptionFunc {
	return func(o *Config) {
		o.HTTPClient = httpClient
	}
}

func WithTransitEngine(engine string) OptionFunc {
	return func(o *Config) {
		o.TransitEngine = engine
	}
}

func WithKeyValueEngine(engine string) OptionFunc {
	return func(o *Config) {
		o.KeyValueEngine = engine
	}
}
