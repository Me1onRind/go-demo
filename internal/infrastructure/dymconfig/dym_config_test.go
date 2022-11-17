package dymconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCfg struct {
	Str      string `json:"str" yaml:"str"`
	IntSlice []int  `json:"int_slice" yaml:"int_slice"`
}

func Test_UnmarshalBySuffix_Json(t *testing.T) {
	t1 := &testCfg{}
	err := unmarshalBySuffix("test.json", []byte(`{"str":"str_test","int_slice":[1,2,3]}`), t1)
	assert.Empty(t, err)
	assert.Equal(t, "str_test", t1.Str)
	assert.Equal(t, []int{1, 2, 3}, t1.IntSlice)
}

func Test_UnmarshalBySuffix_Yaml(t *testing.T) {
	t1 := &testCfg{}
	for _, key := range []string{"test.yaml", "test.yml"} {
		err := unmarshalBySuffix(key, []byte(
			"str: str_test\n"+
				"int_slice:\n"+
				" - 1\n - 2\n - 3\n",
		), t1)
		assert.Empty(t, err)
		assert.Equal(t, "str_test", t1.Str)
		assert.Equal(t, []int{1, 2, 3}, t1.IntSlice)
	}
}
