package begundal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type cc struct {
	Capabilities []Capability `json:"capabilities"`
}

// Capability s
type Capability string

const (
	// PathTokenLookupSelf Used for Policy Creation
	PathTokenLookupSelf = "auth/token/lookup-self"
	// PathTokenRevokeSelf Used for Policy Creation
	PathTokenRevokeSelf = "auth/token/revoke-self"
	// PathSysCapabilitiesSelf Used for Policy Creation
	PathSysCapabilitiesSelf = "sys/capabilities-self"
	// PathTokenRevoke Used for Policy Creation
	PathTokenRevoke = "auth/token/revoke"
	// PathTokenRoot Used for checking if token is root. *Untested*
	PathTokenRoot = "auth/token/root"
)

const (
	// CapabilityRead Determines if the token will be able to access the resource. Used for Policy Creation
	CapabilityRead Capability = "read"
	// CapabilityUpdate Determines if the token will be able to update the resource. Used for Policy Creation
	CapabilityUpdate Capability = "update"
	// CapabilityList Determines if the token will be able to LIST resources (not access them). Implicitly allows Cabality of Create.
	// Used for Policy Creation
	CapabilityList Capability = "list"
	// CapabilityDelete Determines if the token can delete a resource. Used for Policy Creation
	CapabilityDelete Capability = "delete"
	// CapabilityCreate Determines if the token used can create a resource.
	CapabilityCreate Capability = "create"
)

// PolicyResponse struct
type PolicyResponse struct {
	RequestID     string      `json:"request_id"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int         `json:"lease_duration"`
	WrapInfo      interface{} `json:"wrap_info"`
	Warnings      []string    `json:"warnings"`
	Data          PolicyData  `json:"data"`
}

// PolicyData struct
type PolicyData struct {
	Name   string `json:"name"`
	Policy string `json:"policy"`
}

// CheckPolicy checks if policy exists and gets it's data
func (v *Vault) CheckPolicy(ctx context.Context, policy string) (d PolicyData, err error) {
	req, err := v.requestGen(ctx, http.MethodGet, "/sys/policies/acl/"+policy, nil)
	if err != nil {
		return
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	err = checkErrorResponse(res)
	if err != nil {
		return
	}
	var resp PolicyResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return
	}
	return resp.Data, nil
}

// UpsertPolicy creates/updates a policy.
// Token used in the instance must have the permission to even update policy itself.
// Root token have all permissions
func (v *Vault) UpsertPolicy(ctx context.Context, policy string, permissions map[string][]Capability) (err error) {
	var paths = make(map[string]cc)
	for key, cap := range permissions {
		paths[key] = cc{
			Capabilities: cap,
		}
	}
	perms, err := json.MarshalIndent(map[string]interface{}{"path": paths}, "", "  ")
	if err != nil {
		return
	}
	body, err := json.Marshal(map[string]interface{}{"policy": string(perms)})
	if err != nil {
		return
	}
	req, err := v.requestGen(ctx, http.MethodPut, "/sys/policies/acl/"+policy, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	res, err := v.Config.HTTPClient.Do(req)
	if err != nil {
		return
	}
	err = checkErrorResponse(res)
	return
}
