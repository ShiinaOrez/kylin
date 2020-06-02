package kylin

import (
	"github.com/ShiinaOrez/kylin/crawler"
	"github.com/ShiinaOrez/kylin/interceptor"
	"github.com/ShiinaOrez/kylin/logger"
	"github.com/ShiinaOrez/kylin/manager"
)

type Kylin struct {
	manager            *manager.Manager
	inputInterceptors  []*interceptor.Interceptor
	outputInterceptors []*interceptor.Interceptor
	logger             logger.Logger
}

type KylinConfig struct {
	_                  interface{}
}

func NewKylin() Kylin {
	return Kylin{}
}

func NewKylinByConfig(conf KylinConfig) Kylin {
	return Kylin{}
}

func (kylin *Kylin) RegisterInputInterceptor(i *interceptor.Interceptor) error {
	kylin.inputInterceptors = append(kylin.inputInterceptors, i)
	return nil
}

func (kylin *Kylin) RegisterOutputInterceptor(i *interceptor.Interceptor) error {
	kylin.outputInterceptors = append(kylin.outputInterceptors, i)
	return nil
}

func (kylin *Kylin) SetLogger(l logger.Logger) error {
	kylin.logger = l
	return nil
}

func (kylin *Kylin) RegisterCrawler(c *crawler.Crawler) error {
	return kylin.manager.AddCrawler(c)
}