package function

import (
	"os"
	"strings"

	faasflow "github.com/s8sg/faas-flow"

	redisDataStore "github.com/chennqqi/faas-flow-redis-datastore"
	redisStateStore "github.com/chennqqi/faas-flow-redis-statestore"
	consulStateStore "github.com/s8sg/faas-flow-consul-statestore"
	etcdStateStore "github.com/s8sg/faas-flow-etcd-statestore"
	minioDataStore "github.com/s8sg/faas-flow-minio-datastore"
)

// Define provide definiton of the workflow
func Define(flow *faasflow.Workflow, context *faasflow.Context) (err error) {
	flow.SyncNode().Apply("func-ofwatchdog").Apply("func-ofwatchdog")
	//flow.Node("n2").Apply("func-ofwatchdog")
	//flow.Edge("n1", "n2")
	return nil
}

// DefineStateStore provides the override of the default StateStore
func DefineStateStore() (faasflow.StateStore, error) {
	state := os.Getenv("statestore")
	switch strings.ToLower(state) {
	case "etcd":
		return etcdStateStore.GetEtcdStateStore(os.Getenv("etcd_url"))
	case "redis":
		return redisStateStore.GetRedisStateStore(os.Getenv("redis_url"), os.Getenv("redis_master"))
	default:
		fallthrough
	case "consul":
		return consulStateStore.GetConsulStateStore(os.Getenv("consul_url"), os.Getenv("consul_dc"))
	}
}

// ProvideDataStore provides the override of the default DataStore
func DefineDataStore() (faasflow.DataStore, error) {
	state := os.Getenv("datastore")
	switch strings.ToLower(state) {
	case "redis":
		return redisDataStore.InitFromEnv()
	default:
		fallthrough
	case "minio":
		return minioDataStore.InitFromEnv()
	}
}
