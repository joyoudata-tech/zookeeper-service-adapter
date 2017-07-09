package adapter

import (
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
)

type DashboardUrlGenerator struct {
}

func (a *DashboardUrlGenerator) DashboardUrl(instanceID string, plan serviceadapter.Plan, manifest bosh.BoshManifest) (serviceadapter.DashboardUrl, error) {
	return serviceadapter.DashboardUrl{DashboardUrl: "http://zookeeper.com/" + instanceID}, nil
}
