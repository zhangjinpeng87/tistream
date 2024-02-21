package metadata

import (
	"fmt"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// MetaBootstrapper is used to bootstrap the schema, load metadata from the db and init the meta-sever.
type Bootstrapper struct {
	dbPool *utils.DBPool
}

// NewBootstrapper creates a new Bootstrapper.
func NewBootstrapper(dbPool *utils.DBPool) *Bootstrapper {
	return &Bootstrapper{
		dbPool: dbPool,
	}
}

func (b *Bootstrapper) Bootstrap(taskManagement *TaskManagement) error {
	if err := b.initSchema(); err != nil {
		return err
	}

	if err := b.loadTasks(taskManagement); err != nil {
		return err
	}

	return nil
}

// InitSchema initializes the schema.
func (b *Bootstrapper) initSchema() error {
	// Create database tistream if not exists.
	_, err := b.dbPool.Exec("CREATE DATABASE IF NOT EXISTS tistream")
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
	_, err = b.dbPool.Exec("CREATE TABLE IF NOT EXISTS tistream.tenants(" +
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
	// 	  dispatcher int(10), // Which dispatcher is responsible for this task
	// 	  sorter int(10), Which sorter is responsible for this task
	// 	  snapshot_addr varchar(255), // Where to store the task snapshot
	// }
	_, err = b.dbPool.Exec("CREATE TABLE IF NOT EXISTS tistream.tasks(" +
		"id INT AUTO_INCREMENT PRIMARY KEY," +
		"tenant_id INT," +
		"range_start VARCHAR(255)," +
		"range_end VARCHAR(255)," +
		"dispatcher INT," +
		"sorter INT," +
		"snapshot_addr VARCHAR(255)," +
		"created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, " +
		"updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)")
	if err != nil {
		return fmt.Errorf("failed to create table tasks: %v", err)
	}

	// Create table owner if not exists.
	// CREATE TABLE tistream.owner {
	// 	  id int(10) primary key, // There is only one row in this table, its id is 1. If there are multiple TiStream clusters in one region, then there should be multiple rows.
	// 	  owner int(10), // Who is the owner
	// 	  last_heart_beat timestamp, // If the owner doesn't update heart beat for 10s, other nodes reenable can be the owner
	// }
	_, err = b.dbPool.Exec("CREATE TABLE IF NOT EXISTS tistream.owner(" +
		"id INT AUTO_INCREMENT PRIMARY KEY," +
		"owner VARCHAR(255)," +
		"lease TIMESTAMP," +
		"created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, " +
		"updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)")
	if err != nil {
		return fmt.Errorf("failed to create table owner: %v", err)
	}

	return nil
}

// LoadTasks load tasks from the db.
func (b *Bootstrapper) loadTasks(taskManagement *TaskManagement) error {
	rows, err := b.dbPool.Query("SELECT * FROM tistream.tasks")
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

		taskManagement.AddTask(tenantID, task)
	}

	return nil
}
