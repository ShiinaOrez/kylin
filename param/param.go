package param

import (
	"context"
	"encoding/json"
)

type Param interface {
	Resolve() context.Context
}

type JSONParam struct {
	content     string
}

type jsonContent struct {
	Content     map[string]string    `json:"content"`
}

func(jp JSONParam) Resolve() context.Context {
	ctx := context.Background()
	c := jsonContent{}
	json.Unmarshal([]byte(jp.content), &c)
	for k, v := range c.Content {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

func NewJSONParam(content string) JSONParam {
	return JSONParam{content:content}
}