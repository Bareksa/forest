package forest

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Policies(t *testing.T) {
	ss := map[string][]Capability{
		PathTokenLookupSelf:     {CapabilityRead},
		PathTokenRevokeSelf:     {CapabilityUpdate},
		PathSysCapabilitiesSelf: {CapabilityUpdate},
	}
	xx := make(map[string]cc)
	for k, v := range ss {
		xx[k] = cc{Capabilities: v}
	}
	t.Run("Create Policy", func(t *testing.T) {
		err := UpsertPolicy(context.TODO(), "test-policy", ss)
		require.NoError(t, err)
	})

	t.Run("Check Policy", func(t *testing.T) {
		data, err := CheckPolicy(context.TODO(), "test-policy")
		require.NoError(t, err)
		want, err := json.MarshalIndent(map[string]interface{}{"path": xx}, "", "  ")
		require.NoError(t, err)
		assert.Equal(t, string(want), data.Policy)
	})
}
