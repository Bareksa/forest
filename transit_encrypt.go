package begundal

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Plaintext is the required field
type encryptRequest struct {
	Plaintext            string        `json:"plaintext"`
	Context              string        `json:"context,omitempty"`
	KeyVersion           int           `json:"key_version,omitempty"`
	Nonce                string        `json:"nonce,omitempty"`
	BatchInput           []interface{} `json:"batch_input,omitempty"`
	Type                 string        `json:"type,omitempty"`
	ConvergentEncryption string        `json:"convergent_encryption,omitempty"`
}

// The one we want is Data.Ciphertext
type encryptResponse struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int64  `json:"lease_duration"`
	Data          struct {
		Ciphertext string `json:"ciphertext"`
		KeyVersion int    `json:"key_version"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings []string    `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

type encryptResponseStream struct {
	startData bool
	initiate  strings.Builder
	temp      *bytes.Buffer
}

func (e *encryptResponseStream) Read(p []byte) (n int, err error) {
	return e.temp.Read(p)
}

func (e *encryptResponseStream) Write(p []byte) (n int, err error) {
	if e.temp == nil {
		e.temp = bytes.NewBuffer(nil)
	}
	for _, b := range p {
		if e.startData {
			// Stop early since we hit the closing quote
			if b == '"' {
				return 0, nil
			}
			err = e.temp.WriteByte(b)
			if err != nil {
				return 0, err
			}
		} else {
			err = e.initiate.WriteByte(b)
			if err != nil {
				return 0, err
			}
			dd := e.initiate.String()
			if strings.Contains(dd, `"ciphertext":"`) {
				e.startData = true
			}
		}
		n++
	}
	return
}

// TransitEncrypt will encrypt payload for sending somewhere else. 'key' is the encryptor name.
func (v *Vault) TransitEncrypt(ctx context.Context, key string, payload []byte) (cipherText string, err error) {
	encoded := base64.StdEncoding.EncodeToString(payload)
	path := fmt.Sprintf("/%s/encrypt/%s", v.Config.TransitEngine, key)
	body := encryptRequest{
		Plaintext: encoded,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return "", err
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
		return "", err
	}
	response := encryptResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}
	cipherText = response.Data.Ciphertext
	return
}

// TransitEncryptStream will encrypt payload in stream manner to prevent memory overload on huge number of operation.
// Use this function for big files. The returned io Reader is a stream of pure encoded vault data without bells and whistles of JSON.
func (v *Vault) TransitEncryptStream(ctx context.Context, key string, payload io.Reader) (io.Reader, error) {
	// Initiate buf json holder
	buf := bytes.NewBufferString(`{"plaintext":"`)

	// Creates the encoder. Any data written to encoder will be written to buf in encoded fashion
	encoded := base64.NewEncoder(base64.StdEncoding, buf)

	// Starts the payload encryption
	_, err := io.Copy(encoded, payload)
	if err != nil {
		return nil, err
	}
	// Close has to be called now to trigger adding remaining data to buf
	err = encoded.Close()
	if err != nil {
		return nil, err
	}

	// Enclose the json stream
	_, err = buf.WriteString(`"}`)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/%s/encrypt/%s", v.Config.TransitEngine, key)
	req, err := v.requestGen(ctx, http.MethodPost, path, buf)
	if err != nil {
		return nil, err
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	cipher := newStreamInBetween(`"ciphertext":"`, '"')
	_, err = io.Copy(cipher, res.Body)
	if err != nil && err != io.ErrShortWrite {
		return nil, err
	}
	return cipher, nil
}
