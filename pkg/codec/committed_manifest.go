package codec

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/zhangjinpeng87/tistream/pkg/storage"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"google.golang.org/protobuf/proto"
)

const (
	CommittedManifestName    string = "manifest"
	CommittedManifestMagic   string = "TISTREAM-COMMITTED-MANIFEST"
	CommittedManifestVersion uint32 = 1
)

type CommittedManifest struct {
	sync.RWMutex
	// The tenant id.
	TenantID uint64

	// Root directory of the committed data pool for the tenant.
	// In above file org it should be "{prefix}committed_data_pool/Tenant-{1}/"
	RootDir string

	// Backend Storage
	BackendStorage storage.ExternalStorage

	// Ranges Snapshots in different ts.
	// ts -> *pb.RangesSnapshot
	Snapshots []*pb.RangesSnapshot
	latestTs  uint64
}

// NewCommittedDataPoolManifest creates a new CommittedDataPoolManifest.
func NewCommittedDataPoolManifest(tenantID uint64, rootDir string, backendStorage storage.ExternalStorage) *CommittedManifest {
	return &CommittedManifest{
		TenantID:       tenantID,
		RootDir:        rootDir,
		BackendStorage: backendStorage,
		Snapshots:      make(map[uint64]*pb.RangesSnapshot),
	}
}

// AddSnapshot adds a snapshot to the committed manifest.
func (m *CommittedManifest) AddSnapshot(ts uint64, snapshot *pb.RangesSnapshot) error {
	m.Lock()
	defer m.Unlock()

	if ts <= m.latestTs {
		return fmt.Errorf("invalid snapshto with ts %d, it should be larger than the latest ts %d", ts, m.latestTs)
	}

	m.Snapshots = append(m.Snapshots, snapshot)

	return nil
}

// Load from file.
func (m *CommittedManifest) Load() error {
	m.Lock()
	defer m.Unlock()

	// Load the snapshots from the backend storage.
	fname := m.RootDir + CommittedManifestName
	data, err := m.BackendStorage.GetFile(fname)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}

	// Verify checksum
	checksum := binary.LittleEndian.Uint32(data[len(data)-4:])
	if checksum != CalcChecksum(data[:len(data)-4]) {
		return utils.ErrChecksumNotMatch
	}

	// Decode the snapshots.
	r := bytes.NewReader(data[:len(data)-4])
	decoder := NewCommittedManifestDecoder(m.TenantID)
	snapshots, err := decoder.Decode(r)
	if err != nil {
		return err
	}
	m.Snapshots = snapshots
	for _, s := range snapshots {
		if s.Ts > m.latestTs {
			m.latestTs = s.Ts
		}
	}

	return nil
}

// Save to file.
func (m *CommittedManifest) Save() error {
	m.RLock()
	defer m.RUnlock()

	var data bytes.Buffer
	w := bufio.NewWriter(&data)
	encoder := NewCommittedManifestEncoder(m.TenantID)
	if err := encoder.Encode(w, m.Snapshots); err != nil {
		return err
	}

	// Flush the buffer.
	if err := w.Flush(); err != nil {
		return err
	}

	// Calculate the checksum.
	checksum := CalcChecksum(data.Bytes())
	if err := binary.Write(w, binary.LittleEndian, checksum); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	// Save the snapshots to the backend storage.
	fname := m.RootDir + CommittedManifestName
	return m.BackendStorage.PutFile(fname, data.Bytes())
}

// Compact the committed manifest.
// Return true if the manifest is compacted, otherwise return false.
func (m *CommittedManifest) CompactTo(ts uint64) bool {
	m.Lock()
	defer m.Unlock()

	if len(m.Snapshots) == 0 {
		return false
	}

	pos := sort.Search(len(m.Snapshots), func(i int) bool {
		return m.Snapshots[i].Ts > ts
	})
	if pos <= 1 {
		// No need to compact.
		// Pos == 0 means all ts are larger than ts.
		// Pos == 1 allTs[0] is the first snapshot less or eqaul to ts, we should keep it.
		return false
	}
	m.Snapshots = m.Snapshots[pos-1:]

	// todo: remove the old files.

	return true
}

type CommittedManifestEncoder struct {
	// The tenant id.
	TenantID uint64
}

func NewCommittedManifestEncoder(tenantID uint64) *CommittedManifestEncoder {
	return &CommittedManifestEncoder{
		TenantID: tenantID,
	}
}

func (e *CommittedManifestEncoder) Encode(w io.Writer, snapshots []*pb.RangesSnapshot) error {
	// Encode the header.
	// Write the magic header
	if _, err := w.Write([]byte(CommittedManifestMagic)); err != nil {
		return err
	}
	// Write the file version
	if err := binary.Write(w, binary.LittleEndian, CommittedDataFileVersion); err != nil {
		return err
	}
	// Write the tenant id
	if err := binary.Write(w, binary.LittleEndian, e.TenantID); err != nil {
		return err
	}
	// Write the number of snapshots
	if err := binary.Write(w, binary.LittleEndian, uint32(len(snapshots))); err != nil {
		return err
	}
	for _, snapshot := range snapshots {
		// Marsharl the snapshot
		data, err := proto.Marshal(snapshot)
		if err != nil {
			return err
		}

		// Write the data len
		if err := binary.Write(w, binary.LittleEndian, uint32(len(data))); err != nil {
			return err
		}

		// Write the data
		if _, err := w.Write(data); err != nil {
			return err
		}
	}

	return nil
}

type CommittedManifestDecoder struct {
	// The tenant id.
	TenantID uint64
}

func NewCommittedManifestDecoder(tenantID uint64) *CommittedManifestDecoder {
	return &CommittedManifestDecoder{
		TenantID: tenantID,
	}
}

func (d *CommittedManifestDecoder) Decode(r io.Reader) ([]*pb.RangesSnapshot, error) {
	// Decode the header.
	// Read the magic header
	magic := make([]byte, len(CommittedManifestMagic))
	if _, err := io.ReadFull(r, magic); err != nil {
		return nil, err
	}

	// Read the file version
	var version uint32
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return nil, err
	}

	switch version {
	case 1:
		return d.DecodeV1(r)
	default:
		return nil, utils.ErrInvalidCommittedManifestFile
	}
}

func (d *CommittedManifestDecoder) DecodeV1(r io.Reader) ([]*pb.RangesSnapshot, error) {
	// Read the tenant id
	var tenantID uint64
	if err := binary.Read(r, binary.LittleEndian, &tenantID); err != nil {
		return nil, err
	}

	// Read the number of snapshots
	var snapshotCnt uint32
	if err := binary.Read(r, binary.LittleEndian, &snapshotCnt); err != nil {
		return nil, err
	}

	// Read the snapshots
	snapshots := make([]*pb.RangesSnapshot, 0, snapshotCnt)
	for i := uint32(0); i < snapshotCnt; i++ {
		// Read the data len
		var dataLen uint32
		if err := binary.Read(r, binary.LittleEndian, &dataLen); err != nil {
			return nil, err
		}

		// Read the data
		data := make([]byte, dataLen)
		if _, err := io.ReadFull(r, data); err != nil {
			return nil, err
		}

		// Unmarshal the data
		var snapshot pb.RangesSnapshot
		if err := proto.Unmarshal(data, &snapshot); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, &snapshot)
	}

	return snapshots, nil
}
