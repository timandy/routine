package routine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetgLink(t *testing.T) {
	gp := getg()
	assert.NotNil(t, gp)
}
