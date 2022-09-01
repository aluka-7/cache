package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aluka-7/configuration"
	"sync"
	"time"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]Driver, 0)
)

func Register(name string, provider Driver) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if provider == nil {
		panic("cache: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("cache: Register called twice for provider " + name)
	}
	providers[name] = provider
}

func Read(key string) (Driver, bool) {
	prov, ok := providers[key]
	return prov, ok
}

/**
 * 缓存引擎定义，通过缓存引擎获取缓存客户端并进行数据缓存操作。系统根据业务类型划分了几个特定的缓存节点类型，每个类型的
 * 缓存节点可以分别指定自己的缓存实现方式，通过配置中心的配置示例如下：
 */

// 缓存节点服务器的类型，不同节点类型缓存的数据及其目的有所差异，业务系统要根据实际情况进行选择处理。
func Engine(systemId string, cfg configuration.Configuration) (prov Provider) {
	fmt.Println("Loading Cache Engine")
	if data, err := cfg.String("base", "cache", "", systemId); err == nil {
		var conf map[string]string
		if err := json.Unmarshal([]byte(data), &conf); err == nil {
			if v, ok := conf["provider"]; ok {
				if _prov, _ok := Read(v); _ok {
					prov = _prov.New(conf)
				} else {
					panic("加载缓存实现出错!")
				}
			} else {
				panic("没有定义缓存实现信息!")
			}
		} else {
			panic("解析缓存配置信息出错!")
		}
	} else {
		panic("加载缓存引擎配置发生错误!")
	}
	return
}

// Driver 缓存驱动程序定义
type Driver interface {
	// New 初始化缓存接口
	New(ctx context.Context, cfg map[string]string) Provider
}

// Provider 分布式的缓存操作接口。
type Provider interface {
	XAdd(ctx context.Context, key, id string, value interface{}) bool

	XGroup(ctx context.Context, key, group, start string) bool

	XReadGroup(ctx context.Context, group, consumer string, keys []string, count int64, block time.Duration) (interface{}, bool)

	XRead(ctx context.Context, keys []string, count int64, block time.Duration) (interface{}, bool)

	XAck(ctx context.Context, key, group string, ids ...string) bool

	// Exists
	// @description 判断缓存中是否存在指定的key
	// @param key
	// @return string
	Exists(ctx context.Context, key string) bool

	// String
	// @description 根据给定的key从分布式缓存中读取数据并返回，如果不存在或已过期则返回Null。
	// @param key 缓存唯一键
	// @return string
	String(ctx context.Context, key string) string

	// Set
	// @description 使用指定的key将对象存入分布式缓存中，并使用缓存的默认过期设置，注意，存入的对象必须是可序列化的。
	// @param key   缓存唯一键
	// @param value 对应的值
	Set(ctx context.Context, key, value string) bool

	// SetExpires
	// @description 使用指定的key将对象存入分部式缓存中，并指定过期时间，注意，存入的对象必须是可序列化的
	// @param key     缓存唯一键
	// @param value   对应的值
	// @param expires 过期时间，单位秒
	SetExpires(ctx context.Context, key, value string, expires time.Duration) bool

	// Delete
	// @description 从缓存中删除指定key的缓存数据。
	// @param key
	// @return
	Delete(ctx context.Context, key string) bool

	// BatchDelete
	// @description 批量删除缓存中的key
	// @param keys
	BatchDelete(ctx context.Context, keys ...string) bool

	// HSet
	// @description 将指定key的map数据的某个字段设置为给定的值
	// @param key   map数据的键
	// @param field map的字段名称
	// @param value 要设置的字段值
	HSet(ctx context.Context, key, field, value string) bool

	// HGet
	// @description 获取指定key的map数据某个字段的值，如果不存在则返回Null
	// @param key   map数据的键
	// @param field map的字段名称
	// @return
	HGet(ctx context.Context, key, field string) string

	// HGetAll
	// @description 获取指定key的map对象，如果不存在则返回Null
	// @param key map数据的键
	// @return map[string]string
	HGetAll(ctx context.Context, key string) map[string]string

	// HDelete
	// 将指定key的map数据中的某个字段删除。
	// @param key   map数据的键
	// @param field map中的key名称
	HDelete(ctx context.Context, key string, fields ...string) bool

	// HExists
	// 判断缓存中指定key的map是否存在指定的字段，如果key或字段不存在则返回false。
	// @param key
	// @param field
	// @return bool
	HExists(ctx context.Context, key, field string) bool

	// Val
	// 对指定的key结果集执行指定的脚本并返回最终脚本执行的结果。
	// @param script 脚本
	// @param key    要操作的缓存key
	// @param args   脚本的参数列表
	// @return
	Val(ctx context.Context, script string, keys []string, args ...interface{}) string

	// Operate
	// 通过直接调用缓存客户端进行缓存操作，该操作适用于高级操作，如果执行失败会返回Null。
	// @param operator
	// @return err
	Operate(ctx context.Context, cmd interface{}) error

	// Close
	// 关闭客户端
	Close()
}
