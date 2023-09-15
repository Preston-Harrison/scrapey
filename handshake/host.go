package handshake

import (
	"fmt"
	"io"
	"net"
)

// Sends the host as bytes, and wait for a response. The response must
// be exactly 1 byte.
// Returns nil if the host was sent and an ok status was received.
func SendHost(w io.ReadWriter, host string) error {
	_, err := w.Write([]byte(host))
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

// Receives a host, tries to dial the host over TCP, and replies either
// ok or notOk if the host could be received and a connection could
// be made.
func ReceiveAndDialHost(w io.ReadWriter) (net.Conn, error) {
	hostBytes := make([]byte, 1024)
	i, err := w.Read(hostBytes)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", string(hostBytes[:i]))
	if err != nil {
		if _, writeErr := w.Write([]byte{notOk}); writeErr != nil {
			return nil, writeErr
		}
		return nil, err
	}

	if _, err := w.Write([]byte{ok}); err != nil {
		return nil, err
	}

	return conn, nil
}
