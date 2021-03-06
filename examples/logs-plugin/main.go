package main

import (
	"log"

	"go.ligato.io/cn-infra/v2/agent"
	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/cn-infra/v2/logging/logmanager"
)

// *************************************************************************
// This file contains logger use cases. To define a custom logger, use
// PluginLogger.NewLogger(name). The logger is using 6 levels of logging:
// - Debug
// - Info (this one is default)
// - Warn
// - Error
// - Panic
// - Fatal
//
// Global log levels can be changed locally with the Logger.SetLevel()
// or remotely using REST (but different flavor must be used: rpc.RpcFlavor).
// ************************************************************************/

// PluginName represents name of plugin.
const PluginName = "logs-example"

func main() {
	logging.Info("starting logging example")

	// Prepare example plugin and start the agent
	p := &ExamplePlugin{
		exampleFinished: make(chan struct{}),
		Log:             logging.ForPlugin(PluginName),
		LogManager:      &logmanager.DefaultPlugin,
	}
	a := agent.NewAgent(
		agent.AllPlugins(p),
		agent.QuitOnClose(p.exampleFinished),
	)
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}

// ExamplePlugin presents the PluginLogger API.
type ExamplePlugin struct {
	LogManager *logmanager.Plugin

	Log logging.PluginLogger

	exampleFinished chan struct{}
}

func (plugin *ExamplePlugin) String() string {
	return PluginName
}

// Init demonstrates the usage of PluginLogger API.
func (plugin *ExamplePlugin) Init() (err error) {
	exampleString := "example"
	exampleNum := 15

	// Set log level which logs only entries with current severity or above
	plugin.Log.SetLevel(logging.WarnLevel)  // logs warn/error/fatal/panic levels
	plugin.Log.SetLevel(logging.InfoLevel)  // logs info/warn/error//fatal/panic levels
	plugin.Log.SetLevel(logging.TraceLevel) // logs all levels

	// Basic logger options
	plugin.Log.Print("----------- Log examples -----------")
	plugin.Log.Printf("Print with format specifier. String: %s, Digit: %d, Value: %v", exampleString, exampleNum, plugin)
	plugin.Log.Println()

	// Format also available for all 6 levels of log levels
	plugin.Log.Trace("Trace log message")
	plugin.Log.Debug("Debug log message")
	plugin.Log.Info("Info log message")
	plugin.Log.Infof("Infof log message")
	plugin.Log.Warn("Warn log message")
	plugin.Log.Error("Error log message")
	plugin.showPanicLog()
	//log.Fatal("Bye") calls os.Exit(1) after logging

	// Log with field - automatically adds timestamp
	plugin.Log.WithField("exampleString", exampleString).Info("Info log with field example")
	// For multiple fields
	plugin.Log.WithFields(map[string]interface{}{"exampleString": exampleString, "exampleNum": exampleNum}).Info("Info log with field example string and num")

	// Custom (child) logger with name
	childLogger := plugin.Log.NewLogger("childLogger")
	childLogger.Infof("Log using named logger")
	childLogger.Debug("Debug log using childLogger!")

	childLogger2 := plugin.Log.NewLogger("childLogger2")
	childLogger2.Debug("Debug log using childLogger2!")

	return nil
}

// AfterInit demonstrates the usage of PluginLogger API.
func (plugin *ExamplePlugin) AfterInit() (err error) {
	late := plugin.Log.NewLogger("late")
	late.Debugf("late debug message")

	// End the example
	plugin.Log.Info("logs in plugin example finished, sending shutdown ...")
	close(plugin.exampleFinished)

	return nil
}

// Close implements Plugin interface..
func (plugin *ExamplePlugin) Close() (err error) {
	return nil
}

// showPanicLog demonstrates panic log + recovering.
func (plugin *ExamplePlugin) showPanicLog() {
	defer func() {
		if err := recover(); err != nil {
			plugin.Log.Info("Recovered from panic")
		}
	}()
	plugin.Log.Panic("Panic log: calls panic() after log, will be recovered") //calls panic() after logging
}
