package manager

import (
	"context"
	"errors"
	_const "github.com/ShiinaOrez/kylin/const"
	"github.com/ShiinaOrez/kylin/crawler"
	"github.com/ShiinaOrez/kylin/interceptor"
	"github.com/ShiinaOrez/kylin/result"
	"sync"
)

type Manager struct {
	inputInterceptors  []*interceptor.Interceptor
	crawlers           map[string]*crawler.Crawler
	// outputInterceptors []*interceptor.Interceptor
	resultHandler      func(map[string]*chan int) (result.Result, error)
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

func (manager Manager) Dispatch(ctx context.Context, resultCh chan result.Result) error {
	for _, i := range manager.inputInterceptors {
		ctx = (*i).Run(ctx)
	}
	notifyChMap := make(map[string]*chan int)
	defer func() {
		for _, ch := range notifyChMap {
			close(*ch)
		}
	}()

	wg := sync.WaitGroup{}
	for _, c := range manager.crawlers {
		notifyCh := make(chan int, 1)
		notifyChMap[(*c).GetID()] = &notifyCh
		wg.Add(1)
		go func() {
			(*c).Run(ctx, &notifyCh, &wg)
		}()
	}
	wg.Wait()
	result, err := manager.resultHandler(notifyChMap)

	if err != nil {
		return err
	}
	resultCh<- result
	return nil
}

func DefaultResultHandler(notifyChMap map[string]*chan int) (result.Result, error) {
	r := _const.Success
	for _, ch := range notifyChMap {
		result := <-*ch
		if result == _const.StatusFailed {
			r = _const.Failed
		}
	}
	return r, nil
}

func NewManager() Manager {
	manager := Manager{}
	manager.crawlers = make(map[string]*crawler.Crawler)
	manager.resultHandler = DefaultResultHandler
	return manager
}