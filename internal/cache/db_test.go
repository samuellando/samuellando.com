package cache

import (
	"bytes"
	"testing"
	"time"
)

func TestDbCacheGet_Hit(t *testing.T) {
	db := setup()
	defer teardown(db)

	// Insert a valid cache entry
	key := "test_key"
	value := []byte("test_value")
	validTo := time.Now().Add(time.Hour)
	_, err := db.Exec(`INSERT INTO cache (cache_key, cache_value, valid_to) VALUES ($1, $2, $3)`, key, value, validTo)
	if err != nil {
		t.Fatalf("Failed to insert cache entry: %v", err)
	}

	options := CacheOptions{Db: db}
	elem, err := dbCacheGet(key, options)
	if err != nil {
		t.Fatalf("Expected cache hit, got error: %v", err)
	}
	if !bytes.Equal(elem.value, value) {
		t.Fatalf("Expected value %s, got %s", value, elem.value)
	}
}

func TestDbCacheGet_Miss(t *testing.T) {
	db := setup()
	defer teardown(db)

	options := CacheOptions{Db: db}
	_, err := dbCacheGet("non_existent_key", options)
	if err == nil {
		t.Fatal("Expected cache miss error, got nil")
	}
}

func TestDbCacheGet_DbError(t *testing.T) {
	options := CacheOptions{Db: nil}
	_, err := dbCacheGet("any_key", options)
	if err == nil {
		t.Fatal("Expected database error, got nil")
	}
}

func TestDbCacheUpdate_Insert(t *testing.T) {
	db := setup()
	defer teardown(db)

	key := "insert_key"
	value := []byte("insert_value")
	validTo := time.Now().Add(time.Hour)
	elem := cacheElement{value: value, validTo: validTo}

	options := CacheOptions{Db: db}
	dbCacheUpdate(key, elem, options)

	var dbValue []byte
	var dbValidTo time.Time
	err := db.QueryRow(`SELECT cache_value, valid_to FROM cache WHERE cache_key = $1`, key).Scan(&dbValue, &dbValidTo)
	if err != nil {
		t.Fatalf("Failed to retrieve cache entry: %v", err)
	}
	if !bytes.Equal(dbValue, value) {
		t.Fatalf("Expected value %s, got %s", value, dbValue)
	}
}

func TestDbCacheUpdate_Update(t *testing.T) {
	db := setup()
	defer teardown(db)

	key := "update_key"
	initialValue := []byte("initial_value")
	updatedValue := []byte("updated_value")
	validTo := time.Now().Add(time.Hour)

	// Insert initial cache entry
	_, err := db.Exec(`INSERT INTO cache (cache_key, cache_value, valid_to) VALUES ($1, $2, $3)`, key, initialValue, validTo)
	if err != nil {
		t.Fatalf("Failed to insert initial cache entry: %v", err)
	}

	// Update cache entry
	elem := cacheElement{value: updatedValue, validTo: validTo}
	options := CacheOptions{Db: db}
	dbCacheUpdate(key, elem, options)

	var dbValue []byte
	err = db.QueryRow(`SELECT cache_value FROM cache WHERE cache_key = $1`, key).Scan(&dbValue)
	if err != nil {
		t.Fatalf("Failed to retrieve updated cache entry: %v", err)
	}
	if !bytes.Equal(dbValue, updatedValue) {
		t.Fatalf("Expected updated value %s, got %s", updatedValue, dbValue)
	}
}
