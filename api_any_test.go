package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	var value1 any = 1
	acceptAny(value1)
	acceptInterface(value1)
	//
	var value2 interface{} = 2
	acceptAny(value2)
	acceptInterface(value2)
	//
	value1 = value2
	assert.Equal(t, 2, value1)
}

//goland:noinspection GoUnusedParameter
func acceptAny(value any) {
}

//goland:noinspection GoUnusedParameter
func acceptInterface(value interface{}) {
}
