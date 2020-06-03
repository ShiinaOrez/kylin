package tests

import (
	"context"
	"github.com/ShiinaOrez/kylin"
	"github.com/ShiinaOrez/kylin/crawler"
	"github.com/ShiinaOrez/kylin/interceptor"
	"github.com/ShiinaOrez/kylin/param"
)

type Interceptor struct {}

func(i *Interceptor) Run(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, "namespace", "Test")
	return ctx
}

func (i *Interceptor) GetID() string {
	return "namespace"
}

type ImageCrawler struct {
	crawler.BaseCrawler
}

func Test_Main() {
	k := kylin.NewKylin()
	var _ interceptor.Interceptor = &Interceptor{}
	var namespaceInterceptor interceptor.Interceptor = &Interceptor{}
	k.RegisterInputInterceptor(&namespaceInterceptor)

	var _ crawler.Crawler = &ImageCrawler{}
	var imageCrawler crawler.Crawler = &ImageCrawler{}
	imageCrawler.SetID("image-crawler")
	imageCrawler.SetProc(func(ctx context.Context) {
		// main process of crawler
	})
	k.RegisterCrawler(&imageCrawler)

	p := param.NewJSONParam(`{"name": "Computer Network"}`)
	ch := k.StartOn(p)
	defer k.Stop()

	select {
	case result := <-ch:
		switch result {
		case kylin.Success:
			k.GetLogger().Info("Success")
		case kylin.Failed:
			k.GetLogger().Info("Failed")
		}
	}
}