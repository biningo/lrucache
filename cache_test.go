package lrucache

import (
	"strconv"
	"sync"
	"testing"
)

/**
*@Author lyer
*@Date 5/16/21 21:01
*@Describe
**/

func TestCacheGet(t *testing.T) {
	cache := New(100, nil)

	cacheData := map[string]string{
		"name":    "lyer",
		"age":     "18",
		"address": "zj",
		"email":   "abdsads@gmail.com",
	}

	for k, v := range cacheData {
		cache.Add(k, v)
	}

	for gotKey, wantVal := range cacheData {
		if v, ok := cache.Get(gotKey); !ok || v != wantVal {
			t.Errorf("%s got %s but want %s", gotKey, v, wantVal)
		}
	}
}

func TestCacheAdd(t *testing.T) {

	cache := New(100, nil)

	cacheData := map[string]string{
		"name":    "lyer",
		"age":     "18",
		"address": "zj",
		"email":   "abdsads@gmail.com",
	}
	wantLen := 0
	for k, v := range cacheData {
		cache.Add(k, v)
		wantLen += len(k) + len(v)
	}

	if cache.size != int64(wantLen) {
		t.Errorf("want %d but got %d", wantLen, cache.size)
	}
}

func TestCacheRemove(t *testing.T) {
	cache := New(100, nil)

	cacheData := map[string]string{
		"name":    "lyer",
		"age":     "18",
		"address": "zj",
		"email":   "abdsads@gmail.com",
	}
	for k, v := range cacheData {
		cache.Add(k, v)
	}

	for k := range cacheData {
		cache.Remove(k)
		if _, ok := cache.Get(k); ok {
			t.Errorf("value is removed but also exists!")
		}
	}
}

func TestCacheUpdate(t *testing.T) {
	cache := New(100, nil)
	cacheData := map[string]string{
		"name":    "lyer",
		"age":     "18",
		"address": "zj",
		"email":   "abdsads@gmail.com",
	}
	newCacheDate := map[string]string{
		"name":    "lyer2",
		"age":     "19",
		"address": "bj",
		"email":   "AAA@gmail.com",
	}

	for k, v := range cacheData {
		cache.Add(k, v)
	}

	for k, v := range newCacheDate {
		cache.Update(k, v)
	}

	for k, newVal := range newCacheDate {
		if v, _ := cache.Get(k); v != newVal {
			t.Errorf("want %s but got %s", newVal, v)
		}
	}
}

func TestCacheRemoveOldest(t *testing.T) {
	cache := New(2+4+6+8, nil)
	cache.Add("a", "a")
	cache.Add("bb", "bb")
	cache.Add("ccc", "ccc")
	cache.Add("dddd", "dddd")

	cache.Get("a")
	cache.removeOldest()
	if v, ok := cache.Get("bb"); ok {
		t.Errorf("want %s remove but exists", v)
	}

	cache.Get("ccc")
	cache.Add("aaa", "aaa")
	if v, ok := cache.Get("dddd"); ok {
		t.Errorf("want %s remove but exists", v)
	}
}

func TestCacheKeys(t *testing.T) {
	cache := New(2+4+6+8, nil)
	cacheData := map[string]string{
		"a":    "a",
		"bb":   "bb",
		"ccc":  "ccc",
		"dddd": "dddd",
	}
	for k, v := range cacheData {
		cache.Add(k, v)
	}

	gotKeys := cache.Keys()
	if len(gotKeys) != len(cacheData) {
		t.Errorf("keys len is not equal cacheDate")
	}
	for _, k := range gotKeys {
		if _, ok := cacheData[k]; !ok {
			t.Errorf("%s is not exists", k)
		}
	}
}

func TestCacheValues(t *testing.T) {
	cache := New(2+4+6+8, nil)
	cacheData := map[string]string{
		"a":    "a",
		"bb":   "bb",
		"ccc":  "ccc",
		"dddd": "dddd",
	}
	wantValues := map[string]bool{}
	for k, v := range cacheData {
		cache.Add(k, v)
		wantValues[v] = true
	}

	gotValues := cache.Values()
	if len(gotValues) != len(cacheData) {
		t.Errorf("values len is not equal cacheDate")
	}

	for _, v := range gotValues {
		if _, ok := wantValues[v]; !ok {
			t.Errorf("%s is not exists", v)
		}
	}

}

func TestCacheConcurrent(t *testing.T) {
	cache := New(1000, nil)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			cache.Add(k, k)
		}(strconv.Itoa(i))
	}
	wg.Wait()
	for i := 0; i < 100; i++ {
		go func(k string) {
			if _, ok := cache.Get(k); !ok {
				t.Errorf("%s is not exists", k)
			}
		}(strconv.Itoa(i))
	}
}
