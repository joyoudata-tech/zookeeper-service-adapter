package main

import (
	"os"
	"log"

	"github.com/joyoudata-tech/zookeeper-service-adapter/adapter"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

func main() {
	zkCommand :=  os.Getenv("ZK_COMMAND")
	if zkCommand == "" {
		zkCommand = "/var/vcap/packages/zk-client/zkClient"
	}

	stderrLogger := log.New(os.Stderr, "[zookeeper-service-adapter]", log.LstdFlags)

	manifestGenerator := &adapter.ManifestGenerator{
		StderrLogger: stderrLogger,
	}

	binder := &adapter.Binder{
		CommandRunner: adapter.ExternalCommandRunner{},
		ZKCommand: zkCommand,
		StderrLogger: stderrLogger,
	}

	serviceadapter.HandleCommandLineInvocation(os.Args, manifestGenerator, binder, &adapter.DashboardUrlGenerator{})
}
