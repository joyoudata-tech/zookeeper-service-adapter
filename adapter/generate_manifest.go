package adapter

import (
        "github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"

	"errors"
	"fmt"
	"strings"
)

const OnlystemcellAlias = "only-stemcell"

//manifest的job任务清单
func defaultDeploymentInstanceGroupsToJobs() map[string][]string {
	return map[string][]string{
		"zookeeper": []string{"zookeeper"},
	}
}

//创建清单
func (a *ManifestGenerator) GenerateManifest(
	serviceDeployment serviceadapter.ServiceDeployment,
	servicePlan serviceadapter.Plan,
	requestParams serviceadapter.RequestParameters,
	previousManifest *bosh.BoshManifest,
	previousPlan *serviceadapter.Plan,
) (bosh.BoshManifest, error) {
	//判断有无预先设定的计划
	if previousPlan != nil {
		prev := instanceCounts(*previousPlan)
		current := instanceCounts(servicePlan)
		if (prev["zookeeper"] > current["zookeeper"]) {//有多少服务写多少服务
			a.StderrLogger.Println("the current service plan is too small.")
			return bosh.BoshManifest{}, errors.New("")
		}
	}
	//创建相关的service服务release
	var releases []bosh.Release
	loggingRaw, ok := servicePlan.Properties["logging"]
	includeMetron := false
	if ok {
		includeMetron = true
	}
	for _, serviceRelease := range serviceDeployment.Releases {
		releases = append(releases, bosh.Release{
			Name:	 serviceRelease.Name,
			Version: serviceRelease.Version,
		})
	}
	deploymentinstanceGroupsToJobs := defaultDeploymentInstanceGroupsToJobs()
	//在job里加入metron_agent,也可以选择不加
	if includeMetron {
		for instanceGroup, jobs := range deploymentinstanceGroupsToJobs {
			deploymentinstanceGroupsToJobs[instanceGroup] = append(jobs, "metron_agent")
		}
	}
	//检查zookeeper的peers实例是否存在
	err := checkInstanceGroupsPresent([]string{"zookeeper"}, servicePlan.InstanceGroups)
	if err != nil {
		a.StderrLogger.Println(err.Error())
		return bosh.BoshManifest{}, errors.New("Contact your operator, service configuration issue occurred")
	}
	//服务实例组映射
	instanceGroups ,err := InstanceGroupMapper(servicePlan.InstanceGroups, serviceDeployment.Releases, OnlystemcellAlias, deploymentinstanceGroupsToJobs)
	if err != nil {
		a.StderrLogger.Println(err.Error())
		return bosh.BoshManifest{}, errors.New("")
	}
	//zookeeper的instanceGroup
	zookeeperBrokerInstanceGroup := &instanceGroups[0]
	//zookeeper的网络
	if len(zookeeperBrokerInstanceGroup.Networks) != 1 {
		a.StderrLogger.Println(fmt.Sprintf("expected 1 network for %s, got %d", zookeeperBrokerInstanceGroup.Name, len(zookeeperBrokerInstanceGroup.Networks)))
		return bosh.BoshManifest{}, errors.New("")
	}
	//peers的属性配置
	//autopurge_purge_interval default 24
	autopurgePurgeInterval := 24
	if val, ok := servicePlan.Properties["autopurge_purge_interval"]; ok {
		autopurgePurgeInterval = int(val.(float64))
	}
	//autopurge_snap_retain_count
	autopurgeSnapRetainCount := 3
	if val, ok := servicePlan.Properties["autopurge_snap_retain_count"]; ok {
		autopurgeSnapRetainCount = int(val.(float64))
	}
	//client_port
	clientPort := 2181
	if val, ok := servicePlan.Properties["client_port"]; ok {
		clientPort = int(val.(float64))
	}
	//cnx_timeout
	cnxTimeout := 5
	if val, ok := servicePlan.Properties["cnx_timeout"]; ok {
		cnxTimeout = int(val.(float64))
	}
	//election_algorim
	electionAlgorim := 3
	if val, ok := servicePlan.Properties["election_algorim"]; ok {
		electionAlgorim = int(val.(float64))
	}
	//global_outstanding_limit
	globalOutstandingLimit := 1000
	if val, ok := servicePlan.Properties["global_outstanding_limit"]; ok {
		globalOutstandingLimit = int(val.(float64))
	}
	//init_limit
	initLimit := 5
	if val, ok := servicePlan.Properties["init_limit"]; ok {
		initLimit = int(val.(float64))
	}
	//leader_election_port
	leaderElectionPort := 3888
	if val, ok := servicePlan.Properties["leader_election_port"]; ok {
		autopurgeSnapRetainCount = int(val.(float64))
	}
	//max_client_connections
	maxClientConnections := 3
	if val, ok := servicePlan.Properties["max_client_connections"]; ok {
		maxClientConnections = int(val.(float64))
	}
	//max_session_timeout
	maxSessionTimeout := 40000
	if val, ok := servicePlan.Properties["max_session_timeout"]; ok {
		maxSessionTimeout = int(val.(float64))
	}
	//min_session_timeout
	minSessionTimeout := 4000
	if val, ok := servicePlan.Properties["min_session_timeout"]; ok {
		minSessionTimeout = int(val.(float64))
	}
	//pre_allocation_size string
	preAllocationSize := "65536"
	if val, ok := servicePlan.Properties["pre_allocation_size"]; ok {
		preAllocationSize = val.(string)
	}
	//snap_count
	snapCount := 100000
	if val, ok := servicePlan.Properties["snap_count"]; ok {
		snapCount = int(val.(float64))
	}
	//sync_limit
	syncLimit := 2
	if val, ok := servicePlan.Properties["sync_limit"]; ok {
		syncLimit = int(val.(float64))
	}
	//tick_time
	tickTime := 2000
	if val, ok := servicePlan.Properties["tick_time"]; ok {
		syncLimit = int(val.(float64))
	}
	//warning_threshold_ms
	warningThresholdMs := 1000
	if val, ok := servicePlan.Properties["warning_threshold_ms"]; ok {
		syncLimit = int(val.(float64))
	}
	//组合所有的属性到zookeeper的job里
	if zookeeperBrokerJob, ok := getJobFromInstanceGroup("zookeeper", zookeeperBrokerInstanceGroup); ok {
		zookeeperBrokerJob.Properties = map[string]interface{}{
			"autopurge_purge_interval": autopurgePurgeInterval,
			"autopurge_snap_retain_count": autopurgeSnapRetainCount,
			"client_port": clientPort,
			"cnx_timeout": cnxTimeout,
			"election_algorim": electionAlgorim,
			"global_outstanding_limit": globalOutstandingLimit,
			"init_limit": initLimit,
			"leader_election_port": leaderElectionPort,
			"max_client_connections": maxClientConnections,
			"max_session_timeout": maxSessionTimeout,
			"min_session_timeout": minSessionTimeout,
			"pre_allocation_size": preAllocationSize,
			"snap_count": snapCount,
			"sync_limit": syncLimit,
			"tick_time": tickTime,
			"warning_threshold_ms": warningThresholdMs,
			"network": zookeeperBrokerInstanceGroup.Networks[0].Name,
		}
	}

	//处理剩余的组件集成工作
	manifestProperties := map[string]interface{}{}
	//集成metron日志功能
	if includeMetron {
		logging := loggingRaw.(map[string]interface{})
		manifestProperties["syslog_daemon_config"] = map[interface{}]interface{}{
			"address": logging["syslog_address"],
			"port":    logging["syslog_port"],
		}
		manifestProperties["metron_agent"] = map[interface{}]interface{}{
			"zone":       "",
			"deployment": serviceDeployment.DeploymentName,
		}
		manifestProperties["loggregator"] = map[interface{}]interface{}{
			"etcd": map[interface{}]interface{}{
				"machines": logging["loggregator_etcd_addresses"].([]interface{}),
			},
		}
		manifestProperties["metron_endpoint"] = map[interface{}]interface{}{
			"shared_secret": logging["loggregator_shared_secret"],
		}
	}

	var updateBlock = bosh.Update{
		Canaries:        1,
		MaxInFlight:     10,
		CanaryWatchTime: "30000-240000",
		UpdateWatchTime: "30000-240000",
		Serial:          boolPointer(false),
	}

	if servicePlan.Update != nil {
		updateBlock = bosh.Update{
			Canaries:        servicePlan.Update.Canaries,
			MaxInFlight:     servicePlan.Update.MaxInFlight,
			CanaryWatchTime: servicePlan.Update.CanaryWatchTime,
			UpdateWatchTime: servicePlan.Update.UpdateWatchTime,
			Serial:          servicePlan.Update.Serial,
		}
	}

	return bosh.BoshManifest{
		Name:     serviceDeployment.DeploymentName,
		Releases: releases,
		Stemcells: []bosh.Stemcell{{
			Alias:   OnlystemcellAlias,
			OS:      serviceDeployment.Stemcell.OS,
			Version: serviceDeployment.Stemcell.Version,
		}},
		InstanceGroups: instanceGroups,
		Properties:     manifestProperties,
		Update:         updateBlock,
	}, nil

}

