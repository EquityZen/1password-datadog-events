package cache

import (
	"errors"
	"sync"
	"time"
)

type item struct {
	Id       int64
	itemUUID string
}

type cachedItem struct {
	item
	expiredTimestamp int64
}

type cache struct {
	stop  chan struct{}
	wg    sync.WaitGroup
	mu    sync.Mutex
	items map[int64]cachedItem
}

func NewCacheStore(cleanupInterval time.Duration) *cache {
	c := &cache{
		stop:  make(chan struct{}),
		items: make(map[int64]cachedItem),
	}

	c.wg.Add(1)
	go func(cleanupInterval time.Duration) {
		defer c.wg.Done()
		c.cleanup(cleanupInterval)
	}(cleanupInterval)
	return c
}

func (c *cache) cleanup(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-t.C:
			c.mu.Lock()
			for uid, cu := range c.items {
				if cu.expiredTimestamp <= time.Now().Unix() {
					delete(c.items, uid)
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *cache) stopClean() {
	close(c.stop)
	c.wg.Wait()
}

func (c *cache) update(i item, expAt int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[i.Id] = cachedItem{
		item:             i,
		expiredTimestamp: expAt,
	}
}

var (
	errItemNotInCache = errors.New("the user isn't in cache")
)

func (c *cache) read(id int64) (item, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ci, ok := c.items[id]
	if ok {
		return item{}, errItemNotInCache
	}

	return ci.item, nil

}

func (c *cache) delete(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, id)
}
