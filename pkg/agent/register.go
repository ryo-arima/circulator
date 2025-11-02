package agent

import (
	"net"

	"google.golang.org/grpc"

	"github.com/ryo-arima/circulator/pkg/agent/controller"
	"github.com/ryo-arima/circulator/pkg/config"
)

// RegisterGRPCServices registers all gRPC services with Clean Architecture dependencies
// Architecture: Controller -> Usecase -> Repository (API/Local/Pulsar) -> Config
func RegisterGRPCServices(conf config.BaseConfig) *grpc.Server {
	conf.Logger.INFO(config.ARSGSR, "Starting gRPC service registration")

	// Create gRPC server
	server := grpc.NewServer()

	// Initialize controllers (presentation layer)
	conf.Logger.DEBUG(config.ARIC, "Initializing controllers")
	streamController, err := controller.NewStreamController(conf)
	if err != nil {
		conf.Logger.ERROR(config.ARFISC, "Failed to initialize stream controller", map[string]interface{}{
			"error": err.Error(),
		})
		return nil
	}
	commonController := controller.NewCommonController(conf)

	// TODO: Register gRPC services when proto files are available
	// pb.RegisterAgentServiceServer(server, streamController)
	// pb.RegisterCommonServiceServer(server, commonController)

	_ = streamController // TODO: Remove these lines when proto services are registered
	_ = commonController

	conf.Logger.INFO(config.ARGRPC, "gRPC services registration setup completed", map[string]interface{}{
		"stream_controller_type": "initialized",
		"common_controller_type": "initialized",
		"ready_for":              "proto service registration",
	})

	return server
}

// StartGRPCServer starts the gRPC server with all registered services
func StartGRPCServer(conf config.BaseConfig, port string) error {
	conf.Logger.INFO(config.ARSGRPC, "Starting gRPC server", map[string]interface{}{
		"port": port,
	})

	// Register all services
	server := RegisterGRPCServices(conf)

	// Create listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		conf.Logger.ERROR(config.ARFTLOP, "Failed to listen on port", map[string]interface{}{
			"port":  port,
			"error": err.Error(),
		})
		return err
	}

	conf.Logger.INFO(config.ARGRPCS, "gRPC server starting", map[string]interface{}{
		"port": port,
	})

	// Start serving
	if err := server.Serve(lis); err != nil {
		conf.Logger.ERROR(config.ARFTSGRPC, "Failed to serve gRPC server", map[string]interface{}{
			"port":  port,
			"error": err.Error(),
		})
		return err
	}

	return nil
}
