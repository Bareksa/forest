package begundal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateToken(t *testing.T) {
	var sToken string
	t.Run("Create Service Token", func(t *testing.T) {
		token, err := CreateNewToken().
			WithPolicies("default").
			WithDisplayName("testing").
			WithTimeToLive(1 * time.Hour).
			WithRoleName("kunyit").
			WithSetAsServiceToken().
			WithRenewableStatus(false).
			WithCurrentTokenAsParent(true).
			Do(context.Background())
		require.NoError(t, err)
		require.Contains(t, token, "s.")
		require.Len(t, token, 26)

		look, err := LookupOther(context.TODO(), token)
		require.NoError(t, err)
		assert.Equal(t, token, look.Data.ID)
		assert.Less(t, look.Data.TTL, int64(3601))
		assert.NotEqual(t, look.Data.TTL, int64(0))
		assert.NotContains(t, look.Data.Policies, "root")
		assert.Equal(t, false, look.Renewable)
		assert.Equal(t, "service", look.Data.Type)

		sToken = token
	})
	t.Run("Create Batch Token", func(t *testing.T) {
		token, err := CreateNewToken().
			WithPolicies("default").
			WithDisplayName("testing").
			WithTimeToLive(1 * time.Hour).
			WithRoleName("kunyit").
			WithSetAsBatchToken().
			WithCurrentTokenAsParent(true).
			Do(context.Background())
		require.NoError(t, err)
		assert.Contains(t, token, "b.")
		assert.Len(t, token, 172)

		look, err := LookupOther(context.TODO(), token)
		require.NoError(t, err)
		assert.Equal(t, token, look.Data.ID)
		assert.Less(t, look.Data.TTL, int64(3601))
		assert.NotEqual(t, look.Data.TTL, int64(0))
		assert.NotContains(t, look.Data.Policies, "root")
		assert.Equal(t, false, look.Renewable)
		assert.Equal(t, "batch", look.Data.Type)

	})
	t.Run("Revoke tokens", func(t *testing.T) {
		err := RevokeToken(context.TODO(), sToken)
		require.NoError(t, err)
		_, err = LookupOther(context.TODO(), sToken)
		require.NotNil(t, err)
	})
}
