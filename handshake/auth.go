package handshake

import (
	"fmt"
	"io"
)

func SendAuthToken(w io.ReadWriter, token string) error {
	_, err := w.Write([]byte(token))
	if err != nil {
		return err
	}

	status := make([]byte, 1)
	_, err = w.Read(status)
	if err != nil {
		return err
	}

	if status[0] != ok {
		return fmt.Errorf("received non-ok status code %d", status[0])
	}

	return nil
}

func ReceiveAuthToken(w io.ReadWriter) (string, error) {
	tokenBytes := make([]byte, 1024)
	i, err := w.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return string(tokenBytes[:i]), nil
}

func SendAuthStatus(w io.ReadWriter, isValidAuth bool) (err error) {
	if isValidAuth {
		_, err = w.Write([]byte{ok})
	} else {
		_, err = w.Write([]byte{notOk})
	}
	return
}
