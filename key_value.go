package forest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// ConfigResponse struct
type ConfigResponse struct {
	RequestID     string          `json:"request_id"`
	LeaseID       string          `json:"lease_id"`
	Renewable     bool            `json:"renewable"`
	LeaseDuration int             `json:"lease_duration"`
	WrapInfo      interface{}     `json:"wrap_info"`
	Warnings      []string        `json:"warnings"`
	Data          json.RawMessage `json:"data"`
}

type configList struct {
	RequestID     string      `json:"request_id"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int         `json:"lease_duration"`
	WrapInfo      interface{} `json:"wrap_info"`
	Warnings      []string    `json:"warnings"`
	Data          struct {
		Keys []string `json:"keys"`
	} `json:"data"`
}

//  DEPRECATED. USE GetKeyValue instead.
//
// GetConfig returns config from vault. The path format is '/{secret engine name}/{secret name}'
//
// Example: `data, err := vault.GetConfig(context.Background(), "/kv/foo")`
func (v *Vault) GetConfig(ctx context.Context, path string) (data []byte, err error) {
	req, err := v.requestGen(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err := checkErrorResponse(res); err != nil {
		return nil, err
	}
	var ddd ConfigResponse
	err = json.NewDecoder(res.Body).Decode(&ddd)
	if err != nil {
		return nil, err
	}
	return ddd.Data, nil
}

//  DEPRECATED. USE GetKeyValueLoad instead.
//
// GetConfigLoad returns config and loaded into a variable
// The path format is '/{secret engine name}/{secret name}'
func (v *Vault) GetConfigLoad(ctx context.Context, path string, model interface{}) (err error) {
	data, err := v.GetConfig(ctx, path)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, model)
	return
}

// GetEngineKeys returns list of keys for given engine
func (v *Vault) GetEngineKeys(ctx context.Context, engine string) (configs []string, err error) {
	path := fmt.Sprintf("/%s?list=true", engine)
	req, err := v.requestGen(ctx, http.MethodGet, path, nil)
	if err != nil {
		return
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err := checkErrorResponse(res); err != nil {
		return nil, err
	}
	var ddd configList
	err = json.NewDecoder(res.Body).Decode(&ddd)
	if err != nil {
		return nil, err
	}
	return ddd.Data.Keys, nil
}

// GetKeyValue returns key value store from vault. 'Key' is the key name. Like for example 'ms-order-conf'
func (v *Vault) GetKeyValue(ctx context.Context, key string) (data []byte, err error) {
	path := fmt.Sprintf("/%s/%s", v.Config.KeyValueEngine, key)
	req, err := v.requestGen(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err := checkErrorResponse(res); err != nil {
		return nil, err
	}
	var ddd ConfigResponse
	err = json.NewDecoder(res.Body).Decode(&ddd)
	if err != nil {
		return nil, err
	}
	return ddd.Data, nil
}

// GetKeyValueLoad returns config and loaded into a variable
// 'Key' is the key name. Like for example 'ms-order-conf'
func (v *Vault) GetKeyValueLoad(ctx context.Context, key string, model interface{}) (err error) {
	data, err := v.GetKeyValue(ctx, key)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, model)
	return
}

// UpsertKeyValue method
func (v *Vault) UpsertKeyValue(ctx context.Context, key string, data interface{}) (err error) {
	var val []byte
	switch t := data.(type) {
	case []byte:
		val = t
	case string:
		val = []byte(t)
	default:
		kind := reflect.Indirect(reflect.ValueOf(data)).Kind()
		if kind == reflect.Struct || kind == reflect.Map || kind == reflect.Slice {
			val, _ = json.Marshal(data)
		} else {
			return fmt.Errorf("unsupported file type: '%s'", kind)
		}
	}
	path := fmt.Sprintf("/%s/%s", v.Config.KeyValueEngine, key)
	req, err := v.requestGen(ctx, http.MethodPost, path, bytes.NewBuffer(val))
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
