package kylin

import (
	_const "github.com/ShiinaOrez/kylin/const"
	"github.com/ShiinaOrez/kylin/crawler"
	"github.com/ShiinaOrez/kylin/interceptor"
	"github.com/ShiinaOrez/kylin/logger"
	"github.com/ShiinaOrez/kylin/manager"
	"github.com/ShiinaOrez/kylin/param"
	render2 "github.com/ShiinaOrez/kylin/render"
	"github.com/ShiinaOrez/kylin/result"
	"sync"
)

type Kylin struct {
	manager            *manager.Manager
	logger             logger.Logger

	resultCh           chan result.Result
	once               sync.Once
}

type KylinConfig struct {
	_                  interface{}
}

func NewKylin() Kylin {
	kylin := Kylin{}
	kylin.SetLogger(logger.DefaultLogger{})
	kylin.GetLogger().Info("Kylin set up logger: DefaultLogger.")

	m := manager.NewManager()
	kylin.manager = &m
	kylin.manager.SetLogger(kylin.logger)
	kylin.GetLogger().Info("Kylin set up manager.")

	kylin.resultCh = make(chan result.Result, 1)
	return kylin
}

func NewKylinByConfig(conf KylinConfig) Kylin {
	return NewKylin()
}

func (kylin *Kylin) SetLogger(l logger.Logger) error {
	kylin.logger = l
	return nil
}

func (kylin *Kylin) GetLogger() logger.Logger {
	if kylin.logger == nil {
		kylin.once.Do(func() {
			kylin.SetLogger(logger.DefaultLogger{})
		})
	}
	return kylin.logger
}

func (kylin *Kylin) RegisterCrawler(c *crawler.Crawler) error {
	return kylin.manager.AddCrawler(c)
}

func (kylin *Kylin) AddInputInterceptor(i *interceptor.Interceptor, mode string) error {
	return kylin.manager.AddInputInterceptor(i, mode)
}

func (kylin *Kylin) StartOn(p param.Param) <-chan result.Result {
	kylin.GetLogger().Info("Kylin start running...")
	ctx := p.Resolve()
	dataMap, err := kylin.manager.Dispatch(ctx, kylin.resultCh)
	if err != nil {
		kylin.GetLogger().Warning("Call manager Dispatch method error, reason: "+err.Error())
		kylin.resultCh<- _const.Failed
	}
	render := render2.NewRender(dataMap)
	err = render.Do(p, render2.SaveAsFile)
	if err != nil {
		kylin.GetLogger().Fatal(err.Error())
	}
	return kylin.resultCh
}

func (kylin *Kylin) Stop() {
	kylin.safeClose()
	kylin.GetLogger().Info("Kylin running over")
}

func (kylin *Kylin) safeClose() {
	kylin.once.Do(func() {
		close(kylin.resultCh)
	})
}