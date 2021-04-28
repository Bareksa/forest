package forest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TransitDecrypt_TransitDecrypt(t *testing.T) {
	want := `{"ayam":"kuning"}`
	got, err := TransitDecrypt(context.TODO(), "aes", `vault:v1:6Xyg+Yk8VDMLhKQTU7J+1FGAaFyk9qTwyB9SgEySEhv4pNKnZ2jEkbBLCAlW`)
	require.NoError(t, err)

	assert.Equal(t, want, string(got))
}
