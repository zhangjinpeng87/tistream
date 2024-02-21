package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
)

const (
	// Master is the master role of the meta-server.
	Master = "master"
	// Standby is the standby role of the meta-server.
	Standby = "standby"
)

// We use master and standby architecture for the meta-server to achieve HA.
// The master meta-server updates its id to the db every 5 seconds, and the standby
// meta-server will take over the master role if the master meta-server is down.
type Campaign struct {
	sync.Mutex

	ServerID string
	Role string
	dbPool *utils.DBPool

	done chan struct{}
}

func NewCampaign(addr string, port int, dbPool *utils.DBPool) *Campaign {
	serverID := fmt.Sprintf("%s:%d", addr, port)

	return &Campaign{
		ServerID: serverID,
		Role: Standby,
		dbPool: dbPool,
		done: make(chan struct{}),
	}
}

func (c *Campaign) Start() error {
	// Start the campaign.
	// If the campaign is successful, set the role to master.
	// If the campaign is failed, stay as standby role.

	go func() {
		ticker := time.NewTimer(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				if c.Role == Master {
					// Try to update the lease in the db.
					if !c.TryUpdateLease() {
						// Campaign failed, stay as standby role.
						c.Lock()
						c.Role = Standby
						c.Unlock()
					}
				} else {
					// Campaign the master role.
					// If the master lease is expired, the standby meta-server will take over the master role.
				}
			case <-c.done:
				return
		}
	}()

	return nil
}

func (c *Campaign) TryUpdateLease() bool {
	if c.Role != Master {
		return false
	}

	// Update the lease in the db.
	res, err := c.dbPool.Exec("UPDATE tistream.owner SET lease = ? WHERE id = 1 and owner = ?", time.Now()+10*time.Second, c.ServerID)
	if err != nil {
		return false
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return false
	}
	if affectedRows == 1 {
		return true
	}
	return false
}

func (c *Campaign) IsMaster() bool {
	c.Lock()
	defer c.Unlock()

	return c.Role == Master
}

func (c *Campaign) Stop() error {
	// Stop the campaign.
	c.done <- struct{}{}
	return nil
}