package forest

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Stream(t *testing.T) {
	want := `{"ayam":"kuning","harga":29000}`
	var cipher string
	t.Run("Testing Encrypt Stream", func(t *testing.T) {
		ccc := bytes.NewBufferString(want)
		data, err := TransitEncryptStream(context.TODO(), "aes", ccc)
		require.NoError(t, err)
		sss, err := ioutil.ReadAll(data)
		require.NoError(t, err)
		assert.Contains(t, string(sss), "vault:v1")

		cipher = string(sss)
	})

	t.Run("Testing Decrypt Stream", func(t *testing.T) {
		ccc := strings.NewReader(cipher)
		data, err := TransitDecryptStream(context.TODO(), "aes", ccc)
		require.NoError(t, err)
		got, err := ioutil.ReadAll(data)
		require.NoError(t, err)

		assert.Equal(t, want, string(got))
	})
}
