package begundal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TransitDecrypt_TransitDecrypt(t *testing.T) {
	want := `{"ayam":"kuning"}`
	got, err := TransitDecrypt(context.TODO(), "awawa", `vault:v1:VVR7MhvibRHWmmbNB+zpaXFT7vNURkOC1E7E1l6fhUvIoDIGruvUCvOlK6+3SmYOSWjmwFGCA66s5CLnkTFq2D0O2EPbqGgfv5YVOzq2fuon/MI2tskZQulnSfjKeG15O2YHAOj8qimiywzvNKWx7IlS8oCn53dDIelRr2iD0qRQA/zBKQMd+i/3BxXvWiXydvab1TYmpbFrityZAeoe9qBGEnMeGEW68Gz8JbpcD9RA4ACP4eRzmp2I9eE3I0LZHsCrDgG+Y8jbgMMcRLrXaAIsSXR5NfjOoMJqHYCtiVaK3kiDuWrpMjTM9nVNSvp5DgU6jxmSoamviLUvuwFqxQ==`)
	require.NoError(t, err)

	assert.Equal(t, want, string(got))
}
