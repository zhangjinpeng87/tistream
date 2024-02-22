package metadata

import (
	"sync"
	"sync/atomic"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
)

type DataManagement struct {
	sync.RWMutex

	// tenantId -> Tasks
	tasks       map[uint32]*TenantTasks
	backend     *Backend
	dataIsReady atomic.Bool
}

func NewDataManagement(config *utils.MetaServerConfig) (*DataManagement, error) {
	backend, err := NewBackend(config)
	if err != nil {
		return nil, err
	}

	return &DataManagement{
		tasks:   make(map[uint32]*TenantTasks),
		backend: backend,
	}, nil
}

func (t *DataManagement) Close() {
	t.backend.Close()
}

func (t *DataManagement) TryCampaignMaster(who string, leaseDur int) (bool, error) {
	return t.backend.TryCampaignMaster(who, leaseDur)
}

func (t *DataManagement) TryUpdateLease(who string, leaseDur int) (bool, error) {
	return t.backend.TryUpdateLease(who, leaseDur)
}

func (t *DataManagement) GetMaster() (string, error) {
	return t.backend.GetMaster()
}

// DataIsReady returns true if the data is ready.
func (t *DataManagement) DataIsReady() bool {
	return t.dataIsReady.Load()
}

func (t *DataManagement) SetReady(b bool) {
	t.dataIsReady.Store(b)
}

// LoadAllTasks loads all tasks from the db. The db is the source of truth.
// This function should be called when the meta-server campaign as master.
func (t *DataManagement) LoadAllTasks() error {
	// Clean up the tasks.
	t.Lock()
	t.tasks = make(map[uint32]*TenantTasks)
	t.Unlock()

	// Load all tasks from the backend.
	t.backend.LoadAllTasks(func(tenantId uint32, task *TenantTask) {
		t.AddTask(tenantId, *task)
	})

	t.SetReady(true)
	return nil
}

func (t *DataManagement) AddTask(tenantId uint32, task TenantTask) {
	t.Lock()
	defer t.Unlock()

	tasks, ok := t.tasks[tenantId]
	if !ok {
		tasks = NewTenantTasks(tenantId)
		t.tasks[tenantId] = tasks
	} else {
		tasks.AddTask(task)
	}
}

func (t *DataManagement) RemoveTask(tenantId uint32, rangeStart RangeStart) {
	t.Lock()
	defer t.Unlock()

	if tasks, ok := t.tasks[tenantId]; ok {
		tasks.RemoveTask(rangeStart)
	}
}

func (t *DataManagement) RemoveTenant(tenantId uint32) {
	t.Lock()
	defer t.Unlock()

	delete(t.tasks, tenantId)
}

func (t *DataManagement) GetTask(tenantId uint32, rangeStart RangeStart) *TenantTask {
	t.RLock()
	defer t.RUnlock()

	if tasks, ok := t.tasks[tenantId]; ok {
		return tasks.GetTask(rangeStart)
	}

	return nil
}
