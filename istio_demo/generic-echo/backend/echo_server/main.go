package main

import (
	"flag"
	"log"
	"os"
	database "server/db"

	conf "server/conf"
	"server/content"
)

var config conf.Conf
var parserErr error
var configmap conf.ConfigMap

func init() {
	var confFile string
	flag.StringVar(&confFile, "c", os.Args[1], "config file")
	flag.Parse()

	config, parserErr = conf.ConfParser(confFile, &configmap)
	if parserErr != nil {
		log.Fatalf("parser config failed:", parserErr.Error())
	}

	database.InitDbConn(config.Db)
}

func main() {
	IContent := content.Factory.Create(config)
	if IContent == nil {
		panic("content is nil.")
	}
	IContent.Start()
}
