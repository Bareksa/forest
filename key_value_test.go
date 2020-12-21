package forest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_KeyValue_GetConfig(t *testing.T) {
	data, err := GetConfig(context.TODO(), "/kv/ayam")
	require.NoError(t, err)

	assert.Equal(t, `{"bakar":["mentega","solo"],"goreng":["geprek","kfc"]}`, string(data))
}

func Test_KeyValue_GetConfigLoad(t *testing.T) {
	want := map[string]interface{}{
		"bakar":  []interface{}{"mentega", "solo"},
		"goreng": []interface{}{"geprek", "kfc"},
	}
	var got map[string]interface{}
	err := GetConfigLoad(context.TODO(), "/kv/ayam", &got)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func Test_KeyValue_GetKeyValue(t *testing.T) {
	data, err := GetKeyValue(context.TODO(), "ayam")
	require.NoError(t, err)

	want := `{"bakar":["mentega","solo"],"goreng":["geprek","kfc"]}`

	assert.Equal(t, want, string(data))
}

func Test_KeyValue_GetKeyValueLoad(t *testing.T) {
	want := map[string]interface{}{
		"bakar":  []interface{}{"mentega", "solo"},
		"goreng": []interface{}{"geprek", "kfc"},
	}

	var got map[string]interface{}

	err := GetKeyValueLoad(context.TODO(), "ayam", &got)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func Test_KeyValue_GetEngineKeyList(t *testing.T) {
	config, err := GetEngineKeys(context.TODO(), GetKVEngine())
	require.NoError(t, err)

	assert.NotEmpty(t, config)
}
