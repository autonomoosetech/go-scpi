package scpi

import (
	"context"
	"net"
	"time"

	"golang.org/x/xerrors"
)

// Client is a client of a device controlled using SCPI commands.
type Client interface {
	// Close closes the connection.
	Close() error

	// Exec executes a SCPI command.
	Exec(cmd string) error

	// ExecContext executes a SCPI command.
	ExecContext(ctx context.Context, cmd string) error

	// Ping verifies the connection to the device is still alive,
	// establishing a connection if necessary.
	Ping() error

	// PingContext verifies the connection to the device is still alive,
	// establishing a connection if necessary.
	PingContext(ctx context.Context) error

	// Query queries the device for the results of the specified command.
	Query(cmd string) (res []byte, err error)

	// QueryContext queries the device for the results of the specified command.
	QueryContext(ctx context.Context, cmd string) (res []byte, err error)
}

// NewClient returns a new client of a device controlled using SCPI commands.
func NewClient(proto, addr string, timeout time.Duration) (Client, error) {
	switch proto {
	case "tcp":
		return newTCPClient(addr, timeout)
	default:
		return nil, xerrors.Errorf("invalid protocol %s", proto)
	}
}

// TCPClient is an implementation of the Client interface for TCP network connections.
type TCPClient struct {
	conn *net.TCPConn
}

func newTCPClient(addr string, timeout time.Duration) (*TCPClient, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	d := net.Dialer{
		Timeout: timeout,
	}
	conn, err := d.Dial("tcp", tcpAddr.String())
	if err != nil {
		return nil, err
	}
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, xerrors.Errorf("failed to case %T to *net.TCPConn", conn)
	}
	client := &TCPClient{
		conn: tcpConn,
	}
	return client, nil
}

// Close implements the Client Close method.
func (c *TCPClient) Close() error {
	return c.conn.Close()
}

// Exec implements the Client Exec method.
func (c *TCPClient) Exec(cmd string) error {
	return c.ExecContext(context.Background(), cmd)
}

// ExecContext implements the Client ExecContext method.
func (c *TCPClient) ExecContext(ctx context.Context, cmd string) error {
	b := []byte(cmd + "\n")
	if _, err := c.conn.Write(b); err != nil {
		return err
	}
	return nil
}

// Ping implements the Client Ping method.
func (c *TCPClient) Ping() error {
	return c.PingContext(context.Background())
}

// PingContext implements the Client PingContext method.
func (c *TCPClient) PingContext(ctx context.Context) error {
	// BUG(scizorman): PingContext is not implemented yet.
	return nil
}

// Query implements the Client Query method.
func (c *TCPClient) Query(cmd string) (res []byte, err error) {
	return c.QueryContext(context.Background(), cmd)
}

// QueryContext implements the Client QueryContext method.
func (c *TCPClient) QueryContext(ctx context.Context, cmd string) (res []byte, err error) {
	if err := c.ExecContext(ctx, cmd); err != nil {
		return nil, err
	}

	res = make([]byte, 1024)
	l, err := c.conn.Read(res)
	if err != nil {
		return nil, err
	}
	return res[:l], nil
}
