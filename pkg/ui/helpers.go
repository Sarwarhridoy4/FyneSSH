package ui

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/Sarwarhridoy4/FyneSSH/internal/platform"
	"github.com/Sarwarhridoy4/FyneSSH/internal/sshclient"
)

func uploadPublicKey(pubKeyPath, user, host, port, password string) error {
	portNum := parsePort(port)

	authMethods := []ssh.AuthMethod{ssh.Password(password)}

	client, err := sshclient.Dial(context.Background(), user, host, portNum, authMethods)
	if err != nil {
		return fmt.Errorf("connect to %s@%s:%s: %w", user, host, port, err)
	}
	defer client.Close()

	session, err := client.Session()
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	defer session.Close()

	pubKeyBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return fmt.Errorf("read public key file: %w", err)
	}

	pubKeyContent := strings.TrimSpace(string(pubKeyBytes))
	escapedKey := strings.ReplaceAll(pubKeyContent, "'", "'\\''")

	shellCmd := fmt.Sprintf("mkdir -p ~/.ssh && chmod 700 ~/.ssh && echo '%s' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys", escapedKey)

	if err := session.Run(shellCmd); err != nil {
		return fmt.Errorf("append public key: %w", err)
	}

	return nil
}

func addToKnownHostsUI(host string) error {
	portNum := 22
	authMethods := []ssh.AuthMethod{}

	hostKeyCallback, err := platform.HostKeyCallback()
	if err != nil {
		return fmt.Errorf("create host key callback: %w", err)
	}

	client, err := sshclient.DialWithHostKey(context.Background(), platform.CurrentUser(), host, portNum, authMethods, hostKeyCallback)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", host, err)
	}
	defer client.Close()

	return nil
}

func parsePort(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 22
	}
	var p int
	_, err := fmt.Sscan(s, &p)
	if err != nil {
		return 22
	}
	return p
}
