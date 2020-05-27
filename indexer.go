package main

import (
	"github.com/chadwcarlson/gomeili/index"
	"github.com/chadwcarlson/gomeili/utils/config"
)

func main() {

	var config config.ConfigSet
	meilindex := index.Build(config.Load("config.yaml"))
	meilindex.Write("output/index.json")

	// meilindex.Post()

}
