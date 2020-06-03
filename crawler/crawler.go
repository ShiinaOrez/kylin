package crawler

import (
	"context"
	"github.com/ShiinaOrez/kylin/interceptor"
	"sync"
)

type Crawler interface {
	Run(ctx context.Context, notifyCh *chan int, group *sync.WaitGroup)
	GetID() string
	SetProc(func(ctx context.Context, notifyCh *chan int))
}

type BaseCrawler struct {
	inputInterceptors  []*interceptor.Interceptor
	proc               func(ctx context.Context, notifyCh *chan int)
	outputInterceptors []*interceptor.Interceptor

}

func (wc *BaseCrawler) GetID() string {
	return "Base-Crawler"
}

func (wc *BaseCrawler) SetProc(newProc func(ctx context.Context, notify *chan int)) {
	wc.proc = newProc
	return
}

func (wc *BaseCrawler) Run(ctx context.Context, notifyCh *chan int, group *sync.WaitGroup) {
	wc.proc(ctx, notifyCh)
	(*group).Done()
	return
}