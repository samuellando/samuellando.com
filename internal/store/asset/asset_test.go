package asset

import (
	"database/sql"
	"testing"
	"time"

	"samuellando.com/internal/db"
	"samuellando.com/internal/testutil"
)

func setup() (Store, *sql.DB) {
	if err := testutil.ResetDb(); err != nil {
		panic(err)
	}
	con := db.ConnectPostgres(testutil.GetDbCredentials())
	migrations, err := testutil.GetMigrationsPath()
	if err != nil {
		panic(err)
	}
	if err := db.ApplyMigrations(con, func(o *db.Options) {
		o.MigrationsDir = migrations
		o.Logger = testutil.CreateDiscardLogger()
	}); err != nil {
		panic(err)
	}
	return CreateStore(con), con
}

func teardown(s Store) {
	s.db.Close()
	testutil.ResetDb()
}

func TestCreateProto(t *testing.T) {
	asset := CreateProto(func(fields *AssetFields) {
		fields.Name = "Test Asset"
		fields.Content = []byte("Test Content")
	})

	if asset.Name() != "Test Asset" {
		t.Errorf("expected name to be 'Test Asset', got %s", asset.Name())
	}
	if string(asset.content) != "Test Content" {
		t.Errorf("expected content to be 'Test Content', got %s", string(asset.content))
	}
	if !asset.loaded {
		t.Error("expected asset to be loaded")
	}
}

func TestAssetId(t *testing.T) {
	store, _ := setup()
	defer teardown(store)

	asset := CreateProto(func(fields *AssetFields) {
		fields.Name = "Test Asset"
		fields.Content = []byte("Test Content")
	})

	err := store.Add(&asset)
	if err != nil {
		t.Fatalf("failed to add asset: %v", err)
	}

	if asset.Id() == 0 {
		t.Error("expected asset ID to be non-zero")
	}
}

func TestAssetName(t *testing.T) {
	store, _ := setup()
	defer teardown(store)

	asset := CreateProto(func(fields *AssetFields) {
		fields.Name = "Test Asset"
		fields.Content = []byte("Test Content")
	})

	err := store.Add(&asset)
	if err != nil {
		t.Fatalf("failed to add asset: %v", err)
	}

	if asset.Name() != "Test Asset" {
		t.Errorf("expected name to be 'Test Asset', got %s", asset.Name())
	}
}

func TestAssetCreated(t *testing.T) {
	store, _ := setup()
	defer teardown(store)

	asset := CreateProto(func(fields *AssetFields) {
		fields.Name = "Test Asset"
		fields.Content = []byte("Test Content")
	})

	err := store.Add(&asset)
	if err != nil {
		t.Fatalf("failed to add asset: %v", err)
	}

	if time.Since(asset.Created()) > time.Second {
		t.Error("expected created time to be within the last second")
	}
}

func TestAssetContent(t *testing.T) {
	store, _ := setup()
	defer teardown(store)

	asset := CreateProto(func(fields *AssetFields) {
		fields.Name = "Test Asset"
		fields.Content = []byte("Test Content")
	})

	err := store.Add(&asset)
	if err != nil {
		t.Fatalf("failed to add asset: %v", err)
	}

	content, err := asset.Content()
	if err != nil {
		t.Fatalf("failed to get content: %v", err)
	}
	if string(content) != "Test Content" {
		t.Errorf("expected content to be 'Test Content', got %s", string(content))
	}
}

func TestAssetDelete(t *testing.T) {
	store, _ := setup()
	defer teardown(store)

	asset := CreateProto(func(fields *AssetFields) {
		fields.Name = "Test Asset"
		fields.Content = []byte("Test Content")
	})

	err := store.Add(&asset)
	if err != nil {
		t.Fatalf("failed to add asset: %v", err)
	}

	err = asset.Delete()
	if err != nil {
		t.Fatalf("failed to delete asset: %v", err)
	}

	_, err = store.GetById(asset.Id())
	if err == nil {
		t.Error("expected error when getting deleted asset")
	}
}
