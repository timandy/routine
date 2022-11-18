package routine

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
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
		g0 := getg0()
		runtime.GC()
		stackguard0 := reflect.ValueOf(g0).FieldByName("stackguard0")
		assert.Greater(t, stackguard0.Uint(), uint64(0))
	})
}

func TestGetgt(t *testing.T) {
	fmt.Println("*** GOOS:", runtime.GOOS, "***")
	fmt.Println("*** GOARCH:", runtime.GOARCH, "***")
	if GOARM := os.Getenv("GOARM"); len(GOARM) > 0 {
		fmt.Println("*** GOARM:", GOARM, "***")
	}
	if GOMIPS := os.Getenv("GOMIPS"); len(GOMIPS) > 0 {
		fmt.Println("*** GOMIPS:", GOMIPS, "***")
	}
	//
	gt := getgt()
	runtime.GC()
	assert.Equal(t, "g", gt.Name())
	//
	numField := gt.NumField()
	//
	fmt.Println("#numField:", numField)
	fmt.Println("#offsetGoid:", offsetGoid)
	fmt.Println("#offsetPaniconfault:", offsetPaniconfault)
	fmt.Println("#offsetGopc:", offsetGopc)
	fmt.Println("#offsetLabels:", offsetLabels)
	//
	assert.Greater(t, numField, 20)
	assert.Greater(t, int(offsetGoid), 0)
	assert.Greater(t, int(offsetPaniconfault), 0)
	assert.Greater(t, int(offsetGopc), 0)
	assert.Greater(t, int(offsetLabels), 0)
	//
	runTest(t, func() {
		tt := getgt()
		runtime.GC()
		assert.Equal(t, numField, tt.NumField())
		assert.Equal(t, offsetGoid, offset(tt, "goid"))
		assert.Equal(t, offsetPaniconfault, offset(tt, "paniconfault"))
		assert.Equal(t, offsetGopc, offset(tt, "gopc"))
		assert.Equal(t, offsetLabels, offset(tt, "labels"))
	})
}

func runTest(t *testing.T, fun func()) {
	var count int32
	wg := &sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				fun()
			}
			atomic.AddInt32(&count, 1)
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, 10, int(count))
}
