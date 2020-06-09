package manager

import (
	"context"
	"errors"
	"fmt"
	_const "github.com/ShiinaOrez/kylin/const"
	"github.com/ShiinaOrez/kylin/crawler"
	"github.com/ShiinaOrez/kylin/interceptor"
	logger2 "github.com/ShiinaOrez/kylin/logger"
	"github.com/ShiinaOrez/kylin/render"
	"github.com/ShiinaOrez/kylin/result"
	"sync"
)

type Manager struct {
	inputInterceptors  []*interceptor.Interceptor
	crawlers           map[string]*crawler.Crawler
	renderMap          map[string]render.Render
	// outputInterceptors []*interceptor.Interceptor
	resultHandler      func(map[string]*chan int) (result.Result, error)

	once               sync.Once
}

func (manager *Manager) AddInputInterceptor(i *interceptor.Interceptor, mode string) error {
	if mode == "tail" {
		manager.inputInterceptors = append(manager.inputInterceptors, i)
	} else if mode == "head" {
		manager.inputInterceptors = append([]*interceptor.Interceptor{i}, manager.inputInterceptors...)
	} else {
		return errors.New("Add input interceptor mode invalid.")
	}
	return nil
}

func (manager *Manager) AddCrawler(c *crawler.Crawler) error {
	id := (*c).GetID()
	if _, ok := manager.crawlers[id]; !ok {
		manager.crawlers[id] = c
	} else {
		return errors.New("Can't register crawler which ID duplicated ")
	}
	return nil
}

func (manager *Manager) AddRender(id string, r render.Render) error {
	if id == "" {
		return errors.New("Crawler ID not be EMPTY string.")
	}
	if _, ok := manager.crawlers[id]; !ok {
		return errors.New(fmt.Sprintf("Crawler %s not register.", id))
	} else {
		manager.renderMap[id] = r
	}
	return nil
}

func (manager Manager) Dispatch(ctx context.Context, resultCh chan result.Result) (error) {
	for _, i := range manager.inputInterceptors {
		ctx = (*i).Run(ctx)
		if interceptorID := ctx.Value("break").(string); interceptorID != "" {
			logger2.GetLogger(ctx).Warning("Interceptor trigger break! InterceptorID: "+interceptorID)
			return errors.New("Interceptor trigger break ")
		}
	}
	notifyChMap := make(map[string]*chan int)
	dataMap := result.DataMap{
		Lock: new(sync.Mutex),
		Map:  make(map[string]result.Data),
	}
	defer func() {
		for _, ch := range notifyChMap {
			close(*ch)
		}
	}()

	wg := sync.WaitGroup{}
	for id, c := range manager.crawlers {
		notifyCh := make(chan int, 1)
		notifyChMap[(*c).GetID()] = &notifyCh
		wg.Add(1)
		logger2.GetLogger(ctx).Info("Crawler ID: "+id+" start running...")
		go func() {
			(*c).Run(ctx, &notifyCh, &wg, &dataMap)
		}()
	}
	wg.Wait()
	logger2.GetLogger(ctx).Info("All crawler running over.")
	result, err := manager.resultHandler(notifyChMap)

	for id, data := range dataMap.Map {
		nCtx := context.WithValue(ctx, "id", id)
		if err := manager.renderMap[id].Do(nCtx, data); err != nil {
			logger2.GetLogger(ctx).Fatal(fmt.Sprintf("Render crawler %s results error, reason: %s", id, err.Error()))
		}
	}

	if err != nil {
		return err
	}
	resultCh<- result
	return nil
}

func (manager *Manager) SetResultHandler(f func(notifyChMap map[string]*chan int) (result.Result, error)) {
	manager.resultHandler = f
	return
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
	manager.renderMap = make(map[string]render.Render)
	manager.SetResultHandler(DefaultResultHandler)
	logger2.GetLogger(nil).Info("Kylin manager set up ResultHandler: DefaultResultHandler.")
	return manager
}