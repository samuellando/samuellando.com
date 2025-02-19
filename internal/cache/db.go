package cache

import (
	"fmt"
	"log"
	"time"
)

// Attempt to retrieve a cache entry from the external database.
// It returns the cached data if found and valid, otherwise returns an error.
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

// Insert or update a cache entry in the external database.
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
