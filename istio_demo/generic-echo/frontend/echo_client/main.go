package main

import (
	"client/conf"
	"client/content"
	"flag"
	"log"
	"os"
)

var config conf.Conf
var parserErr error
var protocol conf.Protocol

func init() {
	var confFile string
	flag.StringVar(&confFile, "c", os.Args[1], "config file")
	flag.Parse()

	config, parserErr = conf.ConfParser(confFile, &protocol)
	if parserErr != nil {
		log.Fatalf("parser config failed:", parserErr.Error())
	}
}

func main() {
	IFactory := &content.ContentFactory{}
	IContent := IFactory.Create(config)
	if IContent == nil {
		panic("content is nil.")
	}
	IContent.Start()
}
