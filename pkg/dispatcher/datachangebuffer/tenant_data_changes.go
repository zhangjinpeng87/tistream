package datachangebuffer

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/zhangjinpeng87/tistream/pkg/storage"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
)

const (
	schemaSnap = "schema_snap"
)

// The data change buffer files for a specified tenant are organized as follows:
// |____data_change_buffer/cluster-{1}/Tenant-{1}
// |  |____schema_snap
// |  |____{store-id1}
// |  |  |____file-{ts}
// |  |  |____file-{ts}
// |  |____{store-id1}
// |  |  |____file-{ts}
// |  |  |____file-{ts}
// |  |____{store-id1}
// |     |____file-{ts}
//
// There are 2 types of files:
//  1. schema_snap: the schema snapshot file.
//  2. file-{ts}: the data change file. It contains the data change events for a specified
//     time range for a specified store.

// Tenant is the data change buffer for a tenant.
type TenantDataChanges struct {
	// The tenant id.
	tenantID uint64

	// Root directory of the data change buffer for the tenant.
	// In above file org it should be "{prefix}data_change_buffer/cluster-{1}/Tenant-{1}"
	rootDir string

	// Backend Storage
	backendStorage storage.BackendStorage

	// Mutex to protect the storesProgress.
	mu sync.Mutex
	// StoreID -> StoreProgress
	storesProgress map[string]*StoreProgress
}

// NewTenantDataChanges creates a new TenantDataChanges.
func NewTenantDataChanges(tenantID uint64, rootDir string, backendStorage storage.BackendStorage) *TenantDataChanges {
	return &TenantDataChanges{
		tenantID:       tenantID,
		rootDir:        rootDir,
		backendStorage: backendStorage,
	}
}

// initialize initializes the tenant data changes.
func (t *TenantDataChanges) initialize() error {
	// Lock the mutex.
	t.mu.Lock()
	defer t.mu.Unlock()

	// List all the store directories.
	storeDirs, err := t.listStoreDir()
	if err != nil {
		return err
	}

	// Initialize the store progress.
	t.storesProgress = make(map[string]*StoreProgress)
	for storeDir, _ := range storeDirs {
		storeID, err := strconv.ParseUint(storeDir, 10, 64)
		if err != nil {
			return err
		}

		t.storesProgress[storeDir] = NewStoreProgress(storeID, storeDir, 0)
	}

	return nil
}

func (t *TenantDataChanges) updateStores() error {
	// Lock the mutex.
	t.mu.Lock()
	defer t.mu.Unlock()

	latestStores, err := t.listStoreDir()
	if err != nil {
		return err
	}

	for storeDir, _ := range latestStores {
		storeID, err := strconv.ParseUint(storeDir, 10, 64)
		if err != nil {
			return err
		}

		if _, ok := t.storesProgress[storeDir]; !ok {
			t.storesProgress[storeDir] = NewStoreProgress(storeID, storeDir, 0)
		}
	}

	// Remove the stores that are not in the latest stores.
	// Most of the time there is no store changes, so we don't need to remove the store.
	if len(latestStores) == len(t.storesProgress) {
		return nil
	}

	// Collect all old stores.
	oldStores := make([]string, 0, len(t.storesProgress))
	for store, _ := range t.storesProgress {
		oldStores = append(oldStores, store)
	}
	for _, oldStore := range oldStores {
		if _, ok := latestStores[oldStore]; !ok {
			delete(t.storesProgress, oldStore)
		}
	}

	return nil
}

// GetSchemaSnap returns the schema snapshot of this tenant.
func (t *TenantDataChanges) GetSchemaSnap() (*SchemaSnap, error) {
	filePath := t.rootDir + schemaSnap

	content, err := t.backendStorage.GetFile(filePath)
	if err != nil {
		return nil, err
	}

	schemaSnap := NewEmptySchemaSnap()
	reader := bytes.NewReader(content)
	err = schemaSnap.DecodeFrom(reader)
	if err != nil {
		return nil, err
	}

	return schemaSnap, nil
}

