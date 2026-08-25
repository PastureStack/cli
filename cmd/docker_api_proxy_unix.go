//go:build !windows
// +build !windows

package cmd

import (
	"os"

	apiproxy "github.com/PastureStack/cli/internal/dockerapiproxy"
	rancher "github.com/PastureStack/cli/internal/rancherclient/v2"
)

func getDockerHost(client *rancher.RancherClient, host string) (string, *apiproxy.Proxy, error) {
	tempfile, err := os.CreateTemp("", "docker-sock")
	if err != nil {
		return "", nil, err
	}
	defer os.Remove(tempfile.Name())

	if err := tempfile.Close(); err != nil {
		return "", nil, err
	}

	dockerHost := "unix://" + tempfile.Name()
	proxy := apiproxy.NewProxy(client, host, dockerHost)
	return dockerHost, proxy, nil
}
