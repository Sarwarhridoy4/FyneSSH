package platform

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
)

// KnownHostsPath returns the path to the known_hosts file.
func KnownHostsPath() string {
	return filepath.Join(SSHBasePath(), "known_hosts")
}

// HostBlock represents a single Host entry in the SSH config.
type HostBlock struct {
	Patterns           []string
	HostName           string
	User               string
	Port               int
	IdentityFile       string
	AddKeysToAgent     bool
	IdentitiesOnly     bool
	ServerAliveInterval int
	ServerAliveCountMax int
	Comment            string
}

var hostHeaderRe = regexp.MustCompile(`(?i)^\s*Host\s+(.+)$`)
var optionRe = regexp.MustCompile(`(?i)^\s+(HostName|User|Port|IdentityFile|AddKeysToAgent|IdentitiesOnly|ServerAliveInterval|ServerAliveCountMax)\s+(.+)$`)

// ReadSSHConfig reads and parses the SSH config file.
func ReadSSHConfig(path string) ([]HostBlock, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var blocks []HostBlock
	var current *HostBlock
	var currentComment strings.Builder
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			if current != nil {
				current.Comment = strings.TrimSpace(currentComment.String())
			}
			currentComment.WriteString(line)
			currentComment.WriteByte('\n')
			continue
		}

		if h := hostHeaderRe.FindStringSubmatch(line); h != nil {
			if current != nil {
				blocks = append(blocks, *current)
			}
			current = &HostBlock{
				Patterns: strings.Fields(h[1]),
				Comment:  strings.TrimSpace(currentComment.String()),
			}
			currentComment.Reset()
			continue
		}

		if o := optionRe.FindStringSubmatch(line); o != nil && current != nil {
			key := strings.ToLower(o[1])
			val := strings.TrimSpace(o[2])
			val = strings.Trim(val, "\"'")
			switch key {
			case "hostname":
				current.HostName = val
			case "user":
				current.User = val
			case "port":
				var p int
				if _, err := fmt.Sscan(val, &p); err == nil {
					current.Port = p
				}
			case "identityfile":
				current.IdentityFile = val
			case "addkeystoagent":
				current.AddKeysToAgent = strings.EqualFold(val, "yes")
			case "identitiesonly":
				current.IdentitiesOnly = strings.EqualFold(val, "yes")
			case "serveraliveinterval":
				var v int
				if _, err := fmt.Sscan(val, &v); err == nil {
					current.ServerAliveInterval = v
				}
			case "serveralivecountmax":
				var v int
				if _, err := fmt.Sscan(val, &v); err == nil {
					current.ServerAliveCountMax = v
				}
			}
		}
	}
	if current != nil {
		current.Comment = strings.TrimSpace(currentComment.String())
		blocks = append(blocks, *current)
	}

	return blocks, scanner.Err()
}

// WriteSSHConfig writes the SSH config blocks to the file.
func WriteSSHConfig(path string, blocks []HostBlock) error {
	var b strings.Builder
	for i, block := range blocks {
		if block.Comment != "" {
			b.WriteString(block.Comment)
			if !strings.HasSuffix(block.Comment, "\n") {
				b.WriteByte('\n')
			}
		}
		b.WriteString("Host ")
		b.WriteString(strings.Join(block.Patterns, " "))
		b.WriteByte('\n')
		if block.HostName != "" {
			b.WriteString("    HostName ")
			b.WriteString(block.HostName)
			b.WriteByte('\n')
		}
		if block.User != "" {
			b.WriteString("    User ")
			b.WriteString(block.User)
			b.WriteByte('\n')
		}
		if block.Port > 0 {
			b.WriteString("    Port ")
			b.WriteString(fmt.Sprintf("%d", block.Port))
			b.WriteByte('\n')
		}
		if block.IdentityFile != "" {
			b.WriteString("    IdentityFile ")
			b.WriteString(block.IdentityFile)
			b.WriteByte('\n')
		}
		if block.AddKeysToAgent {
			b.WriteString("    AddKeysToAgent yes\n")
		}
		if block.IdentitiesOnly {
			b.WriteString("    IdentitiesOnly yes\n")
		}
		if block.ServerAliveInterval > 0 {
			b.WriteString("    ServerAliveInterval ")
			b.WriteString(fmt.Sprintf("%d", block.ServerAliveInterval))
			b.WriteByte('\n')
		}
		if block.ServerAliveCountMax > 0 {
			b.WriteString("    ServerAliveCountMax ")
			b.WriteString(fmt.Sprintf("%d", block.ServerAliveCountMax))
			b.WriteByte('\n')
		}
		if i < len(blocks)-1 {
			b.WriteByte('\n')
		}
	}
	return writeFile(path, []byte(b.String()), 0600)
}

