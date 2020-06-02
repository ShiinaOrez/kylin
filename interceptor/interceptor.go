package interceptor

import "context"

type Interceptor interface {
	Name() string
	Run(ctx context.Context) context.Context
}