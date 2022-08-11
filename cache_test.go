package cache_test

import (
	"testing"

	"github.com/aluka-7/cache"
	"github.com/aluka-7/configuration"
	"github.com/aluka-7/configuration/backends"
)

func TestCacheApp(t *testing.T) {
	conf := configuration.MockEngine(t, backends.StoreConfig{Exp: map[string]string{
		"/system/base/cache/10000": "{\"provider\":\"mem\",\"mem\":\"1000\"}",
	}})
	biz := cache.Engine("10000", conf)
	expected := "test set"
	biz.Set("test", expected)
	actual := biz.String("test")
	if actual != expected {
		t.Error("生成的结果不匹配\n", "预期:", expected, "|", "实际:", actual)
	}
}
