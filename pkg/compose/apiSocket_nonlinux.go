//go:build !linux

package compose

import (
	"github.com/docker/cli/cli/command"
	"github.com/docker/docker/opts"
)

func dockerEndpoint(_ command.Cli) (string, error) {
	// non-Linux systems cannot run dockerd directly, so they must communicate with
	// a remote engine. Make a best effort to return the default socket.
	return opts.DefaultUnixSocket, nil
}
