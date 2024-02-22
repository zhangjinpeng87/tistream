package metadata

import (
	"bytes"
	"sync"

	"github.com/google/btree"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type RangeStart []byte

type TenantTask struct {
	// RangeStart is the start of the range.
	RangeStart

	// Task contains rangeStart, rangeEnd, tenantId, and other information.
	InternalTask *pb.Task
}

func taskLess(a, b TenantTask) bool {
	return bytes.Compare(a.RangeStart, b.RangeStart) < 0
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
