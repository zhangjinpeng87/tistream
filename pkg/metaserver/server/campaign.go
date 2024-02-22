package server

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/zhangjinpeng87/tistream/pkg/metaserver/metadata"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
)

// We use master and standby architecture for the meta-server to achieve HA.
// The master meta-server updates its id to the db every 5 seconds, and the standby
// meta-server will take over the master role if the master meta-server is down.
type Campaign struct {
	sync.RWMutex
	ServerID string

	config *utils.MetaServerConfig

	// Current Master
	currentMaster string
}

func NewCampaign(config *utils.MetaServerConfig) *Campaign {
	serverID := fmt.Sprintf("%s:%d", config.Addr, config.Port)

	return &Campaign{
		ServerID: serverID,
		config:   config,
	}
}

func (c *Campaign) Start(dataManagement *metadata.DataManagement, isMaster *atomic.Bool, eg *errgroup.Group, ctx context.Context) {
	// Start the campaign.
	// If the campaign is successful, set the role to master.
	// If the campaign is failed, stay as standby role.
	eg.Go(func() error {
		ticker := time.NewTimer(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				if isMaster.Load() {
					// Try to update the lease in the db.
					success, err := dataManagement.TryUpdateLease(c.ServerID, c.config.LeaseDuration)
					if err != nil {
						// Log the error.
						return err
					}
					if !success {
						// Update lease failed, transfer to standby role.
						isMaster.Store(false)
						dataManagement.SetReady(false)
					}
				} else {
					// Campaign the master role.
					// If the master lease is expired, the standby meta-server will take over the master role.
					success, err := dataManagement.TryCampaignMaster(c.ServerID, c.config.LeaseDuration)
					if err != nil {
						// Log the error.
						return err
					}
					if success {
						// Set this meta-server as master.
						isMaster.Store(true)

						// Load tasks from db after this meta-server takes over the master role.
						if err := dataManagement.LoadAllTasks(); err != nil {
							// Log the error.
							return err
						}
					} else {
						// Get the current master from the db.
						// This is used to response the client's request to tell the client who is the current master.
						master, err := dataManagement.GetMaster()
						if err != nil {
							// Log the error.
							return err
						}
						if master != c.CurrentMaster() {
							// Update the current master.
							c.Lock()
							c.currentMaster = master
							c.Unlock()
						}
					}
				}
			case <-ctx.Done():
				ticker.Stop()
				return nil
			}
		}
	})
}

func (c *Campaign) CurrentMaster() string {
	c.RLock()
	defer c.Unlock()

	return c.currentMaster
}
