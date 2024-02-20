package metadata

import (
	"sync"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"

	"github.com/google/btree"
)

type TaskManagement struct {
	sync.RWMutex

	// tenantId -> Tasks
	tasks map[uint32]*TenantTasks
}

func NewTaskManagement() *TaskManagement {
	return &TaskManagement{
		tasks: make(map[uint32]*TenantTasks),
	}
}

func (t *TaskManagement) AddTask(tenantId uint32, task *pb.Task) {
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

func (t *TaskManagement) GetTask(tenantId uint32, rangeStart RangeStart) *pb.Task {
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

func (t *TenantTasks) AddTask(task *pb.Task) {
	t.Lock()
	defer t.Unlock()

	t.taskBtree.ReplaceOrInsert(TenantTask{RangeStart: task.RangeStart, Task: task})
}

func (t *TenantTasks) RemoveTask(rangeStart RangeStart) {
	t.Lock()
	defer t.Unlock()

	t.taskBtree.Delete(TenantTask{RangeStart: rangeStart})
}

func (t *TenantTasks) GetTask(rangeStart RangeStart) *pb.Task {
	t.RLock()
	defer t.RUnlock()

	if item, ok := t.taskBtree.Get(TenantTask{RangeStart: rangeStart}); ok {
		return item.Task
	}

	return nil
}
