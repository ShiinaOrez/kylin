package crawler

import "context"

type Crawler interface {
	Run(ctx context.Context) error
	GetID() string
}