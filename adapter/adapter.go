package adapter

import (
	"log"

	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

type ManifestGenerator struct {
	StderrLogger *log.Logger
}

type Binder struct {
	ZKCommand string
	CommandRunner
	StderrLogger *log.Logger
}

var InstanceGroupMapper = serviceadapter.GenerateInstanceGroupsWithNoProperties