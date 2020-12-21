package begundal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// RevokeToken revokes given token string.
// Cannot revoke root token if instance token is not root
// If getting permission denied error, then it's most likely that reason.
func (v *Vault) RevokeToken(ctx context.Context, token string) (err error) {
	body := sLookupOther{
		Token: token,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return
	}
	req, err := v.requestGen(ctx, http.MethodPost, "/auth/token/revoke", bytes.NewBuffer(b))
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
