package adapter

import (
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"

	"errors"
	"fmt"
	"bytes"
	"os/exec"
	"strings"
	"sort"
)

func (b *Binder) CreateBinding(bindingId string, boshVMs bosh.BoshVMs, manifest bosh.BoshManifest, requestParams serviceadapter.RequestParameters) (serviceadapter.Binding, error){
	//用户自定义根目录
	params := requestParams.ArbitraryParams()
	var invalidParams []string
	for paramKey, _ := range params {
		if paramKey != "director" {
			invalidParams = append(invalidParams, paramKey)
		}
	}
	if len(invalidParams) > 0 {
		sort.Strings(invalidParams)
		errorMessage := fmt.Sprintf("unsupported parameter(s) for this service: %s", strings.Join(invalidParams, ", "))
		b.StderrLogger.Println(errorMessage)
		return serviceadapter.Binding{}, errors.New(errorMessage)
	}

	//判断zookeeper vm是否存在
	zookeeperHosts := boshVMs["peers"]
	if len(zookeeperHosts) == 0 {
		b.StderrLogger.Println("no VMs for instance group peers")
		return serviceadapter.Binding{}, errors.New("")
	}

	//获取zookeeper的地址
	var zookeeperAddress []interface{}
	for _, zookeeperHost := range zookeeperHosts {
		zookeeperAddress = append(zookeeperAddress, fmt.Sprintf("%s:2181", zookeeperHost))
	}

	//创建director目录
	var director string

	director = "/" + bindingId
	if _, errorStream, err := b.Run(b.ZKCommand, "--zookeeperServers", strings.Join(zookeeperHosts, ","), "create", director) ; err != nil {
		if strings.Contains(string(errorStream), "node already exists") {
			b.StderrLogger.Println(fmt.Sprintf("director '%s' already exists", director))
			return serviceadapter.Binding{}, serviceadapter.NewBindingAlreadyExistsError(nil)
		}
		b.StderrLogger.Println("Error creating director: " + err.Error())
		return serviceadapter.Binding{}, errors.New("")
	}

	if params["director"] != nil {//带参数的
		director = params["director"]
		if !strings.HasPrefix(director, "/") {
			director = "/" + director
		}
		if _, _, err := b.Run(b.ZKCommand, "--zookeeperServers", strings.Join(zookeeperHosts, ","), "create", params["director"].(string)); err != nil {
			b.StderrLogger.Println("Error creating director: " + err.Error())
			return serviceadapter.Binding{}, errors.New("")
		}
	}

	return serviceadapter.Binding{
		Credentials: map[string]interface{}{
			"zookeeper_services" : zookeeperAddress,
			"director" : director,
		},
	},nil
}

//go:generate counterfeiter -o fake_command_runner/fake_command_runner.go . CommandRunner
type CommandRunner interface {
	Run(name string, arg ...string) ([]byte, []byte, error)
}

type ExternalCommandRunner struct{}

func (c ExternalCommandRunner) Run(name string, arg ...string) ([]byte, []byte, error) {
	cmd := exec.Command(name, arg...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.Output()
	return stdout, stderr.Bytes(), err
}