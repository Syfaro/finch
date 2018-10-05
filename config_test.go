package finch

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigGet(t *testing.T) {
	conf := make(Config)

	assert.Nil(t, conf["testing"])
	conf["testing"] = true
	res := conf.Get("testing")
	assert.IsType(t, true, res)
	assert.True(t, conf.Get("testing").(bool))

	assert.Nil(t, conf.Get("test_var"))
	os.Setenv("TEST_VAR", "val")
	res = conf.Get("test_var")
	assert.IsType(t, "", res)
	assert.Equal(t, "val", res.(string))
}
