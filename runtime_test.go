package routine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestGetgp(t *testing.T) {
	gp1 := getgp()
	runtime.GC()
	assert.NotNil(t, gp1, "Fail to get gp.")

	runTest(t, func() {
		gp2 := getgp()
		runtime.GC()
		assert.NotNil(t, gp2, "Fail to get gp.")
		assert.NotEqual(t, gp1, gp2, "Every living g must be different. [gp1:%p] [gp2:%p]", gp1, gp2)
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
	fmt.Println("#offsetLabels:", offsetLabels)
	//
	assert.Greater(t, numField, 20)
	assert.Greater(t, int(offsetGoid), 0)
	assert.Greater(t, int(offsetPaniconfault), 0)
	assert.Greater(t, int(offsetLabels), 0)
	//
	runTest(t, func() {
		tt := getgt()
		runtime.GC()
		switch runtime.GOARCH {
		case "386":
			fallthrough
		case "amd64":
			fallthrough
		case "arm":
			fallthrough
		case "arm64":
			assert.Equal(t, numField, tt.NumField())
			assert.Equal(t, offsetGoid, offset(tt, "goid"))
			assert.Equal(t, offsetPaniconfault, offset(tt, "paniconfault"))
			assert.Equal(t, offsetLabels, offset(tt, "labels"))

		default:
			panic("Not support GOARCH: " + runtime.GOARCH)
		}
	})
}

func runTest(t *testing.T, fun func()) {
	var count int32
	wg := &sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				fun()
			}
			atomic.AddInt32(&count, 1)
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, 100, int(count))
}
