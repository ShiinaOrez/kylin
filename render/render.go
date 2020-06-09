package render

import (
	"context"
	"github.com/ShiinaOrez/kylin/result"
)

type Render interface {
	Do(ctx context.Context, data result.Data) error
}

type FileRender struct {}

func (r FileRender) Do(ctx context.Context, data result.Data) error {
	SaveAsFile(ctx, data.Format())
	return nil
}