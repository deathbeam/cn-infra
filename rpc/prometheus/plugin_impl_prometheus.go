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

package prometheus

import (
	"errors"
	"github.com/ligato/cn-infra/flavors/local"
	"github.com/ligato/cn-infra/rpc/rest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/unrolled/render"
	"net/http"
	"strings"
	"sync"
)

// DefaultMetricsPath default Prometheus metrics URL
const DefaultRegistry = "/metrics"

var (
	// PathInvalidFormatError is returned if the path doesn't start with slash
	PathInvalidFormatError = errors.New("path is invalid, it must start with '/' character")
	// PathAlreadyRegistryError is returned on attempt to register a path used by a registry
	PathAlreadyRegistryError = errors.New("registry with the path is already registered")
	// RegistryNotFoundError is returned on attempt to use register that has not been created
	RegistryNotFoundError = errors.New("registry was not found")
)

// Plugin struct holds all plugin-related data.
type Plugin struct {
	Deps
	sync.Mutex
	// regs is a map of URL path(symbolic names) to registries. Registries group metrics and can be exposed at different urls.
	regs map[string]*registry
}

// Deps lists dependencies of the plugin.
type Deps struct {
	local.PluginInfraDeps // inject
	// HTTP server used to expose metrics
	HTTP rest.HTTPHandlers // inject
}

type registry struct {
	prometheus.Gatherer
	prometheus.Registerer
	// httpOpts applied when exposing registry using http
	httpOpts promhttp.HandlerOpts
}

// Init initializes the internal structures
func (p *Plugin) Init() (err error) {

	p.regs = map[string]*registry{}

	// add default registry
	p.regs[DefaultRegistry] = &registry{
		Gatherer:   prometheus.DefaultGatherer,
		Registerer: prometheus.DefaultRegisterer,
	}

	return nil
}

// AfterInit registers HTTP handlers.
func (p *Plugin) AfterInit() error {
	if p.HTTP != nil {
		p.Lock()
		defer p.Unlock()
		for path, reg := range p.regs {
			p.HTTP.RegisterHTTPHandler(path, p.createHandlerHandler(reg.Gatherer), "GET")
			p.Log.Infof("Serving %s on port %d", path, p.HTTP.GetPort())

		}
	} else {
		p.Log.Info("Unable to register Prometheus metrics handlers, HTTP is nil")
	}

	return nil
}

// Close cleans up the allocated resources.
func (p *Plugin) Close() error {
	return nil
}

// NewRegistry creates new registry exposed at defined URL path (must begin with '/' character), path is used to reference
// registry while adding new metrics into registry, opts adjust the behavior of exposed registry. Must be called before
// AfterInit phase of the Prometheus plugin. An attempt to create  a registry with path that is already used
// by different registry returns an error.
func (p *Plugin) NewRegistry(path string, opts promhttp.HandlerOpts) error {
	p.Lock()
	defer p.Unlock()

	if !strings.HasPrefix(path, "/") {
		return PathInvalidFormatError
	}
	if _, found := p.regs[path]; found {
		return PathAlreadyRegistryError
	}
	newReg := prometheus.NewRegistry()
	p.regs[path] = &registry{
		Registerer: newReg,
		Gatherer:   newReg,
		httpOpts:   opts,
	}
	return nil
}

// Register registers prometheus metric to a specified registry. In order to add metrics
// to default registry use prometheus.DefaultRegistry const.
func (p *Plugin) Register(registryPath string, collector prometheus.Collector) error {
	p.Lock()
	defer p.Unlock()

	reg, found := p.regs[registryPath]
	if !found {
		return RegistryNotFoundError
	}
	return reg.Register(collector)
}

// RegisterGauge registers custom gauge with specific valueFunc to report status when invoked. RegistryPath identifies
// the registry where gauge is added.
func (p *Plugin) RegisterGaugeFunc(registryPath string, namespace string, subsystem string, name string, help string,
	labels prometheus.Labels, valueFunc func() float64) error {

	p.Lock()
	defer p.Unlock()

	reg, found := p.regs[registryPath]
	if !found {
		return RegistryNotFoundError
	}

	gaugeName := name
	if subsystem != "" {
		gaugeName = subsystem + "_" + gaugeName
	}
	if namespace != "" {
		gaugeName = namespace + "_" + gaugeName
	}

	err := reg.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Subsystem:   subsystem,
			Name:        name,
			Help:        help,
			ConstLabels: labels,
		},
		valueFunc,
	))
	if err != nil {
		p.Log.Errorf("GaugeFunc('%s') registration failed: %s", gaugeName, err)
		return err
	}
	p.Log.Infof("GaugeFunc('%s') registered.", gaugeName)
	return nil
}

func (p *Plugin) createHandlerHandler(gatherer prometheus.Gatherer) func(formatter *render.Render) http.HandlerFunc {
	return func(formatter *render.Render) http.HandlerFunc {
		return promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}).ServeHTTP
	}
}
