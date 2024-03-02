package sorter

import (
	"bytes"

	"github.com/huandu/skiplist"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// OrderedEventMap is a map with ordered keys.
type OrderedEventMap interface {
	// Put puts the key-value pair to the map.
	Put(key []byte, value *pb.EventRow)

	// Get gets the value by the key.
	Get(key []byte) *pb.EventRow

	// Delete deletes the key-value pair by the key.
	Delete(key []byte) *pb.EventRow

	// GetRange gets the key-value pairs by the range.
	GetRange(start, end []byte) []*pb.EventRow

	// DelRange deletes the key-value pairs by the range.
	DelRange(start, end []byte)

	// IterateRange iterates the key-value pairs by the range.
	IterateRange(start, end []byte, f func(key []byte, value *pb.EventRow) bool)

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
		list: skiplist.New(skiplist.BytesAsc),
	}
}

// Put puts the key-value pair to the map.
func (m *SkipListEventMap) Put(key []byte, value *pb.EventRow) {
	m.list.Set(key, value)
}

// Get gets the value by the key.
func (m *SkipListEventMap) Get(key []byte) *pb.EventRow {
	v := m.list.Get(key)
	if v == nil {
		return nil
	}
	return v.Value.(*pb.EventRow)
}

// Delete deletes the key-value pair by the key.
func (m *SkipListEventMap) Delete(key []byte) *pb.EventRow {
	e := m.list.Remove(key)
	if e == nil {
		return nil
	}
	return e.Value.(*pb.EventRow)
}

// GetRange gets the key-value pairs by the range.
func (m *SkipListEventMap) GetRange(start, end []byte) []*pb.EventRow {
	var res []*pb.EventRow
	e := m.list.Find(start)
	for e != nil {
		if bytes.Compare(e.Key().([]byte), end) >= 0 {
			break
		}
		res = append(res, e.Value.(*pb.EventRow))
		e = e.Next()
	}
	return res
}

// DelRange deletes the key-value pairs by the range.
func (m *SkipListEventMap) DelRange(start, end []byte) {
	e := m.list.Find(start)
	for e != nil {
		if bytes.Compare(e.Key().([]byte), end) >= 0 {
			break
		}
		m.list.Remove(e.Key().([]byte))
		e = e.Next()
	}
}

// IterateRange iterates the key-value pairs by the range.
func (m *SkipListEventMap) IterateRange(start, end []byte, f func(key []byte, value *pb.EventRow) bool) {
	e := m.list.Find(start)
	for e != nil {
		if bytes.Compare(e.Key().([]byte), end) >= 0 {
			break
		}
		if !f(e.Key().([]byte), e.Value.(*pb.EventRow)) {
			break
		}
		e = e.Next()
	}
}

// Reset resets the map.
func (m *SkipListEventMap) Reset() {
	m.list = skiplist.New(skiplist.BytesAsc)
}

type DiskBasedOrderMap struct {
	// Todo: implement the disk based ordered event map.
}

type SkiplistFactory struct{}

func (f *SkiplistFactory) Create() OrderedEventMap {
	return NewSkipListEventMap()
}
