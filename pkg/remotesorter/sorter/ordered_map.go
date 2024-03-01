package sorter

import (
	"github.com/huandu/skiplist"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// OrderedEventMap is a map with ordered keys.
type OrderedEventMap interface {
	// Put puts the key-value pair to the map.
	Put(key string, value *pb.EventRow)

	// Get gets the value by the key.
	Get(key string) *pb.EventRow

	// Delete deletes the key-value pair by the key.
	Delete(key string) *pb.EventRow

	// GetRange gets the key-value pairs by the range.
	GetRange(start, end string) []*pb.EventRow

	// DelRange deletes the key-value pairs by the range.
	DelRange(start, end string)

	// Reset resets the map.
	Reset()
}

// SkipListEventMap is the skip list based ordered event map.
type SkipListEventMap struct {
	list *skiplist.SkipList
}

// NewSkipListEventMap creates a new SkipListEventMap.
func NewSkipListEventMap() *SkipListEventMap {
	return &SkipListEventMap{
		list: skiplist.New(skiplist.StringAsc),
	}
}

// Put puts the key-value pair to the map.
func (m *SkipListEventMap) Put(key string, value *pb.EventRow) {
	m.list.Set(key, value)
}

// Get gets the value by the key.
func (m *SkipListEventMap) Get(key string) *pb.EventRow {
	v := m.list.Get(key)
	if v == nil {
		return nil
	}
	return v.Value.(*pb.EventRow)
}

// Delete deletes the key-value pair by the key.
func (m *SkipListEventMap) Delete(key string) *pb.EventRow {
	e := m.list.Remove(key)
	if e == nil {
		return nil
	}
	return e.Value.(*pb.EventRow)
}

// GetRange gets the key-value pairs by the range.
func (m *SkipListEventMap) GetRange(start, end string) []*pb.EventRow {
	var res []*pb.EventRow
	e := m.list.Find(start)
	for e != nil {
		if e.Key().(string) >= end {
			break
		}
		res = append(res, e.Value.(*pb.EventRow))
		e = e.Next()
	}
	return res
}

// DelRange deletes the key-value pairs by the range.
func (m *SkipListEventMap) DelRange(start, end string) {
	e := m.list.Find(start)
	for e != nil {
		if e.Key().(string) >= end {
			break
		}
		m.list.Remove(e.Key().(string))
		e = e.Next()
	}
}

// Reset resets the map.
func (m *SkipListEventMap) Reset() {
	m.list = skiplist.New(skiplist.StringAsc)
}

// LSMTreeEventMap is the LSM tree based ordered event map.
type LSMTreeEventMap struct {
	// Todo: implement the LSM tree based ordered event map.
}
