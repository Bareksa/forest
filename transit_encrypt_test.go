package forest

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TransitEncrypt_TransitEncrypt(t *testing.T) {
	data, err := TransitEncrypt(context.TODO(), "awawa", []byte(`{"ayam":"kuning","harga":29000}`))
	require.NoError(t, err)

	assert.Contains(t, string(data), "vault:v1")
}

func Test_TransitEncrypt_TransitEncryptStream(t *testing.T) {
	ccc := bytes.NewBufferString(`{"ayam":"kuning","harga":29000}`)
	data, err := TransitEncryptStream(context.TODO(), "awawa", ccc)
	require.NoError(t, err)
	sss, err := ioutil.ReadAll(data)
	require.NoError(t, err)
	assert.Contains(t, string(sss), "vault:v1")
}
