package cache

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"log"
	"reflect"
	"runtime"
	"time"
)

// Configuration options for the cache.
// MaxAge specifies the duration for which a cache entry is considered valid.
// Db is an optional database connection for external caching, can be omitted to
// disable the external cache.
type CacheOptions struct {
	MaxAge time.Duration
	Db     *sql.DB
}

// A cache entry.
type cacheElement struct {
	validTo time.Time
	value   []byte
}

var localCache = make(map[string]cacheElement)

func resetCache() {
	for k := range localCache {
		delete(localCache, k)
	}
}

// Function that caches the result of the provided function f.
// It uses in-memory and optional external database caching.
// The cache key is derived from the function's name.
// opts allows customization of cache options such as MaxAge and Db.
func Cached(f func() ([]byte, error), opts ...func(*CacheOptions)) func() ([]byte, error) {
	// In-memory local instance cache to avoid hitting the db with every request.
	cacheOptions := CacheOptions{MaxAge: time.Hour}

	for _, opt := range opts {
		opt(&cacheOptions)
	}

	// Generate a unique key for the function f
	funcPointer := reflect.ValueOf(f).Pointer()
	funcDetails := runtime.FuncForPC(funcPointer)
	hasher := sha256.New()
	hasher.Write([]byte(funcDetails.Name()))
	cacheKey := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Return a wrapped function
	return func() ([]byte, error) {
		log.Println("Checking cache for function:", funcDetails.Name(), "with key:", cacheKey)

		if cachedElem, exists := localCache[cacheKey]; exists && time.Until(cachedElem.validTo) > 0 {
			log.Println("Cache hit, valid for:", time.Until(cachedElem.validTo))
			return cachedElem.value, nil
		}

		// Cache miss, check the external cache.
		cachedElem, err := dbCacheGet(cacheKey, cacheOptions)
		if err == nil {
			localCache[cacheKey] = cachedElem
			return cachedElem.value, nil
		}

		log.Println("Cache miss, executing function")
		data, err := f()
		if err != nil {
			return nil, err
		}

		newElem := cacheElement{validTo: time.Now().Add(cacheOptions.MaxAge), value: data}
		dbCacheUpdate(cacheKey, newElem, cacheOptions)
		localCache[cacheKey] = newElem

		return data, nil
	}
}

// Function that caches the result of the provided function f, with a paramter key.
// It uses in-memory and optional external database caching.
// The cache key is derived from the function's name.
// opts allows customization of cache options such as MaxAge and Db.
func ParamCached(f func() ([]byte, error), paramKey string, opts ...func(*CacheOptions)) func() ([]byte, error) {
	// In-memory local instance cache to avoid hitting the db with every request.
	cacheOptions := CacheOptions{MaxAge: time.Hour}

	for _, opt := range opts {
		opt(&cacheOptions)
	}

	// Generate a unique key for the function f
	funcPointer := reflect.ValueOf(f).Pointer()
	funcDetails := runtime.FuncForPC(funcPointer)
	hasher := sha256.New()
	hasher.Write([]byte(funcDetails.Name()))
	hasher.Write([]byte(paramKey))
	cacheKey := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Return a wrapped function
	return func() ([]byte, error) {
		log.Println("Checking cache for function:", funcDetails.Name(), "with key:", cacheKey)

		if cachedElem, exists := localCache[cacheKey]; exists && time.Until(cachedElem.validTo) > 0 {
			log.Println("Cache hit, valid for:", time.Until(cachedElem.validTo))
			return cachedElem.value, nil
		}

		// Cache miss, check the external cache.
		cachedElem, err := dbCacheGet(cacheKey, cacheOptions)
		if err == nil {
			localCache[cacheKey] = cachedElem
			return cachedElem.value, nil
		}

		log.Println("Cache miss, executing function")
		data, err := f()
		if err != nil {
			return nil, err
		}

		newElem := cacheElement{validTo: time.Now().Add(cacheOptions.MaxAge), value: data}
		dbCacheUpdate(cacheKey, newElem, cacheOptions)
		localCache[cacheKey] = newElem

		return data, nil
	}
}
