package metadata

import (
	"fmt"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// Backend is the backend of the metadata.
type Backend struct {
	db *utils.DBPool
}

// NewBackend creates a new Backend.
func NewBackend(cfg *utils.MetaServerConfig) (*Backend, error) {
	db, err := utils.NewDBPool(cfg.MysqlHost, cfg.MysqlPort, cfg.MysqlUser, cfg.MysqlPassword, "")
	if err != nil {
		return nil, err
	}
	return &Backend{db: db}, nil
}

// Close closes the backend.
func (b *Backend) Close() {
	b.db.Close()
}

func (b *Backend) BootstrapSchema() error {
	// Create database tistream if not exists.
	_, err := b.db.Exec("CREATE DATABASE IF NOT EXISTS tistream")
	if err != nil {
		return fmt.Errorf("failed to create database tistream: %v", err)
	}

	// Create table tenants if not exists.
	// CREATE TABLE tistream.tenants {
	// 	  id int(10) primary key,
	// 	  cluster_id int(10), // Belongs to which cluster
	// 	  data_change_buffer_addr varchar(255),
	// 	  kms varchar(255),
	// 	  range_start varchar(255),
	// 	  range_end varchar(255),
	// }
	_, err = b.db.Exec("CREATE TABLE IF NOT EXISTS tistream.tenants(" +
		"id INT AUTO_INCREMENT PRIMARY KEY," +
		"cluster_id INT," +
		"data_change_buffer_addr VARCHAR(255)," +
		"kms VARCHAR(255)," +
		"range_start VARCHAR(255)," +
		"range_end VARCHAR(255)," +
		"created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, " +
		"updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)")
	if err != nil {
		return fmt.Errorf("failed to create table tenants: %v", err)
	}

	// Create table tasks if not exists.
	// CREATE TABLE tistream.tasks {
	// 	  id int(10) primary key, // Task id
	// 	  tenant_id int(10), // Task belongs to which tenant. Typically there is just one task for each tenant, there might be multiple tasks if the tenant is large.
	// 	  range_start varchar(255), // Start key of this continuous range
	// 	  range_end varchar(255),
	// 	  uuid varchar(255), // UUID of this task, the UUID will be used as the directory name in the storage
	// 	  sorter_addr varchar(255), Which sorter is responsible for this task
	// 	  snapshot_addr varchar(255), // Where to store the task snapshot
	//    unique (tenant_id, range_start), // There should be only one task for each tenant and range
	// }
	_, err = b.db.Exec("CREATE TABLE IF NOT EXISTS tistream.tasks(" +
		"id INT AUTO_INCREMENT PRIMARY KEY," +
		"tenant_id INT," +
		"range_start VARCHAR(255)," +
		"range_end VARCHAR(255)," +
		"uuid VARCHAR(255), " +
		"sorter_addr VARCHAR(255)," +
		"snapshot_addr VARCHAR(255)," +
		"created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, " +
		"updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)," +
		"UNIQUE (tenant_id, range_start)")
	if err != nil {
		return fmt.Errorf("failed to create table tasks: %v", err)
	}

	// Create table owner if not exists.
	// CREATE TABLE tistream.owner {
	// 	  id int(10) primary key, // There is only one row in this table, its id is 1. If there are multiple TiStream clusters in one region, then there should be multiple rows.
	// 	  master int(10), // Who is the master
	// 	  lease timestamp, // If the master doesn't update lease for 10s, other nodes can take over the master role
	// }
	_, err = b.db.Exec("CREATE TABLE IF NOT EXISTS tistream.owner(" +
		"id INT AUTO_INCREMENT PRIMARY KEY," +
		"master VARCHAR(255)," +
		"lease TIMESTAMP," +
		"created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, " +
		"updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)")
	if err != nil {
		return fmt.Errorf("failed to create table owner: %v", err)
	}

	// Init row 1 in table owner if not exists.
	_, err = b.db.Exec("INSERT IGNORE INTO tistream.owner (id, master, lease) VALUES (1, \"dummy\", NOW())")
	if err != nil {
		return fmt.Errorf("failed to init row 1 in table owner: %v", err)
	}

	return nil
}

func (b *Backend) GetMaster() (string, error) {
	var master string
	err := b.db.QueryRow("SELECT master FROM tistream.owner WHERE id = 1").Scan(&master)
	if err != nil {
		return "", fmt.Errorf("failed to get owner: %v", err)
	}
	return master, nil
}

func (b *Backend) TryUpdateLease(who string, leaseDurSec int) (bool, error) {
	res, err := b.db.Exec("UPDATE tistream.owner SET lease = NOW() + INTERVAL ? SECOND WHERE id = 1 AND master = ?",
		leaseDurSec, who)
	if err != nil {
		return false, err
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if affectedRows == 1 {
		return true, nil
	}
	return false, nil
}

func (b *Backend) TryCampaignMaster(who string, leaseDurSec int) (bool, error) {
	res, err := b.db.Exec("UPDATE tistream.owner SET master = ?, lease = NOW() + INTERVAL ? SECOND WHERE id = 1 AND lease < NOW()",
		who, leaseDurSec)
	if err != nil {
		return false, err
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if affectedRows == 1 {
		return true, nil
	}
	return false, nil
}

// Add a new task to backend.
// Return error if the task already exists. TenantID and RangeStart are unique to indicate a task.
func (b *Backend) AddTask(task *pb.Task) error {
	res, err := b.db.Exec("INSERT INTO tistream.tasks (tenant_id, range_start, range_end, sorter, snapshot_addr) VALUES (?, ?, ?, ?, ?, ?)",
		task.TenantId, string(task.Range.Start), string(task.Range.End), task.SorterAddr, task.SnapAddr)
	if err != nil {
		return fmt.Errorf("failed to insert task: %v", err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %v", err)
	}
	return nil
}

// Remove a task from backend.
// Return true if the task is removed, false if the task doesn't exist.
func (b *Backend) RemoveTask(tenantID uint32, rangeStart []byte) (bool, error) {
	res, err := b.db.Exec("DELETE FROM tistream.tasks WHERE tenant_id = ? AND range_start = ?", tenantID, string(rangeStart))
	if err != nil {
		return false, fmt.Errorf("failed to remove task: %v", err)
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get affected rows: %v", err)
	}
	if affectedRows == 1 {
		return true, nil
	}

	return false, nil
}

// Modify a task in backend.
// Return true if the task is modified, false if the task doesn't exist.
func (b *Backend) ModifyTask(task *pb.Task) (bool, error) {
	res, err := b.db.Exec("UPDATE tistream.tasks SET sorter = ?, snapshot_addr = ?, range_end = ? WHERE tenant_id = ? AND range_start = ?",
		task.SorterAddr, task.SnapAddr, task.Range.End, task.TenantId, string(task.Range.Start))
	if err != nil {
		return false, fmt.Errorf("failed to modify task: %v", err)
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get affected rows: %v", err)
	}
	if affectedRows == 1 {
		return true, nil
	}

	return false, nil
}

// LoadAllTasks loads all tasks from backend.
func (b *Backend) LoadAllTasks(f func(t uint64, task *pb.Task)) error {
	rows, err := b.db.Query("SELECT id, tenant_id, range_start, range_end, sorter_addr, snapshot_addr FROM tistream.tasks")
	if err != nil {
		return fmt.Errorf("failed to query tasks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, tenantID uint64
		var rangeStart, rangeEnd, snapshotAddr, sorterAddr string
		if err := rows.Scan(&id, &tenantID, &rangeStart, &rangeEnd, &sorterAddr, &snapshotAddr); err != nil {
			return fmt.Errorf("failed to scan tasks: %v", err)
		}

		pbTask := &pb.Task{
			TenantId:   tenantID,
			Range:      &pb.Task_Range{Start: []byte(rangeStart), End: []byte(rangeEnd)},
			SorterAddr: sorterAddr,
			SnapAddr:   snapshotAddr,
		}

		f(tenantID, pbTask)
	}

	return nil
}
