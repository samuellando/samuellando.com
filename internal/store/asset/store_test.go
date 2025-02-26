package asset

import (
	"testing"
)

func TestStoreAdd(t *testing.T) {
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
		t.Error("expected asset ID to be non-zero after adding")
	}
}

func TestStoreGetById(t *testing.T) {
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

	retrievedAsset, err := store.GetById(asset.Id())
	if err != nil {
		t.Fatalf("failed to get asset by ID: %v", err)
	}

	if retrievedAsset.Name() != asset.Name() {
		t.Errorf("expected asset Name to be %s, got %s", asset.Name(), retrievedAsset.Name())
	}
}

func TestStoreGetByName(t *testing.T) {
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

	retrievedAsset, err := store.GetByName(asset.Name())
	if err != nil {
		t.Fatalf("failed to get asset by name: %v", err)
	}

	if retrievedAsset.Id() != asset.Id() {
		t.Errorf("expected asset id to be %d, got %d", asset.Id(), retrievedAsset.Id())
	}
}

func TestStoreGetAll(t *testing.T) {
	store, _ := setup()
	defer teardown(store)

	asset1 := CreateProto(func(fields *AssetFields) {
		fields.Name = "Asset 1"
		fields.Content = []byte("Content 1")
	})
	asset2 := CreateProto(func(fields *AssetFields) {
		fields.Name = "Asset 2"
		fields.Content = []byte("Content 2")
	})

	store.Add(&asset1)
	store.Add(&asset2)

	assets, err := store.GetAll()
	if err != nil {
		t.Fatalf("failed to get all assets: %v", err)
	}

	if len(assets) != 2 {
		t.Errorf("expected 2 assets, got %d", len(assets))
	}
}

func TestStoreFilter(t *testing.T) {
	store, _ := setup()
	defer teardown(store)

	asset1 := CreateProto(func(fields *AssetFields) {
		fields.Name = "Asset 1"
		fields.Content = []byte("Content 1")
	})
	asset2 := CreateProto(func(fields *AssetFields) {
		fields.Name = "Asset 2"
		fields.Content = []byte("Content 2")
	})

	store.Add(&asset1)
	store.Add(&asset2)

	filteredStore := store.Filter(func(a *Asset) bool {
		return a.Name() == "Asset 1"
	})

	assets, err := filteredStore.GetAll()
	if err != nil {
		t.Fatalf("failed to filter assets: %v", err)
	}

	if len(assets) != 1 || assets[0].Name() != "Asset 1" {
		t.Errorf("expected 1 asset with name 'Asset 1', got %d with name %s", len(assets), assets[0].Name())
	}
}

func TestStoreGroup(t *testing.T) {
	store, _ := setup()
	defer teardown(store)

	asset1 := CreateProto(func(fields *AssetFields) {
		fields.Name = "Asset 1"
		fields.Content = []byte("Content 1")
	})
	asset2 := CreateProto(func(fields *AssetFields) {
		fields.Name = "Asset 2"
		fields.Content = []byte("Content 2")
	})

	store.Add(&asset1)
	store.Add(&asset2)

	grouped := store.Group(func(a *Asset) string {
		return a.Name()
	})

	if grouped.Len() != 2 {
		t.Errorf("expected 2 groups, got %d", grouped.Len())
	}
}

func TestStoreSort(t *testing.T) {
	store, _ := setup()
	defer teardown(store)

	asset1 := CreateProto(func(fields *AssetFields) {
		fields.Name = "Asset B"
		fields.Content = []byte("Content B")
	})
	asset2 := CreateProto(func(fields *AssetFields) {
		fields.Name = "Asset A"
		fields.Content = []byte("Content A")
	})

	store.Add(&asset1)
	store.Add(&asset2)

	sortedStore := store.Sort(func(a1, a2 *Asset) bool {
		return a1.Name() < a2.Name()
	})

	assets, err := sortedStore.GetAll()
	if err != nil {
		t.Fatalf("failed to sort assets: %v", err)
	}

	if assets[0].Name() != "Asset A" || assets[1].Name() != "Asset B" {
		t.Errorf("expected assets sorted by name, got %s and %s", assets[0].Name(), assets[1].Name())
	}
}
