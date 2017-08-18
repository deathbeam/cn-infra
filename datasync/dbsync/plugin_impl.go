// Copyright (c) 2017 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbsync

import (
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/datasync/adapters"
	"github.com/ligato/cn-infra/datasync/persisted/dbsync"
	"github.com/ligato/cn-infra/db/keyval"
	"github.com/ligato/cn-infra/servicelabel"
	"github.com/ligato/cn-infra/utils/safeclose"
)

// PluginID used in the Agent Core flavors
const PluginID core.PluginName = "db-sync"

// Plugin dbsync implements Plugin interface
type Plugin struct {
	Aggregator   *adapters.TransportAggregator // to communicate with transport aggregator
	KvPlugin     *keyval.KvBytesPlugin         // connection type
	ServiceLabel *servicelabel.Plugin
}

// Init uses provided connection to build new transport adapter
func (plugin *Plugin) Init() error {
	etcdConnection := *plugin.KvPlugin
	if etcdConnection != nil {
		broker := etcdConnection.NewBroker(plugin.ServiceLabel.GetAgentPrefix())
		watcher := etcdConnection.NewWatcher(plugin.ServiceLabel.GetAgentPrefix())
		adapter := dbsync.NewAdapter(string(PluginID), broker, watcher)
		plugin.Aggregator.Adapters = append(plugin.Aggregator.Adapters, adapter)
	}

	return nil
}

// Close resources
func (plugin *Plugin) Close() error {
	err := safeclose.Close(plugin.KvPlugin)
	return err
}
