package metadata

import (
	"sync"
	"sync/atomic"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type TaskManagement struct {
	sync.RWMutex

	// tenantId -> Tasks
	tasks       map[uint64]*TenantTasks
	backend     *Backend
	dataIsReady atomic.Bool
}

func NewTaskManagement(config *utils.MetaServerConfig, backend *Backend) *TaskManagement {
	return &TaskManagement{
		tasks:   make(map[uint64]*TenantTasks),
		backend: backend,
	}
}

func (t *TaskManagement) Prepare() error {
	return t.backend.BootstrapSchema()
}

// DataIsReady returns true if the data is ready.
func (t *TaskManagement) DataIsReady() bool {
	return t.dataIsReady.Load()
}

func (t *TaskManagement) SetReady(b bool) {
	t.dataIsReady.Store(b)
}

// LoadAllTasks loads all tasks from the db. The db is the source of truth.
// This function should be called when the meta-server campaign as master.
func (t *TaskManagement) LoadAllTasks() error {
	// Clean up the tasks.
	t.Lock()
	t.tasks = make(map[uint64]*TenantTasks)
	t.Unlock()

	// Load all tasks from the backend.
	t.backend.LoadAllTasks(func(tenantId uint64, task *pb.Task) {
		t.AddTask(tenantId, task)
	})

	t.SetReady(true)
	return nil
}

func (t *TaskManagement) AddTask(tenantId uint64, task *pb.Task) {
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

func (t *TaskManagement) RemoveTask(tenantId uint64, task *pb.Task) {
	t.Lock()
	defer t.Unlock()

	if tasks, ok := t.tasks[tenantId]; ok {
		tasks.RemoveTask(task)
	}
}

func (t *TaskManagement) RemoveTenant(tenantId uint64) {
	t.Lock()
	defer t.Unlock()

	delete(t.tasks, tenantId)
}

func (t *TaskManagement) GetTask(tenantId uint64, rangeStart []byte) *pb.Task {
	t.RLock()
	defer t.RUnlock()

	if tasks, ok := t.tasks[tenantId]; ok {
		return tasks.GetTask(rangeStart)
	}

	return nil
}