// Run runs the tenant data changes.
func (t *TenantDataChanges) Run(ctx context.Context, checkStoreInterval, checkFileInterval int) error {
	// Initialize the tenant data changes.
	if err := t.initialize(); err != nil {
		return err
	}

	// Start the data change file watcher.
	t1 := time.NewTicker(time.Duration(checkStoreInterval) * time.Second)
	defer t1.Stop()
	t2 := time.NewTicker(time.Duration(checkFileInterval) * time.Second)
	defer t2.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t1.C:
			if err := t.updateStores(); err != nil {
				return err
			}
		case <-t2.C:
			if err := t.iterateNewFileChanges(); err != nil {
				return err
			}
		}
	}
}

func (t *TenantDataChanges) iterateNewFileChanges() error {
	// Lock the mutex.
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, storeProgress := range t.storesProgress {
		// Iterate the new changes for each store
		files, err := t.backendStorage.ListFiles(t.rootDir + storeProgress.storeDir)
		if err != nil {
			return err
		}

		// Sort the files.
		sort.Slice(files, func(i, j int) bool {
			return files[i] < files[j]
		})

		// Iterate the files.
		for _, file := range files {
			// Parse the file timestamp.
			ts, err := strconv.ParseUint(file, 10, 64)
			if err != nil {
				return err
			}

			// If the file is not handled, handle it.
			if ts > storeProgress.LatestHandledTs() {
				// Handle the file.
				if err := t.handleFile(storeProgress, ts); err != nil {
					return err
				}
				storeProgress.Advance(ts)
			}
		}
	}

	return nil
}

// handleFile handles the file.
func (t *TenantDataChanges) handleFile(storeProgress *StoreProgress, ts uint64) error {
	// Read the file.
	filePath := t.rootDir + storeProgress.storeDir + "/" + strconv.FormatUint(ts, 10)
	content, err := t.backendStorage.GetFile(filePath)
	if err != nil {
		return err
	}

	// Verify the checksum of the file.
	fileLen := len(content)
	if fileLen <= 4 {
		return utils.ErrInvalidDataChangeFile
	}
	checksum := binary.LittleEndian.Uint32(content[fileLen-4:])
	if !utils.IsChecksumMatch(checksum, content[:fileLen-4]) {
		return fmt.Errorf("file {} checksum not match", filePath)
	}

	// Decode the file.
	decoder := NewDataChangesFileDecoder()
	reader := bytes.NewReader(content[:fileLen-4])
	if err := decoder.DecodeFrom(reader); err != nil {
		return err
	}

	// Dispatch the event rows to the ranges.
	// Todo: dispatch the event rows to the ranges.

	return nil
}

// listStoreDirs returns the list of store directories.
func (t *TenantDataChanges) listStoreDir() (map[string]struct{}, error) {
	list, err := t.backendStorage.ListSubDir(t.rootDir)
	if err != nil {
		return nil, err
	}

	res := make(map[string]struct{})
	for _, dir := range list {
		res[dir] = struct{}{}
	}

	return res, nil
}

type StoreProgress struct {
	// The store id.
	storeID uint64

	// The store sub directory.
	storeDir string

	// Latest handled file timestamp.
	latestHandledTs uint64
}

// NewStoreProgress creates a new StoreProgress.
func NewStoreProgress(storeID uint64, storeDir string, latestHandledTs uint64) *StoreProgress {
	return &StoreProgress{
		storeID:         storeID,
		storeDir:        storeDir,
		latestHandledTs: latestHandledTs,
	}
}

func (s *StoreProgress) Advance(ts uint64) {
	s.latestHandledTs = ts
}

func (s *StoreProgress) LatestHandledTs() uint64 {
	return s.latestHandledTs
}
