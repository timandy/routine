package main

import (
	"fmt"
	"github.com/timandy/routine"
	"time"
)

var nameVar = routine.NewInheritableThreadLocal()

func main() {
	nameVar.Set("hello world")
	fmt.Println("name: ", nameVar.Get())

	// other goroutine cannot read it
	go func() {
		fmt.Println("name1: ", nameVar.Get())
	}()

	// but, the new goroutine could inherit/copy all local data from the current goroutine like this:
	routine.Go(func() {
		fmt.Println("name2: ", nameVar.Get())
	})

	time.Sleep(time.Second)
}
