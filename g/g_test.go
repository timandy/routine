// Copyright 2018 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.

package g

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestG(t *testing.T) {
	gp1 := G()
	assert.NotNil(t, gp1, "fail to get G.")

	run := false
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		gp2 := G()
		assert.NotNil(t, gp2, "fail to get G.")
		assert.NotEqual(t, gp1, gp2, "every living G must be different. [gp1:%p] [gp2:%p]", gp1, gp2)
		run = true
		wg.Done()
	}()
	wg.Wait()
	assert.True(t, run)
}
