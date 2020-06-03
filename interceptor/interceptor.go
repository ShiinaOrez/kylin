package interceptor

import "context"

type Interceptor interface {
	GetID() string
	Run(ctx context.Context) context.Context
}