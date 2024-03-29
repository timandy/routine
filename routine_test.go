package routine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInheritedTask(t *testing.T) {
	tls := NewInheritableThreadLocal[string]()
	it := inheritedTask{context: nil, function: func() {
		assert.Equal(t, "", tls.Get())
	}}
	task := NewFutureTask(it.run)
	go task.Run()
	assert.Nil(t, task.Get())
	assert.True(t, task.IsDone())
	//
	it2 := inheritedTask{context: nil, function: func() {
		assert.Equal(t, "", tls.Get())
	}}
	task2 := NewFutureTask(it2.run)
	go func() {
		tls.Set("hello")
		task2.Run()
	}()
	assert.Nil(t, task2.Get())
	assert.True(t, task2.IsDone())
	//
	tls.Set("world")
	it3 := inheritedTask{context: createInheritedMap(), function: func() {
		assert.Equal(t, "world", tls.Get())
	}}
	task3 := NewFutureTask(it3.run)
	go task3.Run()
	assert.Nil(t, task3.Get())
	assert.True(t, task3.IsDone())
	//
	it4 := inheritedTask{context: createInheritedMap(), function: func() {
		assert.Equal(t, "world", tls.Get())
	}}
	task4 := NewFutureTask(it4.run)
	go func() {
		tls.Set("hello")
		task4.Run()
	}()
	assert.Nil(t, task4.Get())
	assert.True(t, task4.IsDone())
}

func TestInheritedWaitTask(t *testing.T) {
	tls := NewInheritableThreadLocal[string]()
	it := inheritedWaitTask{context: nil, function: func(token CancelToken) {
		assert.Equal(t, "", tls.Get())
	}}
	task := NewFutureTask(it.run)
	go task.Run()
	assert.Nil(t, task.Get())
	assert.True(t, task.IsDone())
	//
	it2 := inheritedWaitTask{context: nil, function: func(token CancelToken) {
		assert.Equal(t, "", tls.Get())
	}}
	task2 := NewFutureTask(it2.run)
	go func() {
		tls.Set("hello")
		task2.Run()
	}()
	assert.Nil(t, task2.Get())
	assert.True(t, task2.IsDone())
	//
	tls.Set("world")
	it3 := inheritedWaitTask{context: createInheritedMap(), function: func(token CancelToken) {
		assert.Equal(t, "world", tls.Get())
	}}
	task3 := NewFutureTask(it3.run)
	go task3.Run()
	assert.Nil(t, task3.Get())
	assert.True(t, task3.IsDone())
	//
	it4 := inheritedWaitTask{context: createInheritedMap(), function: func(token CancelToken) {
		assert.Equal(t, "world", tls.Get())
	}}
	task4 := NewFutureTask(it4.run)
	go func() {
		tls.Set("hello")
		task4.Run()
	}()
	assert.Nil(t, task4.Get())
	assert.True(t, task4.IsDone())
}

func TestInheritedWaitResultTask(t *testing.T) {
	tls := NewInheritableThreadLocal[string]()
	it := inheritedWaitResultTask[string]{context: nil, function: func(token CancelToken) string {
		assert.Equal(t, "", tls.Get())
		return tls.Get()
	}}
	task := NewFutureTask(it.run)
	go task.Run()
	assert.Equal(t, "", task.Get())
	assert.True(t, task.IsDone())
	//
	it2 := inheritedWaitResultTask[string]{context: nil, function: func(token CancelToken) string {
		assert.Equal(t, "", tls.Get())
		return tls.Get()
	}}
	task2 := NewFutureTask(it2.run)
	go func() {
		tls.Set("hello")
		task2.Run()
	}()
	assert.Equal(t, "", task2.Get())
	assert.True(t, task2.IsDone())
	//
	tls.Set("world")
	it3 := inheritedWaitResultTask[string]{context: createInheritedMap(), function: func(token CancelToken) string {
		assert.Equal(t, "world", tls.Get())
		return tls.Get()
	}}
	task3 := NewFutureTask(it3.run)
	go task3.Run()
	assert.Equal(t, "world", task3.Get())
	assert.True(t, task3.IsDone())
	//
	it4 := inheritedWaitResultTask[string]{context: createInheritedMap(), function: func(token CancelToken) string {
		assert.Equal(t, "world", tls.Get())
		return tls.Get()
	}}
	task4 := NewFutureTask(it4.run)
	go func() {
		tls.Set("hello")
		task4.Run()
	}()
	assert.Equal(t, "world", task4.Get())
	assert.True(t, task4.IsDone())
}
