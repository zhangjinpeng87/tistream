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
	Role     string

	config *utils.MetaServerConfig
	dbPool *utils.DBPool
	done   chan struct{}
	wg     sync.WaitGroup
}

func NewCampaign(addr string, port int, dbPool *utils.DBPool, config *utils.MetaServerConfig) *Campaign {
	serverID := fmt.Sprintf("%s:%d", addr, port)

	return &Campaign{
		ServerID: serverID,
		Role:     Standby,
		config:   config,
		dbPool:   dbPool,
		done:     make(chan struct{}),
	}
}

func (c *Campaign) Start() error {
	// Start the campaign.
	// If the campaign is successful, set the role to master.
	// If the campaign is failed, stay as standby role.
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		ticker := time.NewTimer(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				if c.IsMaster() {
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
					if c.TryCampaignMaster() {
						c.Lock()
						c.Role = Master
						c.Unlock()
					}
				}
			case <-c.done:
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func (c *Campaign) TryUpdateLease() bool {
	if c.Role != Master {
		return false
	}

	// Only update lease when the `master` colunn value is this server itself.
	res, err := c.dbPool.Exec("UPDATE tistream.owner SET lease = NOW() + INTERVAL ? SECOND WHERE id = 1 and master = ?",
		c.config.LeaseDuration, c.ServerID)
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

func (c *Campaign) TryCampaignMaster() bool {
	if c.IsMaster() {
		return true
	}

	// Campaign the master role.
	// If the master lease is expired, the standby meta-server will take over the master role.
	res, err := c.dbPool.DB.Exec("UPDATE tistream.owner SET master = ?, lease = NOW() + INTERVAL ? SECOND WHERE id = 1 AND lease < NOW()",
		c.ServerID, c.config.LeaseDuration)
	if err != nil {
		return false
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return false
	}
	if affectedRows == 1 {
		// The standby meta-server takes over the master role.
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
	c.wg.Wait()
	return nil
}
