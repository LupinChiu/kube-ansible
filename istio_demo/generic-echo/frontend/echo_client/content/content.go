package content

import (
	"client/conf"
)

type IContent interface {
	Start()
}

type IContentFactory interface {
	Create(config conf.Conf) IContent
}

type ContentFactory struct {
}

func (cf ContentFactory) Create(config conf.Conf) IContent {
	var content IContent
	switch config.Name {
	case "tcp":
	case "http":
		content = &Http{
			config: config,
		}
	case "https":
	case "grpc":
	}
	return content
}
