package cache

import (
	"bytes"
	"database/sql"
	"errors"
	"testing"
	"time"

	"samuellando.com/internal/db"
	"samuellando.com/internal/testutil"
)

// setup initializes the test database and applies migrations.
func setup() *sql.DB {
	fetchDataCallCount = 0
	con := db.ConnectPostgres(testutil.GetDbCredentials())
	if err := testutil.ResetDb(con, "cacheTests"); err != nil {
		panic(err)
	}
	return con
}

// teardown closes the test database connection.
func teardown(db *sql.DB) {
	db.Close()
}

// Test function to simulate data fetching
var fetchDataCallCount int

func fetchData() ([]byte, error) {
	fetchDataCallCount++
	return []byte("test data"), nil
}

// Test function to simulate data fetching with an error
func fetchDataWithError() ([]byte, error) {
	return nil, errors.New("fetch error")
}

// Test in-memory caching
func TestCached_InMemory(t *testing.T) {
	db := setup()
	defer teardown(db)

	cachedFunc := Cached(fetchData)

	// First call should fetch data
	data, err := cachedFunc()
	if err != nil || !bytes.Equal(data, []byte("test data")) {
		t.Fatalf("Expected 'test data', got %s, error: %v", data, err)
	}
	if fetchDataCallCount != 1 {
		t.Fatal("The call shouldve been a miss")
	}

	// Second call should hit the cache
	data, err = cachedFunc()
	if err != nil || !bytes.Equal(data, []byte("test data")) {
		t.Fatalf("Expected 'test data' from cache, got %s, error: %v", data, err)
	}
	if fetchDataCallCount != 1 {
		t.Fatal("The call shouldve been a hit")
	}
}

// Test cache expiry
func TestCached_Expiry(t *testing.T) {
	db := setup()
	defer teardown(db)

	cachedFunc := Cached(fetchData, func(o *CacheOptions) {
		o.MaxAge = time.Microsecond
	})
	// First call should fetch data
	data, err := cachedFunc()
	if err != nil || !bytes.Equal(data, []byte("test data")) {
		t.Fatalf("Expected 'test data', got %s, error: %v", data, err)
	}
	if fetchDataCallCount != 1 {
		t.Fatalf("The call shouldve been a miss, got %d", fetchDataCallCount)
	}

	// Wait for cache to expire
	time.Sleep(2 * time.Microsecond)

	// Next call should fetch data again, not from cache
	data, err = cachedFunc()
	if err != nil || !bytes.Equal(data, []byte("test data")) {
		t.Fatalf("Expected 'test data' after expiry, got %s, error: %v", data, err)
	}
	if fetchDataCallCount != 2 {
		t.Fatal("The call shouldve been a miss")
	}
}

// Test external database caching
func TestCached_ExternalDB(t *testing.T) {
	db := setup()
	defer teardown(db)

	cachedFunc := Cached(fetchData, func(o *CacheOptions) {
		o.Db = db
	})

	// First call should fetch data
	data, err := cachedFunc()
	if err != nil || !bytes.Equal(data, []byte("test data")) {
		t.Fatalf("Expected 'test data', got %s, error: %v", data, err)
	}
	if fetchDataCallCount != 1 {
		t.Fatalf("The call shouldve been a miss, count %d", fetchDataCallCount)
	}
	// Reset internal cache
	cachedFunc = Cached(fetchData, func(o *CacheOptions) {
		o.Db = db
	})

	// Second call should hit the external cache
	data, err = cachedFunc()
	if err != nil || !bytes.Equal(data, []byte("test data")) {
		t.Fatalf("Expected 'test data' from external cache, got %s, error: %v", data, err)
	}
	if fetchDataCallCount != 1 {
		t.Fatal("The call shouldve been a hit")
	}
}

// Test cache miss handling
func TestCached_CacheMiss(t *testing.T) {
	db := setup()
	defer teardown(db)

	cachedFunc := Cached(fetchDataWithError)

	// First call should attempt to fetch data and fail
	data, err := cachedFunc()
	if err == nil || data != nil {
		t.Fatalf("Expected error, got data: %s, error: %v", data, err)
	}
}
