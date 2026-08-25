package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	apiproxy "github.com/PastureStack/cli/internal/dockerapiproxy"
	rancher "github.com/PastureStack/cli/internal/rancherclient/v2"
	"github.com/sirupsen/logrus"
)

const (
	DefaultNamedPipeline = "//./pipe/docker-sock"
	PipelineNamePrefix   = "npipe://" + DefaultNamedPipeline
)

func getDockerHost(client *rancher.RancherClient, host string) (string, *apiproxy.Proxy, error) {
	dockerHost := getRandNPipe()
	logrus.Info(dockerHost)
	proxy := apiproxy.NewProxy(client, host, dockerHost)
	return dockerHost, proxy, nil
}

func getRandNPipe() string {
	var suffix [6]byte
	if _, err := rand.Read(suffix[:]); err == nil {
		return PipelineNamePrefix + hex.EncodeToString(suffix[:])
	}
	return fmt.Sprintf(PipelineNamePrefix+"%d", time.Now().UnixNano())
}
