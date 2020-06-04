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

type NameInterceptor struct {}

func(i *NameInterceptor) Run(ctx context.Context) context.Context {
	if name := ctx.Value("name").(string); name == "" {
		ctx = context.WithValue(ctx, "break", i.GetID())
	}
	return ctx
}

func (i *NameInterceptor) GetID() string {
	return "name-interceptor"
}

type ImageCrawler struct {
	crawler.BaseCrawler
}

func (ic ImageCrawler) GetID() string {
	return "ArtStation-Crawler"
}

func main() {
	var (
		k            kylin.Kylin             = kylin.NewKylin()
		i            interceptor.Interceptor = &NameInterceptor{}
		imageCrawler crawler.Crawler         = &ImageCrawler{}
	)

	k.RegisterInputInterceptor(&i)
	imageCrawler.SetProc(artStationCrawler)

	err := k.RegisterCrawler(&imageCrawler)
	if err != nil {
		k.GetLogger().Fatal(err.Error())
		return
	}

	p := param.NewJSONParam(`{"content": {"name": "timbougami"}}`)

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

func artStationCrawler(ctx context.Context, notifyCh *chan int) {
	c := new(http.Client)
	req := request.NewRequest(c)

	artistName := ctx.Value("name")
	resp, err := req.Get("https://www.artstation.com/"+artistName.(string))
	defer resp.Body.Close()

	if err == nil {
		*notifyCh<- _const.StatusSuccess
	} else {
		*notifyCh<- _const.StatusFailed
	}
	return
}