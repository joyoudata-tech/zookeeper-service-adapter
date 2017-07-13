package adapter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

func (b *Binder) DeleteBinding(bindingID string, boshVMs bosh.BoshVMs, manifest bosh.BoshManifest, requestParams serviceadapter.RequestParameters) error {
	zookeeperServers := boshVMs["zookeeper"]
	if len(zookeeperServers) == 0 {
		b.StderrLogger.Println("no VMs for job peers")
		return errors.New("")
	}
	if _, errorStream, err := b.Run(b.ZKCommand, "--zookeeperServers", strings.Join(zookeeperServers, ","), "delete", "/" + bindingID); err != nil {
		if strings.Contains(string(errorStream), "node does not exist") {
			b.StderrLogger.Println(fmt.Sprintf("director '%s' not found", bindingID))
			return serviceadapter.NewBindingNotFoundError(nil)
		}
		b.StderrLogger.Println("Error deleting director: " + err.Error())
		return errors.New("")
	}

	return nil
}
