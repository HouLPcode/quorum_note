package tests

import (
	"testing"
	"time"
	"fmt"
	"runtime"
)

// go的调度不是抢占的
func TestTicker(t *testing.T){
	runtime.GOMAXPROCS(1)
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("ticker!!!")
		}
	}()

	for i := 0; ;{
		i++
	}

}
