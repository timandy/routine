package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneable(t *testing.T) {
	//struct can not be cast to interface
	var value interface{} = personCloneable{Id: 1, Name: "Hello"}
	_, ok := value.(Cloneable)
	assert.False(t, ok)
	//pointer can be cast to interface
	var pointer interface{} = &personCloneable{Id: 1, Name: "Hello"}
	_, ok2 := pointer.(Cloneable)
	assert.True(t, ok2)
}

func TestCloneable_Clone(t *testing.T) {
	//clone struct
	pc := &personCloneable{Id: 1, Name: "Hello"}
	assert.NotSame(t, pc, pc.Clone())
	assert.Equal(t, *pc, *(pc.Clone().(*personCloneable)))
	//copy pointer
	pcs := make([]*personCloneable, 1)
	pcs[0] = pc
	pcs2 := make([]*personCloneable, 1)
	copy(pcs2, pcs)
	assert.Same(t, pc, pcs2[0])
}

type personCloneable struct {
	Id   int
	Name string
}

func (p *personCloneable) Clone() any {
	return &personCloneable{Id: p.Id, Name: p.Name}
}
