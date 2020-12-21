package forest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LookupSelf(t *testing.T) {
	look, err := LookupSelf(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, *testToken, look.Data.ID)
}

func Test_LookupOther(t *testing.T) {
	token, err := CreateNewToken().WithPolicies("default").WithTimeToLive(1 * time.Hour).Do(context.TODO())
	require.NoError(t, err)
	look, err := LookupOther(context.TODO(), token)
	require.NoError(t, err)
	assert.Equal(t, token, look.Data.ID)
	assert.Less(t, look.Data.TTL, int64(3601))
}
