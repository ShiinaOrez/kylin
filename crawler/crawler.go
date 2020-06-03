package crawler

import (
	"context"
	"github.com/ShiinaOrez/kylin/interceptor"
	"sync"
)

type Crawler interface {
	Run(ctx context.Context, group *sync.WaitGroup) error
	GetID() string
	SetID(string)
	SetProc(func(ctx context.Context))
}

type BaseCrawler struct {
	id                 string
	inputInterceptors  []*interceptor.Interceptor
	proc               func(ctx context.Context)
	outputInterceptors []*interceptor.Interceptor

}

func (wc BaseCrawler) GetID() string {
	return wc.id
}

func (wc BaseCrawler) SetID(newID string) {
	wc.id = newID
	return
}

func (wc BaseCrawler) SetProc(newProc func(ctx context.Context)) {
	wc.proc = newProc
	return
}

func (wc BaseCrawler) Run(ctx context.Context, group *sync.WaitGroup) error {
	wc.proc(ctx)
	(*group).Done()
	return nil
}