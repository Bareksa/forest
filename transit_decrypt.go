package forest

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type decryptResponse struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int64  `json:"lease_duration"`
	Data          struct {
		Plaintext string `json:"plaintext"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

type decryptRequest struct {
	Ciphertext string        `json:"ciphertext"`
	Context    string        `json:"context,omitempty"`
	Nonce      string        `json:"nonce,omitempty"`
	BatchInput []interface{} `json:"batch_input,omitempty"`
}

// TransitDecrypt decrypts a transit encrypted payload.
// If decyprting a big ciphertext like if decrypted it's actually an image, please use TransitDecryptStream.
func (v *Vault) TransitDecrypt(ctx context.Context, key, cipherText string) (data []byte, err error) {
	path := fmt.Sprintf("/%s/decrypt/%s", v.Config.TransitEngine, key)
	trimmedCipher := strings.TrimSpace(cipherText)
	body := decryptRequest{
		Ciphertext: trimmedCipher,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := v.requestGen(ctx, http.MethodPost, path, bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if err := checkErrorResponse(res); err != nil {
		return nil, err
	}
	var response decryptResponse
	json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	s, err := base64.StdEncoding.DecodeString(response.Data.Plaintext)
	if err != nil {
		return data, errors.New("Fail to decode payload")
	}
	trimmed := strings.TrimSpace(string(s))
	return []byte(trimmed), nil
}

// TransitDecryptStream decrypts a transit encrypted payload in streaming manner.
// Best usage is if you expect a big ciphertext from whatever your source is.
func (v *Vault) TransitDecryptStream(ctx context.Context, key string, cipher io.Reader) (payload io.Reader, err error) {
	path := fmt.Sprintf("/%s/decrypt/%s", v.Config.TransitEngine, key)
	// Initiate buf json holder
	buf := bytes.NewBufferString(`{"ciphertext":"`)
	// Starts the payload encryption
	_, err = io.Copy(buf, cipher)
	if err != nil {
		return nil, err
	}

	_, err = buf.WriteString(`"}`)
	if err != nil {
		return nil, err
	}
	req, err := v.requestGen(ctx, http.MethodPost, path, buf)
	if err != nil {
		return nil, err
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	p := newStreamInBetween(`"plaintext":"`, '"')
	_, err = io.Copy(p, res.Body)
	if err != nil && err != io.ErrShortWrite {
		return
	}
	payload = base64.NewDecoder(base64.StdEncoding, p)
	return payload, nil
}
