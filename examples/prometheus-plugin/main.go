package main

import (
	"fmt"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/flavors/local"
	prom "github.com/ligato/cn-infra/rpc/prometheus"
	"github.com/ligato/cn-infra/rpc/rest"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"time"
)

// *************************************************************************
// This example demonstrates the usage of prometheus plugin that allows
// to expose metrics.
//
// Default metrics are exposed at path /metrics on port specified for http plugin
// to access these metrics following command can be used:
//       curl localhost:9191/metrics
//
// There is also created custom metrics registry exposed accessible:
//       curl localhost:9191/custom
// ************************************************************************/

func main() {
	// Init close channel used to stop the example.
	exampleFinished := make(chan struct{}, 1)

	// Start Agent with ExamplePlugin, REST, prometheus plugin & FlavorLocal (reused cn-infra plugins).
	agent := local.NewAgent(local.WithPlugins(func(flavor *local.FlavorLocal) []*core.NamedPlugin {
		httpPlugin := &rest.Plugin{}

		rest.DeclareHTTPPortFlag("http")
		httpPlugDeps := *flavor.InfraDeps("http", local.WithConf())
		httpPlugin.Deps.Log = httpPlugDeps.Log
		httpPlugin.Deps.PluginConfig = httpPlugDeps.PluginConfig
		httpPlugin.Deps.PluginName = httpPlugDeps.PluginName

		prometheusPlugin := &prom.Plugin{}
		prometheusPlugin.Deps.PluginInfraDeps = *flavor.InfraDeps("prometheus")
		prometheusPlugin.Deps.HTTP = httpPlugin

		examplePlug := &ExamplePlugin{}
		examplePlug.Deps.PluginInfraDeps = *flavor.InfraDeps("prometheus-example")
		examplePlug.Deps.Prometheus = prometheusPlugin

		return []*core.NamedPlugin{
			{httpPlugin.PluginName, httpPlugin},
			{prometheusPlugin.PluginName, prometheusPlugin},
			{examplePlug.PluginName, examplePlug}}
	}))
	core.EventLoopWithInterrupt(agent, exampleFinished)
}

// Deps group dependencies of the ExamplePlugin
type Deps struct {
	local.PluginInfraDeps
	Prometheus prom.API
}

// ExamplePlugin demonstrates the usage of datasync API.
type ExamplePlugin struct {
	Deps

	temporaryCounter prometheus.Gauge
	counterVal       int

	gaugeVec *prometheus.GaugeVec
}

// Identifier of the custom registry
const customRegistry = "/custom"

const orderLabel = "order"

// Init creates metric registries and adds gauges
func (plugin *ExamplePlugin) Init() error {

	// add new metric to default registry (accessible at the path /metrics)
	//
	// the current value is returned by provided callback
	// created gauge is identified by tuple(namespace, subsystem, name) only the name field is mandatory
	// additional properties can be defined using labels - key-value pairs. They do not change over time for the given gauge.
	err := plugin.Prometheus.RegisterGaugeFunc(prom.DefaultRegistry, "ns", "sub", "gaugeOne",
		"this metrics represents randomly generated numbers", prometheus.Labels{"Property1": "ABC", "Property2": "DEF"}, func() float64 {
			return rand.Float64()
		})
	if err != nil {
		return err
	}

	// create new registry that will be exposed at /custom path
	err = plugin.Prometheus.NewRegistry(customRegistry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError})
	if err != nil {
		return err
	}

	// create gauge using prometheus API
	plugin.temporaryCounter = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Countdown",
		Help: "This gauge is decremented by 1 each second, once it reaches 0 the gauge is removed.",
	})
	plugin.counterVal = 60
	plugin.temporaryCounter.Set(float64(plugin.counterVal))

	// register created gauge to the custom registry
	err = plugin.Prometheus.Register(customRegistry, plugin.temporaryCounter)
	if err != nil {
		return err
	}

	// create gauge vector and register it
	plugin.gaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "Vector",
		Help:        "This gauge groups multiple similar metrics.",
		ConstLabels: prometheus.Labels{"type": "vector", "answer": "42"},
	}, []string{orderLabel})
	err = plugin.Prometheus.Register(customRegistry, plugin.gaugeVec)

	return err

}

// AfterInit starts go routines that modifies metrics
func (plugin *ExamplePlugin) AfterInit() error {

	go plugin.decrementCounter()

	go plugin.addNewGaugesToVector()

	return nil
}

// Close cleanup resources allocated by plugin
func (plugin *ExamplePlugin) Close() error {
	return nil
}

func (plugin *ExamplePlugin) addNewGaugesToVector() {
	for i := 1; i < 10; i++ {
		// add gauge with given labels to the vector
		g, err := plugin.gaugeVec.GetMetricWith(prometheus.Labels{orderLabel: fmt.Sprint(i)})
		if err != nil {
			plugin.Log.Error(err)
		} else {
			g.Set(1)
		}
		time.Sleep(2 * time.Second)
	}

}

func (plugin *ExamplePlugin) decrementCounter() {
	for {
		select {
		case <-time.After(time.Second):
			if plugin.counterVal == 0 {
				// once the countdown reaches zero remove gauge from registry+
				plugin.Prometheus.Unregister(customRegistry, plugin.temporaryCounter)
				return
			}
			plugin.counterVal--
			plugin.temporaryCounter.Set(float64(plugin.counterVal))

		}
	}
}