// UpdateOrAddHost updates or adds a Host block to the SSH config.
func UpdateOrAddHost(path string, block HostBlock) error {
	blocks, err := ReadSSHConfig(path)
	if err != nil {
		return err
	}

	found := false
	for i := range blocks {
		for _, p := range blocks[i].Patterns {
			if p == block.Patterns[0] {
				blocks[i] = block
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		blocks = append(blocks, block)
	}

	return WriteSSHConfig(path, blocks)
}

// KnownHostEntry represents a single known_hosts entry.
type KnownHostEntry struct {
	Hostnames []string
	Key       ssh.PublicKey
}

// ReadKnownHosts reads the known_hosts file.
func ReadKnownHosts(path string) ([]KnownHostEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open known_hosts: %w", err)
	}
	defer f.Close()

	var entries []KnownHostEntry
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		fields := strings.Split(line, " ")
		if len(fields) < 3 {
			continue
		}

		var hostnames []string
		var key ssh.PublicKey

		if fields[0] == "|1|" {
			hostnames = strings.Split(fields[1], ",")
			key, err = ssh.ParsePublicKey([]byte(strings.Join(fields[2:], " ")))
			if err != nil {
				continue
			}
		} else {
			hostnames = strings.Split(fields[0], ",")
			key, err = ssh.ParsePublicKey([]byte(strings.Join(fields[1:], " ")))
			if err != nil {
				continue
			}
		}

		entries = append(entries, KnownHostEntry{
			Hostnames: hostnames,
			Key:       key,
		})
	}

	return entries, scanner.Err()
}

// WriteKnownHosts writes known host entries to the file.
func WriteKnownHosts(path string, entries []KnownHostEntry) error {
	var b strings.Builder
	for _, entry := range entries {
		b.WriteString(strings.Join(entry.Hostnames, ","))
		b.WriteByte(' ')
		b.Write(ssh.MarshalAuthorizedKey(entry.Key))
	}
	return writeFile(path, []byte(b.String()), 0600)
}

// AddToKnownHosts adds a hostname and public key to known_hosts.
func AddToKnownHosts(path, hostname string, key ssh.PublicKey) error {
	entries, err := ReadKnownHosts(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		for _, h := range e.Hostnames {
			if h == hostname {
				return nil
			}
		}
	}

	entries = append(entries, KnownHostEntry{
		Hostnames: []string{hostname},
		Key:       key,
	})

	return WriteKnownHosts(path, entries)
}

// HostKeyCallback returns an ssh.HostKeyCallback that checks and adds hosts to known_hosts.
func HostKeyCallback() (ssh.HostKeyCallback, error) {
	path := KnownHostsPath()
	entries, err := ReadKnownHosts(path)
	if err != nil {
		return nil, err
	}

	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		for _, e := range entries {
			for _, h := range e.Hostnames {
				if h == hostname {
					if string(e.Key.Marshal()) != string(key.Marshal()) {
						return fmt.Errorf("host key mismatch for %s", hostname)
					}
					return nil
				}
			}
		}

		entries = append(entries, KnownHostEntry{
			Hostnames: []string{hostname},
			Key:       key,
		})

		return WriteKnownHosts(path, entries)
	}, nil
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

// OpenTerminal opens the system terminal and runs the given command.
func OpenTerminal(command string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s <<< \"%s\"", GetDefaultShell(), strings.ReplaceAll(command, "\"", "\\\"")))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}
