package iocopy

import (
	"fmt"
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
			fmt.Println("Ws to client error", err)
			return
		}
		a.Close()
	}()
	go func() {
		defer wg.Done()
		if _, err := io.Copy(a, b); err != nil {
			fmt.Println("Client to ws error", err)
			return
		}
		b.Close()
	}()
	wg.Wait()
}
