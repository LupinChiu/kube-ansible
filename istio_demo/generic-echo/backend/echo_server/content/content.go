package content

import (
	"server/conf"
)

type IContent interface {
	Start()
	GetName() string
	SetConf(config conf.Conf)
}

type IContentFactory interface {
	Create(config conf.Conf) IContent
	Add(IContent)
}

type ContentFactory struct {
	Container map[string]IContent
}

func (cf ContentFactory) Create(config conf.Conf) IContent {
	var content IContent
	var found bool

	if content, found = cf.Container[config.Proto.Name]; !found {
		panic("content not implement.")
	}

	content.SetConf(config)
	return content
}

func (cf ContentFactory) Add(content IContent) {
	if _, found := cf.Container[content.GetName()]; found {
		return
	}
	cf.Container[content.GetName()] = content
}

var factory IContentFactory

func GetFactory() IContentFactory {
	if factory == nil {
		factory = &ContentFactory{
			Container: make(map[string]IContent),
		}
	}
	return factory
}
