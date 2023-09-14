package iocopy

import (
	"io"
	"sync"
)

type RWCloser interface {
	io.Reader
	io.Writer
	io.Closer
}

func Between(a, b RWCloser) {
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
