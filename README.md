# 缓存引擎

缓存引擎定义，通过缓存引擎获取缓存客户端并进行数据缓存操作。系统根据业务类型划分了几个特定的缓存节点类型，每个类型的缓存节点可以分别指定自己的缓存实现方式，通过配置中心的配置示例如下：
```
{
     "Node" : "provider",
     "Node" : "provider"
}
```

其中NodeType对应缓存节点类型{@link CacheNode}的枚举值（要大写），其值对应缓存的实现，有如下选项：
```
* redis：基于单Redis实例的实现；
* redis-sentinel：基于Redis实例主从模式的高可用集群实现；
* mem:基于内存实现的缓存部分接口,以供测试使用.
```

配置中心的配置路径为：/system/base/cache/provider

## 缓存Redis单实例

基于Redis单实例的缓存实现(单机），配置中心或构造方法中参数的配置格式如下：
```
{
    "host" : "缓存服务器主机地址，必须",
    "port" : "缓存服务器端口号",
    "sasl" : "是否开启安全认证，true|false，可选，默认是没有",
    "password" : "开启安全认证后的登录密码，sasl如果指定为true则必须"
 }
```

## 快速使用
<font color=red size=72>注意</font>:
1. 一定在导入缓存具体实现如redis
2. 注意缓存接口与缓存实现的导入顺序

示例：
```go
    _ "github.com/aluka-7/cache-redis"
    "github.com/aluka-7/cache"
```

### 获取缓存实体
```go
prov:=cache.Engine().OptProvider()
```

### 判断缓存中是否存在指定的key
```go
prov.Exists(key string) bool
```

### 根据给定的key从分布式缓存中读取数据并返回，如果不存在或已过期则返回Null。
```go
prov.String(key string) string
```

### 使用给定的key从缓存中查询数据，如果查询不到则使用给定的数据提供器来查询数据，然后将数据存入缓存中再返回。
```go
prov.GetByProvider(key string, provider DataProvider) string
```

### 使用指定的key将对象存入分布式缓存中，并使用缓存的默认过期设置，注意，存入的对象必须是可序列化的。
```go
prov.Set(key, value string)
```

### 使用指定的key将对象存入分部式缓存中，并指定过期时间，注意，存入的对象必须是可序列化的
```go
prov.SetExpires(key, value string, expires time.Duration) bool
```

### 从缓存中删除指定key的缓存数据。
```go
prov.Delete(key string)
```

### 批量删除缓存中的key。
```go
prov.BatchDelete(keys ...string)
```

### 将指定key的map数据的某个字段设置为给定的值。
```go
prov.HSet(key, field, value string)
```

### 获取指定key的map数据某个字段的值，如果不存在则返回Null
```go
prov.HGet(key, field string) string
```

### 获取指定key的map对象，如果不存在则返回Null
```go
prov.HGetAll(key string) map[string]string
```

### 将指定key的map数据中的某个字段删除。
```go
prov.HDelete(key string, fields ...string)
```

### 判断缓存中指定key的map是否存在指定的字段，如果key或字段不存在则返回false。
```go
prov.HExists(key, field string) bool
```

### 对指定的key结果集执行指定的脚本并返回最终脚本执行的结果。
```go
prov.Val(script string, keys []string, args ...interface{})
```
### 根据指定的key进行自增1
```go
prov.AtomicIncrement(key string)int64
```

### 根据指定的key进行自减1
```go
prov.AtomicDecrement(key string)int64
```

### 关闭客户端
```go
prov.Close()
```
