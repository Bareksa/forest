package begundal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type errorResponse struct {
	Errors []string `json:"errors"`
}

func checkErrorResponse(res *http.Response) (err error) {
	if res.StatusCode >= 400 {
		var errRes errorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			return
		}
		errs := strings.Join(errRes.Errors, ", ")
		if errs == "" {
			errs = "invalid or wrong path. Most likely the requested resource does not exist."
		}
		return fmt.Errorf("Error: %s", errs)
	}
	return
}
