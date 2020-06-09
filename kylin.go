package kylin

import (
	"context"
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

	resultCh           chan result.Result
	once               sync.Once
}

type KylinConfig struct {
	_                  interface{}
}

func NewKylin() Kylin {
	kylin := Kylin{}
	logger.GetLogger(nil).Info("Kylin set up logger: DefaultLogger.")

	m := manager.NewManager()
	kylin.manager = &m
	logger.GetLogger(nil).Info("Kylin set up manager.")

	kylin.resultCh = make(chan result.Result, 1)
	return kylin
}

func NewKylinByConfig(conf KylinConfig) Kylin {
	return NewKylin()
}

func (kylin *Kylin) RegisterCrawler(c *crawler.Crawler) error {
	return kylin.manager.AddCrawler(c)
}

func (kylin *Kylin) RegisterCrawlerWithRender(c *crawler.Crawler, r render2.Render) error {
	err := kylin.manager.AddCrawler(c)
	if err != nil {
		return err
	}
	err = kylin.manager.AddRender((*c).GetID(), r)
	if err != nil {
		return err
	}
	return nil
}

func (kylin *Kylin) AddInputInterceptor(i *interceptor.Interceptor, mode string) error {
	return kylin.manager.AddInputInterceptor(i, mode)
}

func (kylin *Kylin) StartOn(p param.Param) <-chan result.Result {
	ctx := p.Resolve()
	ctx = context.WithValue(ctx, "kylin-logger", logger.DefaultLogger{})
	logger.GetLogger(ctx).Info("Kylin start running...")

	err := kylin.manager.Dispatch(ctx, kylin.resultCh)
	if err != nil {
		logger.GetLogger(ctx).Warning("Call manager Dispatch method error, reason: "+err.Error())
		kylin.resultCh<- _const.Failed
	}
	if err != nil {
		logger.GetLogger(ctx).Fatal(err.Error())
	}
	return kylin.resultCh
}

func (kylin *Kylin) Stop() {
	kylin.safeClose()
	logger.GetLogger(nil).Info("Kylin running over")
}

func (kylin *Kylin) safeClose() {
	kylin.once.Do(func() {
		close(kylin.resultCh)
	})
}