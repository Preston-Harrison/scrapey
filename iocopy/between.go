package iocopy

import (
	"io"
	"sync"
)

func Between(a, b io.ReadWriteCloser) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if _, err := io.Copy(b, a); err != nil {
			return
		}
		a.Close()
	}()
	go func() {
		defer wg.Done()
		if _, err := io.Copy(a, b); err != nil {
			return
		}
		b.Close()
	}()
	wg.Wait()
}
