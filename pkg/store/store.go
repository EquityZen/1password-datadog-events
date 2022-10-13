package store

import (
	"errors"
	"sync"
	"time"
)

type Item struct {
	Id       int64
	ItemUUID string
}

type cachedItem struct {
	Item
	expiredTimestamp int64
}

type Cache struct {
	stop  chan struct{}
	wg    sync.WaitGroup
	mu    sync.Mutex
	items map[int64]cachedItem
}

func NewCacheStore(cleanupInterval time.Duration) *Cache {
	c := &Cache{
		stop:  make(chan struct{}),
		items: make(map[int64]cachedItem),
	}

	c.wg.Add(1)
	go func(cleanupInterval time.Duration) {
		defer c.wg.Done()
		c.Cleanup(cleanupInterval)
	}(cleanupInterval)
	return c
}

func (c *Cache) Cleanup(interval time.Duration) {
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

func (c *Cache) StopClean() {
	close(c.stop)
	c.wg.Wait()
}

func (c *Cache) Update(i Item, expAt int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[i.Id] = cachedItem{
		Item:             i,
		expiredTimestamp: expAt,
	}
}

var (
	errItemNotInCache = errors.New("the user isn't in Cache")
)

func (c *Cache) Read(id int64) (Item, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ci, ok := c.items[id]
	if ok {
		return Item{}, errItemNotInCache
	}

	return ci.Item, nil

}

func (c *Cache) Delete(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, id)
}
