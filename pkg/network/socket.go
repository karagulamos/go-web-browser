package network

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
)

type Socket interface {
	Connect(address string) error
	Writef(format string, args ...interface{}) error
	WritefLine(format string, args ...interface{}) error
	ReadLine() (string, error)

	io.ReadCloser
}

type SocketType string

const (
	Tcp SocketType = "tcp"
	Udp SocketType = "udp"
)

type socket struct {
	socketType SocketType
	secure     bool
	conn       net.Conn
}

type Config struct {
	Secure bool
}

func NewSocket(socketType SocketType, conf Config) Socket {
	s := &socket{
		socketType: socketType,
		secure:     conf.Secure,
	}

	return s
}

func (s *socket) Connect(address string) error {
	var err error

	if s.secure {
		s.conn, err = tls.Dial(string(s.socketType), address, &tls.Config{})
	} else {
		s.conn, err = net.Dial(string(s.socketType), address)
	}

	return err
}

func (s *socket) Writef(format string, args ...interface{}) error {
	_, err := s.conn.Write([]byte(fmt.Sprintf(format, args...)))
	return err
}

func (s *socket) WritefLine(format string, args ...interface{}) error {
	if err := s.Writef(format, args...); err != nil {
		return err
	}

	return s.Writef("\r\n")
}

func (s *socket) ReadLine() (string, error) {
	var line []byte
	var err error

	for char := make([]byte, 1); ; {
		if _, err = s.Read(char); err != nil {
			break
		}

		if char[0] == '\r' {
			if _, err = s.Read(char); err != nil {
				break
			}
			if char[0] == '\n' {
				break
			}
		}

		line = append(line, char[0])
	}

	return string(line), nil
}

func (s *socket) Read(p []byte) (int, error) {
	return s.conn.Read(p)
}

func (s *socket) Close() error {
	return s.conn.Close()
}
