package metadata

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"

	"github.com/google/btree"
)

type TaskManagement struct {
	sync.RWMutex

	// tenantId -> Tasks
	tasks       map[uint32]*TenantTasks
	dataIsReady atomic.Bool
}

func NewTaskManagement() *TaskManagement {
	return &TaskManagement{
		tasks: make(map[uint32]*TenantTasks),
	}
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
func (t *TaskManagement) LoadAllTasks(db *utils.DBPool) error {
	// Clean up the tasks.
	t.Lock()
	t.tasks = make(map[uint32]*TenantTasks)
	t.Unlock()

	// Load tasks from the db.
	rows, err := db.Query("SELECT * FROM tistream.tasks")
	if err != nil {
		return fmt.Errorf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, tenantID, dispatcher, sorter uint32
		var rangeStart, rangeEnd, snapshotAddr string
		if err := rows.Scan(&id, &tenantID, &rangeStart, &rangeEnd, &dispatcher, &sorter, &snapshotAddr); err != nil {
			return fmt.Errorf("failed to scan tasks: %v", err)
		}

		pbTask := &pb.Task{
			TenantId:   uint32(tenantID),
			RangeStart: []byte(rangeStart),
			RangeEnd:   []byte(rangeEnd),
		}

		task := TenantTask{
			RangeStart:     []byte(rangeStart),
			Task:           pbTask,
			AssignedSorter: sorter,
		}

		t.AddTask(tenantID, task)
	}

	t.SetReady(true)
	return nil
}

func (t *TaskManagement) AddTask(tenantId uint32, task TenantTask) {
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

func (t *TaskManagement) RemoveTask(tenantId uint32, rangeStart RangeStart) {
	t.Lock()
	defer t.Unlock()

	if tasks, ok := t.tasks[tenantId]; ok {
		tasks.RemoveTask(rangeStart)
	}
}

func (t *TaskManagement) RemoveTenant(tenantId uint32) {
	t.Lock()
	defer t.Unlock()

	delete(t.tasks, tenantId)
}

func (t *TaskManagement) GetTask(tenantId uint32, rangeStart RangeStart) *TenantTask {
	t.RLock()
	defer t.RUnlock()

	if tasks, ok := t.tasks[tenantId]; ok {
		return tasks.GetTask(rangeStart)
	}

	return nil
}

type TenantTasks struct {
	sync.RWMutex

	tenantId  uint32
	taskBtree *btree.BTreeG[TenantTask]
}

func NewTenantTasks(tenantId uint32) *TenantTasks {
	return &TenantTasks{
		tenantId:  tenantId,
		taskBtree: btree.NewG(32, taskLess),
	}
}

func (t *TenantTasks) AddTask(task TenantTask) {
	t.Lock()
	defer t.Unlock()

	t.taskBtree.ReplaceOrInsert(task)
}

func (t *TenantTasks) RemoveTask(rangeStart RangeStart) {
	t.Lock()
	defer t.Unlock()

	t.taskBtree.Delete(TenantTask{RangeStart: rangeStart})
}

func (t *TenantTasks) GetTask(rangeStart RangeStart) *TenantTask {
	t.RLock()
	defer t.RUnlock()

	if item, ok := t.taskBtree.Get(TenantTask{RangeStart: rangeStart}); ok {
		return &item
	}

	return nil
}
