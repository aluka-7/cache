package cache

import (
	"container/list"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aluka-7/utils"
)

func init() {
	Register("mem", &memDriver{})
}

// AtomicInt 是要原子访问的int64。
type AtomicInt int64

// Add atomically adds n to i.
func (i *AtomicInt) Add(n int64) {
	atomic.AddInt64((*int64)(i), n)
}

// Get atomically gets the value of i.
func (i *AtomicInt) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}

type memDriver struct{}

func (d memDriver) New(cfg map[string]string) Provider {
	fmt.Println("Loading Memory Cache Engine")
	return &defaultProvider{
		maxItemSize: utils.StrTo(cfg["mem"]).MustInt(),
		cacheList:   list.New(),
		cache:       make(map[interface{}]*list.Element),
	}
}

type defaultProvider struct {
	mutex       sync.RWMutex
	maxItemSize int
	cacheList   *list.List
	cache       map[interface{}]*list.Element
	hits, gets  AtomicInt
}
type entry struct {
	key   interface{}
	value interface{}
}

func (d *defaultProvider) Exists(key string) bool {
	return d.String(key) != ""
}

func (d *defaultProvider) String(key string) string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	d.gets.Add(1)
	if ele, hit := d.cache[key]; hit {
		d.hits.Add(1)
		d.cacheList.MoveToFront(ele)
		return utils.ToStr(ele.Value.(*entry).value)
	}
	return ""
}

func (d *defaultProvider) GetByProvider(key string, provider DataProvider) string {
	panic("implement me")
}

func (d *defaultProvider) Set(key, value string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.cache == nil {
		d.cache = make(map[interface{}]*list.Element)
		d.cacheList = list.New()
	}

	if ele, ok := d.cache[key]; ok {
		d.cacheList.MoveToFront(ele)
		ele.Value.(*entry).value = value
		return
	}

	ele := d.cacheList.PushFront(&entry{key: key, value: value})
	d.cache[key] = ele
	if d.maxItemSize != 0 && d.cacheList.Len() > d.maxItemSize {
		d.RemoveOldest()
	}
}

// RemoveOldest remove the oldest key
func (d *defaultProvider) RemoveOldest() {
	if d.cache == nil {
		return
	}
	ele := d.cacheList.Back()
	if ele != nil {
		d.cacheList.Remove(ele)
		key := ele.Value.(*entry).key
		delete(d.cache, key)
	}
}
func (d *defaultProvider) SetExpires(key, value string, expires time.Duration) bool {
	d.Set(key, value)
	return true
}

func (d *defaultProvider) Delete(key string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.cache == nil {
		return
	}
	if ele, ok := d.cache[key]; ok {
		d.cacheList.Remove(ele)
		key := ele.Value.(*entry).key
		delete(d.cache, key)
		return
	}
}

func (d *defaultProvider) BatchDelete(keys ...string) {
	panic("implement me")
}

func (d *defaultProvider) HSet(key, field, value string) {
	d.Set(key+field, value)
}

func (d *defaultProvider) HGet(key, field string) string {
	return d.String(key + field)
}

func (d *defaultProvider) HGetAll(key string) map[string]string {
	panic("implement me")
}

func (d *defaultProvider) HDelete(key string, fields ...string) {
	panic("implement me")
}

func (d *defaultProvider) HExists(key, field string) bool {
	panic("implement me")
}

func (d *defaultProvider) Val(script string, keys []string, args ...interface{}) {
	panic("implement me")
}

func (d *defaultProvider) Operate(interface{}) error {
	panic("implement me")
}

func (d *defaultProvider) Close() {
	panic("implement me")
}
