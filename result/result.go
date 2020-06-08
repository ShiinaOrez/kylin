package result

import "sync"

type Result struct {
	Status      int
	Description string
}

type Data interface {
	Format() string
}

type DataMap struct {
	Lock      *sync.Mutex
	Map       map[string]Data
}