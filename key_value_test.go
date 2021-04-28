package forest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_KeyValue_GetConfig(t *testing.T) {
	data, err := GetConfig(context.TODO(), "/forest_kv/forest-data")
	require.NoError(t, err)

	assert.Equal(t, `{"bakar":["mentega","solo"],"goreng":["geprek","kfc"]}`, string(data))
}

func Test_KeyValue_GetConfigLoad(t *testing.T) {
	want := map[string]interface{}{
		"bakar":  []interface{}{"mentega", "solo"},
		"goreng": []interface{}{"geprek", "kfc"},
	}
	var got map[string]interface{}
	err := GetConfigLoad(context.TODO(), "/forest_kv/forest-data", &got)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func Test_KeyValue_GetKeyValue(t *testing.T) {
	data, err := GetKeyValue(context.TODO(), "forest-data")
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

	err := GetKeyValueLoad(context.TODO(), "forest-data", &got)
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func Test_KeyValue_GetEngineKeyList(t *testing.T) {
	config, err := GetEngineKeys(context.TODO(), GetKVEngine())
	require.NoError(t, err)

	assert.NotEmpty(t, config)
}

func Test_KeyValue_UpsertKeyValue(t *testing.T) {
	payload := `{"bakar":["mentega","solo"],"goreng":["geprek","kfc","tepung"]}`
	origin := `{"bakar":["mentega","solo"],"goreng":["geprek","kfc"]}`
	err := UpsertKeyValue(context.TODO(), "forest-data", payload)
	require.Nil(t, err)

	data, err := GetKeyValue(context.TODO(), "forest-data")
	require.NoError(t, err)

	require.Equal(t, payload, string(data))

	err = UpsertKeyValue(context.TODO(), "forest-data", origin)
	require.NoError(t, err)
}
