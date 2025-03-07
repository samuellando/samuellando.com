package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"samuellando.com/data"
)

// Attempt to retrieve a cache entry from the external database.
// It returns the cached data if found and valid, otherwise returns an error.
func dbCacheGet(key string, options CacheOptions) (cacheElement, error) {
	if options.Db != nil {
		log.Println("Checking external cache")
		ctx := context.TODO()
		queries := data.New(options.Db)
		row, err := queries.GetCacheByKey(ctx, key)
		if err != nil {
			log.Println("Cache from db error: ", err)
			return cacheElement{}, err
		}
		if time.Since(row.ValidTo) <= 0 {
			log.Println("External cache hit", time.Since(row.ValidTo))
			return cacheElement{value: row.CacheValue, validTo: row.ValidTo}, nil
		}
		return cacheElement{}, fmt.Errorf("db cache miss")
	} else {
		return cacheElement{}, fmt.Errorf("No db provided")
	}
}

// Insert or update a cache entry in the external database.
func dbCacheUpdate(key string, elem cacheElement, options CacheOptions) {
	if options.Db != nil {
		ctx := context.TODO()
		queries := data.New(options.Db)
		err := queries.SetCacheByKey(ctx, data.SetCacheByKeyParams{
			CacheKey:   key,
			CacheValue: elem.value,
			ValidTo:    elem.validTo,
		})
		if err != nil {
			log.Println(err)
		}
	}
}