func boolPointer(b bool) *bool {
	return &b
}

//从实例组里获取Job信息
func getJobFromInstanceGroup(jobName string, instanceGroup *bosh.InstanceGroup) (*bosh.Job, bool) {
	for index, job:= range instanceGroup.Jobs {
		if job.Name == jobName {
			return &instanceGroup.Jobs[index], true
		}
	}
	return &bosh.Job{}, false
}

//检查实例组是否存在
func checkInstanceGroupsPresent(names []string, instanceGroups []serviceadapter.InstanceGroup) error{
	var missingNames []string
	for _, name := range names {
		if !containsInstanceGroup(name, instanceGroups) {
			missingNames = append(missingNames, name)
		}
	}
	if len(missingNames) > 0 {
		return fmt.Errorf("Invalid instance group configuration: expected to find: '%s' in list: '%s'",
			strings.Join(missingNames, ", "),
			strings.Join(getInstanceGroupNames(instanceGroups), ", "))
	}
	return nil;
}

//是否包含实例组名字
func containsInstanceGroup(name string, instanceGroups []serviceadapter.InstanceGroup) bool {
	for _, instanceGroup := range instanceGroups {
		if instanceGroup.Name == name {
			return true
		}
	}
	return false
}

//获取实例组的所有名字
func getInstanceGroupNames(instanceGroups []serviceadapter.InstanceGroup) []string{
	var instanceGroupNames []string
	for _, instanceGroup := range instanceGroups {
		instanceGroupNames = append(instanceGroupNames, instanceGroup.Name)
	}
	return instanceGroupNames
}

func instanceCounts(plan serviceadapter.Plan)  map[string]int{
	val := map[string]int{}
	for _, instanceGroup := range plan.InstanceGroups {
		val[instanceGroup.Name] = instanceGroup.Instances
	}
	return val
}