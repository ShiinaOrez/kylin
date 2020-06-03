package manager

import (
	"context"
	"errors"
	"github.com/ShiinaOrez/kylin"
	"github.com/ShiinaOrez/kylin/crawler"
	"github.com/ShiinaOrez/kylin/interceptor"
	"sync"
)

type Manager struct {
	inputInterceptors  []*interceptor.Interceptor
	crawlers           map[string]*crawler.Crawler
	// outputInterceptors []*interceptor.Interceptor
}

func (manager Manager) AddCrawler(c *crawler.Crawler) error {
	id := (*c).GetID()
	if _, ok := manager.crawlers[id]; !ok {
		manager.crawlers[id] = c
	} else {
		return errors.New("Can't register crawler which ID duplicated ")
	}
	return nil
}

func (manager Manager) Dispatch(ctx context.Context, resultCh chan kylin.Result) {
	for _, i := range manager.inputInterceptors {
		ctx = (*i).Run(ctx)
	}
	wg := &sync.WaitGroup{}
	for _, c := range manager.crawlers {
		wg.Add(1)
		go (*c).Run(ctx, wg)
	}
	wg.Done()
	resultCh<- kylin.Success
	return
}