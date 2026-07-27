package sshclient

import (
	"context"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/ssh"
)

// Client wraps an SSH connection and session.
type Client struct {
	config *ssh.ClientConfig
	conn   *ssh.Client
}

// Dial connects to a remote SSH server using the provided configuration.
func Dial(ctx context.Context, user, host string, port int, auth []ssh.AuthMethod) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	cfg := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: allow known_hosts
	}

	// TODO: use a net.Dialer with context cancellation instead of DialTimeout.
	conn, err := ssh.Dial("tcp", addr, cfg)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", addr, err)
	}
	return &Client{config: cfg, conn: conn}, nil
}

// Session creates a new remote session.
func (c *Client) Session() (*ssh.Session, error) {
	sess, err := c.conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("new session: %w", err)
	}
	return sess, nil
}

// RunCommand executes a command remotely and streams stdout/stderr to writers.
func (c *Client) RunCommand(ctx context.Context, cmd string, stdout, stderr io.Writer) error {
	sess, err := c.Session()
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.Stdout = stdout
	sess.Stderr = stderr
	sess.Stdin = nil

	return sess.Run(cmd)
}

// Close terminates the underlying SSH connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
