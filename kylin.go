package kylin

import (
	"errors"
	"github.com/ShiinaOrez/kylin/crawler"
	"github.com/ShiinaOrez/kylin/interceptor"
	"github.com/ShiinaOrez/kylin/logger"
	"github.com/ShiinaOrez/kylin/manager"
	"github.com/ShiinaOrez/kylin/param"
	"sync"
)

type Kylin struct {
	manager            *manager.Manager
	inputInterceptors  map[string]*interceptor.Interceptor
	outputInterceptors map[string]*interceptor.Interceptor
	logger             logger.Logger

	resultCh           chan Result
	once               sync.Once
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
	id := (*i).GetID()
	if _, ok := kylin.inputInterceptors[id]; !ok {
		kylin.inputInterceptors[id] = i
	} else {
		return errors.New("Can't register input interceptor with duplicate ID ")
	}
	return nil
}

func (kylin *Kylin) RegisterOutputInterceptor(i *interceptor.Interceptor) error {
	id := (*i).GetID()
	if _, ok := kylin.outputInterceptors[id]; !ok {
		kylin.outputInterceptors[id] = i
	} else {
		return errors.New("Can't register output interceptor with duplicate ID ")
	}
	return nil
}

func (kylin *Kylin) SetLogger(l logger.Logger) error {
	kylin.logger = l
	return nil
}

func (kylin *Kylin) GetLogger() logger.Logger {
	kylin.once.Do(func() {
		kylin.SetLogger(logger.DefaultLogger{})
	})
	return kylin.logger
}

func (kylin *Kylin) RegisterCrawler(c *crawler.Crawler) error {
	return kylin.manager.AddCrawler(c)
}

func (kylin *Kylin) StartOn(p param.Param) <-chan Result {
	ctx := p.Resolve()
	kylin.manager.Dispatch(ctx, kylin.resultCh)
	return kylin.resultCh
}

func (kylin *Kylin) Stop() {
	kylin.safeClose()
}

func (kylin *Kylin) safeClose() {
	kylin.once.Do(func() {
		close(kylin.resultCh)
	})
}