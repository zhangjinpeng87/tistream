package metadata

import (
	"sync"

	"github.com/huandu/skiplist"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type TenantTasks struct {
	sync.RWMutex

	tenantId uint32
	tasks    *skiplist.SkipList
}

func NewTenantTasks(tenantId uint32) *TenantTasks {
	return &TenantTasks{
		tenantId: tenantId,
		tasks:    skiplist.New(skiplist.BytesAsc),
	}
}

func (t *TenantTasks) AddTask(task *pb.Task) {
	t.Lock()
	defer t.Unlock()

	t.tasks.Set(task.Range.Start, task)
}

func (t *TenantTasks) RemoveTask(task *pb.Task) {
	t.Lock()
	defer t.Unlock()

	t.tasks.Remove(task.Range.Start)
}

func (t *TenantTasks) GetTask(rangeStart []byte) *pb.Task {
	t.RLock()
	defer t.RUnlock()

	if item := t.tasks.Get(rangeStart); item != nil {
		return item.Value.(*pb.Task)
	}

	return nil
}
