package version

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReleaseVersion_NotEmpty(t *testing.T) {
	v := ReleaseVersion()
	assert.NotEmpty(t, v)
}

func TestAPI_StringAndJSONConsistent(t *testing.T) {
	s := String()
	j := JSON()

	assert.NotEmpty(t, s)
	assert.NotEmpty(t, j)

	var infoFromJSON Info
	err := json.Unmarshal([]byte(j), &infoFromJSON)
	assert.NoError(t, err)

	assert.Equal(t, infoFromJSON.String(), s)
}
