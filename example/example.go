package main

import (
	"context"
	"github.com/ShiinaOrez/kylin"
	_const "github.com/ShiinaOrez/kylin/const"
	"github.com/ShiinaOrez/kylin/crawler"
	"github.com/ShiinaOrez/kylin/interceptor"
	logger2 "github.com/ShiinaOrez/kylin/logger"
	"github.com/ShiinaOrez/kylin/param"
	"github.com/ShiinaOrez/kylin/render"
	"github.com/ShiinaOrez/kylin/result"
	"net/http"
	"os"

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

type ReturnData struct {}

func (d ReturnData) Format() string {
	return "result"
}

func main() {
	var (
		k            kylin.Kylin             = kylin.NewKylin()
		i            interceptor.Interceptor = &NameInterceptor{}
		imageCrawler crawler.Crawler         = &crawler.BaseCrawler{ID:"ArtStation-Crawler"}
	)

	err := k.AddInputInterceptor(&i, "tail")
	if err != nil {
		logger2.GetLogger(nil).Fatal(err.Error())
		return
	}
	imageCrawler.SetProc(artStationCrawler)

	err = k.RegisterCrawlerWithRender(&imageCrawler, render.FileRender{})
	if err != nil {
		logger2.GetLogger(nil).Fatal(err.Error())
		return
	}
	wd, _ := os.Getwd()
	p := param.NewJSONParam(`{"content": {"name": "fine", "path": "` + wd + `"}}`)

	ch := k.StartOn(p)
	defer k.Stop()

	select {
	case result := <-ch:
		switch result {
		case _const.Success:
			logger2.GetLogger(nil).Info("Success")
		case _const.Failed:
			logger2.GetLogger(nil).Info("Failed")
		}
	}
}

func artStationCrawler(ctx context.Context, notifyCh *chan int) result.Data {
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

	var d result.Data = ReturnData{}
	return d
}