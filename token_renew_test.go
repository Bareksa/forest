package begundal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_RenewTokenOther(t *testing.T) {
	token, err := CreateNewToken().WithPolicies("default").WithTimeToLive(1 * time.Hour).Do(context.TODO())
	require.NoError(t, err)
	err = RenewTokenOther(context.TODO(), token)
	require.NoError(t, err)
}
