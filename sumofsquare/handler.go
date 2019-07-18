package function

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	faasflow "github.com/s8sg/faas-flow"
	consulStateStore "github.com/s8sg/faas-flow-consul-statestore"
	minioDataStore "github.com/s8sg/faas-flow-minio-datastore"
)

// Define provide definiton of the workflow
func Define(flow *faasflow.Workflow, context *faasflow.Context) (err error) {
	dag := flow.Dag()
	foreachDag := dag.ForEachBranch("square-each", func(data []byte) map[string][]byte {
		log.Println("square-each:", string(data))
		r := strings.NewReader("square-each:" + string(data))
		http.Post("http://10.143.143.19:5888", "plain/text", r)
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
		r := strings.NewReader("agg:" + buf.String())
		http.Post("http://10.143.143.19:5888", "plain/text", r)
		return buf.Bytes(), nil
	}))
	foreachDag.Node("square").Apply("square")

	dag.Node("add").Apply("add").Modify(func(data []byte) ([]byte, error) {
		r := strings.NewReader("add:" + string(data))
		http.Post("http://10.143.143.19:5888", "plain/text", r)
		return data, nil
	})
	dag.Edge("square-each", "add")
	return nil
}

// DefineStateStore provides the override of the default StateStore
func DefineStateStore() (faasflow.StateStore, error) {
	return consulStateStore.GetConsulStateStore(os.Getenv("consul_url"), os.Getenv("consul_dc"))
}

// ProvideDataStore provides the override of the default DataStore
func DefineDataStore() (faasflow.DataStore, error) {
	return minioDataStore.InitFromEnv()
}
