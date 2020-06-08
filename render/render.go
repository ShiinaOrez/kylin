package render

import (
	"context"
	logger2 "github.com/ShiinaOrez/kylin/logger"
	"github.com/ShiinaOrez/kylin/param"
	"github.com/ShiinaOrez/kylin/result"
)

type Render struct {
	dataMap     map[string]result.Data
	logger      logger2.Logger
}

func NewRender(dataMap map[string]result.Data) Render {
	return Render{dataMap:dataMap}
}

func (r *Render) Do(p param.Param, way func(ctx context.Context, content string) error ) error {
	ctx := p.Resolve()
	for k, v := range r.dataMap {
		ctx = context.WithValue(ctx, "id", k)
		err := way(ctx, v.Format())
		if err != nil {
			return err
		}
	}
	return nil
}