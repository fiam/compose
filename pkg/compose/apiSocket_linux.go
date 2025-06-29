//go:build linux

package compose

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/opts"
)

func dockerEndpoint(cli command.Cli) (string, error) {
	isDockerdRunning, err := isProcessRunning("dockerd")
	if err != nil {
		return "", fmt.Errorf("checking if dockerd is running: %w", err)
	}
	if !isDockerdRunning {
		// Assume we are using a remote Docker Engine
		return opts.DefaultUnixSocket, nil
	}
	socket := cli.DockerEndpoint().Host
	if !strings.HasPrefix(socket, "unix://") {
		return "", fmt.Errorf("use_api_socket can only be used with unix sockets: docker endpoint %s is incompatible", socket)
	}
	socket = strings.TrimPrefix(socket, "unix://") // should we confirm absolute path?
	return socket, nil
}

func isProcessRunning(target string) (bool, error) {
	const procDir = "/proc"

	entries, err := os.ReadDir(procDir)
	if err != nil {
		return false, err
	}

	targetBytes := []byte(target + "\n")
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name()[0] < '0' || entry.Name()[0] > '9' {
			continue // skip non-PID entries
		}

		cmdPath := filepath.Join(procDir, entry.Name(), "comm")
		cmdBytes, err := os.ReadFile(cmdPath)
		if err != nil {
			continue // permission error or short-lived process
		}
		if bytes.Equal(cmdBytes, targetBytes) {
			return true, nil
		}
	}

	return false, nil
}
