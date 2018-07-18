package main

import (
	"errors"
	"log"

	"github.com/ligato/cn-infra/agent"
	"github.com/ligato/cn-infra/logging"
	"github.com/ligato/cn-infra/logging/logrus"
	"github.com/ligato/cn-infra/rpc/grpc"
	"github.com/ligato/cn-infra/rpc/rest"
	"golang.org/x/net/context"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

// *************************************************************************
// This file contains GRPC service exposure example. To register service use
// Server.RegisterService(descriptor, service)
// ************************************************************************/

const PluginName = "myPlugin"

func main() {
	// --------------------
	// ALL DEFAULT
	// --------------------

	/*p := &ExamplePlugin{
		Deps: Deps{
			PluginName: PluginName,
			Log:        logging.ForPlugin(PluginName),
			GRPC:       &grpc.DefaultPlugin,
		},
	}*/

	// --------------------
	// CUSTOM INSTANCE
	// --------------------

	p := &ExamplePlugin{
		GRPC: grpc.NewPlugin(
			//grpc.UseName("myGRPC"),
			grpc.UseHTTP(&rest.DefaultPlugin),
			grpc.UseDeps(func(deps *grpc.Deps) {
				deps.HTTP = &rest.DefaultPlugin //rest.NewPlugin()
				//deps.PluginName = core.PluginName("myGRPC")
				deps.SetName("myGRPC")
			}),
		),
		//GRPC: &grpc.DefaultPlugin,
		Log: logging.ForPlugin(PluginName),
	}

	// --------------------
	// CHANGE GLOBAL DEFAULT
	// --------------------

	/*rest.DefaultPlugin = *rest.NewPlugin(
		rest.UseConf(rest.Config{
			Endpoint: ":1234",
		}),
	)
	p := &ExamplePlugin{
		Deps: Deps{
			PluginName: PluginName,
			Log:        logging.ForPlugin(PluginName),
			GRPC:       grpc.DefaultPlugin,
		},
	}*/

	// --------------------
	// DISABLE DEP
	// --------------------

	/*myGRPC := grpc.NewPlugin(
		grpc.UseDeps(grpc.Deps{
			HTTP: rest.Disabled,
		}),
	)

	//rest.DefaultPlugin = rest.NewPlugin(rest.UseDisabled())*/

	// --------------------
	// INIT AGENT
	// --------------------

	/*myGRPC := grpc.NewPlugin(
		//grpc.UseCustom(grpc.PluginDeps{}),
		//grpc.UseDefaults(),
		grpc.UseDeps(grpc.Deps{
			//Log: logging.ForPlugin("myGRPC"),
			//HTTP: httpPlug,
			//HTTP: rest.Disabled,
			//HTTP: NewPlugin(UseDisabled()),
		}),
	)

	p := &ExamplePlugin{
		Deps: Deps{
			PluginName: PluginName,
			Log:        logging.ForPlugin(PluginName),
			GRPC:       myGRPC,
			//GRPC: grpc.DefaultPlugin,
		},
	}*/

	a := agent.NewAgent(agent.AllPlugins(p))

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}

// ExamplePlugin presents main plugin.
type ExamplePlugin struct {
	Log  logging.PluginLogger
	GRPC grpc.Server
}

func (plugin *ExamplePlugin) String() string {
	return PluginName
}

// Init demonstrates the usage of PluginLogger API.
func (plugin *ExamplePlugin) Init() error {
	plugin.Log.Info("Registering greeter")

	helloworld.RegisterGreeterServer(plugin.GRPC.GetServer(), &GreeterService{})

	return nil
}

func (plugin *ExamplePlugin) Close() error {
	return nil
}

// GreeterService implements GRPC GreeterServer interface (interface generated from protobuf definition file).
// It is a simple implementation for testing/demo only purposes.
type GreeterService struct{}

// SayHello returns error if request.name was not filled otherwise: "hello " + request.Name
func (*GreeterService) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if request.Name == "" {
		return nil, errors.New("not filled name in the request")
	}
	logrus.DefaultLogger().Infof("greeting client: %v", request.Name)

	return &helloworld.HelloReply{Message: "hello " + request.Name}, nil
}
