package adapter

import (
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"errors"
	"fmt"
)

func (b *Binder) CreateBinding(bindingId string, boshVMs bosh.BoshVMs, manifest bosh.BoshManifest, requestParams serviceadapter.RequestParameters) (serviceadapter.Binding, error){
	zookeeperHosts := boshVMs["peers"]
	if len(zookeeperHosts) == 0 {
		b.StderrLogger.Println("no VMs for instance group peers")
		return serviceadapter.Binding{}, errors.New("")
	}
	var zookeeperAddress []interface{}
	for _, zookeeperHost := range zookeeperHosts {
		zookeeperAddress = append(zookeeperAddress, fmt.Sprintf("%s:2181", zookeeperHost))
	}
	return serviceadapter.Binding{
		Credentials: map[string]interface{}{
			"zookeeper_services" : zookeeperAddress,
		},
	},nil
}
