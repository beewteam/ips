package main

import (
	"net"
	"time"

	"github.com/stretchr/testify/mock"
)

type MAddr struct {
	mock.Mock
}

func (m MAddr) Network() string {
	return "test"
}

func (m MAddr) String() string {
	return "test:test"
}

type connMockObject struct {
	mock.Mock
}

func (m connMockObject) Read(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m connMockObject) Write(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m connMockObject) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m connMockObject) LocalAddr() net.Addr {
	args := m.Called()
	_ = args
	return MAddr{}
}

func (m connMockObject) RemoteAddr() net.Addr {
	args := m.Called()
	_ = args
	return MAddr{}
}

func (m connMockObject) SetDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m connMockObject) SetReadDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m connMockObject) SetWriteDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}
