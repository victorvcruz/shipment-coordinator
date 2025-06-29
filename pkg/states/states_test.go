package states

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState_MarshalJSON(t *testing.T) {
	state := SP
	data, err := state.MarshalJSON()
	assert.NoError(t, err, "expected no error on MarshalJSON")
	assert.Equal(t, `"SP"`, string(data), "expected 'SP' as marshaled data")
}

func TestState_UnmarshalJSON_Valid(t *testing.T) {
	var s State
	err := s.UnmarshalJSON([]byte(`"SP"`))
	assert.NoError(t, err, "expected no error for valid state code")
}

func TestState_UnmarshalJSON_Invalid(t *testing.T) {
	var s State
	err := s.UnmarshalJSON([]byte("XX"))
	assert.Error(t, err, "expected error for invalid state code")
}

func TestState_JSONRoundTrip(t *testing.T) {
	original := RJ
	data, err := json.Marshal(&original)
	assert.NoError(t, err, "expected no error on marshal")
	var decoded State
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err, "expected no error on unmarshal")
	assert.Equal(t, original.Sigla, decoded.Sigla, "expected same state code after round trip")
}

func TestState_JSONRoundTrip2(t *testing.T) {
	type testStruct struct {
		State State `json:"state"`
	}

	original := testStruct{
		State: SP,
	}

	data, err := json.Marshal(&original)
	assert.NoError(t, err, "expected no error on marshal")
	var decoded map[string]any
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err, "expected no error on unmarshal")
	assert.Equal(
		t,
		original.State.Sigla,
		decoded["state"].(string),
		"expected same state code after round trip",
	)
}

func TestRegion_MarshalJSON(t *testing.T) {
	region := CentroOeste
	data, err := region.MarshalJSON()
	assert.NoError(t, err, "expected no error on MarshalJSON")
	assert.Equal(t, `"Centro-Oeste"`, string(data), "expected 'Centro-Oeste' as marshaled data")
}

func TestRegion_UnmarshalJSON_Valid(t *testing.T) {
	var s Region
	err := s.UnmarshalJSON([]byte(`"Centro-Oeste"`))
	assert.NoError(t, err, "expected no error for valid state code")
}

func TestRegion_UnmarshalJSON_Invalid(t *testing.T) {
	var s Region
	err := s.UnmarshalJSON([]byte("XX"))
	assert.Error(t, err, "expected error for invalid state code")
}

func TestRegion_JSONRoundTrip(t *testing.T) {
	original := Sudeste
	data, err := json.Marshal(&original)
	assert.NoError(t, err, "expected no error on marshal")
	var decoded Region
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err, "expected no error on unmarshal")
	assert.Equal(t, original.Name, decoded.Name, "expected same region code after round trip")
}

func TestRegion_JSONRoundTrip2(t *testing.T) {
	type testStruct struct {
		Region Region `json:"region"`
	}

	original := testStruct{
		Region: Sudeste,
	}

	data, err := json.Marshal(&original)
	assert.NoError(t, err, "expected no error on marshal")
	var decoded map[string]any
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err, "expected no error on unmarshal")
	assert.Equal(
		t,
		original.Region.Name,
		decoded["region"].(string),
		"expected same state code after round trip",
	)
}
