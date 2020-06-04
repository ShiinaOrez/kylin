package param

import (
	"context"
	"encoding/json"
)

type Param interface {
	Resolve() context.Context
}

type EmptyParam struct {}

type JSONParam struct {
	content     string
}

type jsonContent struct {
	Content     map[string]string    `json:"content"`
}

func (ep EmptyParam) Resolve() context.Context {
	return context.WithValue(context.Background(), "break", "")
}

func(jp JSONParam) Resolve() context.Context {
	ctx := context.WithValue(context.Background(), "break", "")
	c := jsonContent{}
	json.Unmarshal([]byte(jp.content), &c)
	for k, v := range c.Content {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

func NewEmptyParam() EmptyParam {
	return EmptyParam{}
}

func NewJSONParam(content string) JSONParam {
	return JSONParam{content:content}
}
