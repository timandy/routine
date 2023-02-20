// Copyright 2022 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

package g

import (
	"reflect"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetgp(t *testing.T) {
	gp0 := getgp()
	runtime.GC()
	assert.NotNil(t, gp0)
	//
	runTest(t, func() {
		gp := getgp()
		runtime.GC()
		assert.NotNil(t, gp)
		assert.NotEqual(t, gp0, gp)
	})
}

func TestGetg0(t *testing.T) {
	runTest(t, func() {
		g := getg0()
		runtime.GC()
		stackguard0 := reflect.ValueOf(g).FieldByName("stackguard0")
		assert.Greater(t, stackguard0.Uint(), uint64(0))
	})
}

func TestGetgt(t *testing.T) {
	runTest(t, func() {
		gt := getgt()
		runtime.GC()
		assert.Equal(t, "g", gt.Name())
		//
		assert.Greater(t, gt.NumField(), 20)
	})
}

func runTest(t *testing.T, fun func()) {
	run := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		fun()
		run = true
		wg.Done()
	}()
	wg.Wait()
	assert.True(t, run)
}
