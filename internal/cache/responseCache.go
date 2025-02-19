package cache

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"
)

/*
Configuration options for the cache.
MaxAge specifies the duration for which a cache entry is considered valid.
Db is an optional database connection for external caching, can be ommited validTo
disable the external cache.
*/
type CacheOptions struct {
	MaxAge time.Duration
	Db     *sql.DB
}

/*
A single cache entry.
*/
type cacheElement struct {
	validTo time.Time
	value   []byte
}

/*
Function that caches the result of the provided function f.
It uses in-memory and optional external database caching.
The cache key is derived from the function's name.
opts allows customization of cache options such as MaxAge and Db.
*/
func Cached(f func() ([]byte, error), opts ...func(*CacheOptions)) func() ([]byte, error) {
    // In memory local instance cache, so we don't hit the db with every request.
    cache := make(map[string]cacheElement)
	o := CacheOptions{
		MaxAge: time.Hour,
	}
	for _, opt := range opts {
		opt(&o)
	}
	options := o
	// Get a unique key for the function f
	fp := reflect.ValueOf(f).Pointer()
	fnDetails := runtime.FuncForPC(fp)
	hasher := sha256.New()
	hasher.Write([]byte(fnDetails.Name()))
	key := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	// Return a wrapped function
	return func() ([]byte, error) {
		log.Println("Checking cache", fnDetails.Name(), key)
		if cached, ok := cache[key]; ok && time.Since(cached.validTo) <= 0 {
			log.Println("Cache hit", time.Since(cached.validTo))
			return cached.value, nil
		} else {
			// Miss, check the external cache.
			elem, err := dbCacheGet(key, options)
			if err == nil {
                cache[key] = elem
				return elem.value, nil
			}
			log.Println("Cache miss")
            data, err := f()
			elem = cacheElement{validTo: time.Now().Add(options.MaxAge), value: data}
			dbCacheUpdate(key, elem, options)
			cache[key] = elem
			return data, err
		}
	}
}

/*
Attempt to retrieve a cache entry from the external database.
It returns the cached data if found and valid, otherwise returns an error.
*/
func dbCacheGet(key string, options CacheOptions) (cacheElement, error) {
	if options.Db != nil {
		log.Println("Checking external cache")
		var value []byte
		var validto time.Time
		row := options.Db.QueryRow(`
            SELECT 
                cache_value, validto 
            FROM cache 
            WHERE cache_key = $1
        `, key)
		if err := row.Scan(&value, &validto); err != nil {
			log.Println("Cache from db error: ", err)
			return cacheElement{}, err
		}
		if time.Since(validto) <= 0 {
			log.Println("External cache hit", time.Since(validto))
            return cacheElement{value: value, validTo: validto}, nil
		}
		return cacheElement{}, fmt.Errorf("db cache miss")
	} else {
		return cacheElement{}, fmt.Errorf("No db provided")
	}
}

/*
Insert or update a cache entry in the external database.
*/
func dbCacheUpdate(key string, elem cacheElement, options CacheOptions) {
	if options.Db != nil {
		tx, err := options.Db.Begin()
		if err != nil {
			log.Println(err)
			return
		}
		defer tx.Rollback()
		_, err = tx.Exec(`
			INSERT INTO cache (cache_key, cache_value, validto) VALUES
                ($1, $2, $3)
			ON CONFLICT (cache_key) DO UPDATE SET
                cache_value = EXCLUDED.cache_value,
                validto = EXCLUDED.validto
        `, key, elem.value, elem.validTo)
		if err != nil {
			log.Println(err)
		}
		err = tx.Commit()
		if err != nil {
			log.Println(err)
		}
	}
}
