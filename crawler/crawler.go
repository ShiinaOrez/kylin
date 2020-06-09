package crawler

import (
	"context"
	"github.com/ShiinaOrez/kylin/interceptor"
	"github.com/ShiinaOrez/kylin/result"
	"sync"
)

type Crawler interface {
	Run(ctx context.Context, notifyCh *chan int, group *sync.WaitGroup, dataMap *result.DataMap)
	GetID() string
	SetProc(func(ctx context.Context, notifyCh *chan int) result.Data)
}

type BaseCrawler struct {
	ID                 string
	inputInterceptors  []*interceptor.Interceptor
	proc               func(ctx context.Context, notifyCh *chan int) result.Data
	outputInterceptors []*interceptor.Interceptor

}

func (wc *BaseCrawler) GetID() string {
	return wc.ID
}

func (wc *BaseCrawler) SetProc(newProc func(ctx context.Context, notify *chan int) result.Data) {
	wc.proc = newProc
	return
}

func (wc *BaseCrawler) Run(ctx context.Context, notifyCh *chan int, group *sync.WaitGroup, dataMap *result.DataMap) {
	data := wc.proc(ctx, notifyCh)
	dataMap.Lock.Lock()
	dataMap.Map[wc.GetID()] = data
	defer dataMap.Lock.Unlock()

	(*group).Done()
	return
}