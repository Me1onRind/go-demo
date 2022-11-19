package goroutine

import (
	"context"
	"testing"
)

func Test_SafeGo(t *testing.T) {
	ch := make(chan struct{}, 1)

	f := func() {
		defer func() {
			ch <- struct{}{}
		}()
		panic("test")
	}

	SafeGo(context.Background(), f)
	<-ch
}
