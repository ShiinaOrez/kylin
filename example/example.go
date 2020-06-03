package main

import (
	"context"
	"github.com/ShiinaOrez/kylin"
	_const "github.com/ShiinaOrez/kylin/const"
	"github.com/ShiinaOrez/kylin/crawler"
	"github.com/ShiinaOrez/kylin/interceptor"
	"github.com/ShiinaOrez/kylin/param"
	"net/http"

	"github.com/mozillazg/request"
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

func (ic ImageCrawler) GetID() string {
	return "Image-Crawler"
}

func main() {
	k := kylin.NewKylin()
	var _ interceptor.Interceptor = &Interceptor{}
	var namespaceInterceptor interceptor.Interceptor = &Interceptor{}
	k.RegisterInputInterceptor(&namespaceInterceptor)

	var _ crawler.Crawler = &ImageCrawler{}
	var imageCrawler crawler.Crawler = &ImageCrawler{}
	imageCrawler.SetProc(func(ctx context.Context, notifyCh *chan int) {
		c := new(http.Client)
		req := request.NewRequest(c)
		resp, err := req.Get("https://github.com")
		defer resp.Body.Close()

		if err == nil {
			*notifyCh<- _const.StatusSuccess
		} else {
			*notifyCh<- _const.StatusFailed
		}
		return
	})
	err := k.RegisterCrawler(&imageCrawler)
	if err != nil {
		k.GetLogger().Fatal(err.Error())
		return
	}

	p := param.NewJSONParam(`{"name": "Computer Network"}`)
	ch := k.StartOn(p)
	defer k.Stop()

	select {
	case result := <-ch:
		switch result {
		case _const.Success:
			k.GetLogger().Info("Success")
		case _const.Failed:
			k.GetLogger().Info("Failed")
		}
	}
}