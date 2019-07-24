# project for testing faas-flow

	testing examples for [faas-flow](https://github.com/s8sg/faas-flow)

## install 

1. datastore
	redis/minio
2. statestore
	redis/consul/etcd

## deploy

```bash
	cd flow-ofwatchdog 
	dep init
	cd ..
	cd sumofsquare
	dep init
	faas-cli up
```