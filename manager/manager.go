package manager

import (
	"errors"
	"github.com/ShiinaOrez/kylin/crawler"
)

type Manager struct {
	crawlers map[string]*crawler.Crawler
}

func (manager Manager) AddCrawler(c *crawler.Crawler) error {
	id := (*c).GetID()
	if _, ok := manager.crawlers[id]; !ok {
		manager.crawlers[id] = c
	} else {
		return errors.New("Can't register crawler which ID duplicated.")
	}
	return nil
}