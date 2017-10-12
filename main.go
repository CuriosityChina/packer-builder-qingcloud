package main

import (
	"github.com/CuriosityChina/packer-builder-qingcloud/qingcloud"
	"github.com/hashicorp/packer/packer/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}

	server.RegisterBuilder(new(qingcloud.Builder))
	server.Serve()
}
