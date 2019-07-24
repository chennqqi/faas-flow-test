package function

import (
	"bytes"
	"fmt"
	"log"
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
	dag := flow.Dag()
	foreachDag := dag.ForEachBranch("square-each", func(data []byte) map[string][]byte {
		log.Println("square-each:", string(data))
		values := bytes.SplitN(data, []byte(","), -1)
		rmap := make(map[string][]byte)
		for i := 0; i < len(values); i++ {
			rmap[fmt.Sprintf("%d", i)] = values[i]
		}
		return rmap
	}, faasflow.Aggregator(func(vmap map[string][]byte) ([]byte, error) {
		var buf bytes.Buffer
		var nonestart bool
		for _, v := range vmap {
			if nonestart {
				buf.WriteByte(',')
			} else {
				nonestart = true
			}
			buf.Write(v)
		}
		log.Println("agg:" + buf.String())
		return buf.Bytes(), nil
	}))
	foreachDag.Node("square").Apply("square")

	dag.Node("add").Apply("add").Modify(func(data []byte) ([]byte, error) {
		log.Println("add:" + string(data))
		return data, nil
	})
	dag.Edge("square-each", "add")
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
