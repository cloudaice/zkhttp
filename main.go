package main

import (
	"github.com/cloudaice/zkhttp/api"
	"github.com/cloudaice/zkhttp/zkop"
	"log"
	"net/url"
)

func router(api *api.API) {
	api.AddResource(new(zookeeper), "/zookeeper")
}

type zookeeper struct {
	api.DefaultResource
}

/*
   {
       "action": "stat",
       "cluster": "",
       "cluster": "",
       "cluster": "",

*/
func (zk zookeeper) Get(values url.Values) (int, interface{}) {
	action, ok := values["action"]
	if !ok {
		return 200, "no action"
	}

	cluster, ok := values["cluster"]
	if !ok {
		return 200, "no cluster"
	}

	zh := zkop.NewZh(cluster)
	defer zh.Close()

	switch action[0] {
	case "stat":
		stats := zh.State()
		return 200, stats
	case "ls":
		znode, ok := values["znode"]
		if !ok {
			return 200, "no znode"
		}

		nodes, err := zh.Ls(znode[0])
		if err != nil {
			log.Println(err)
		}

		return 200, nodes
	case "get":

		znode, ok := values["znode"]
		if !ok {
			return 200, "no znode"
		}

		data, err := zh.Get(znode[0])
		if err != nil {
			log.Println(err)
		}

		return 200, data
	}
	return 200, "hello world"
}

func main() {
	apiServer := api.NewAPI()
	router(apiServer)
	apiServer.Start(8080)
}
