package agent

import (
	"fmt"
	agentConfig "github.com/LanceLRQ/deer-executor/v2/agent/config"
	"github.com/LanceLRQ/deer-executor/v2/agent/logic"
	"github.com/LanceLRQ/deer-executor/v2/agent/rpc"
	"github.com/LanceLRQ/deer-executor/v2/client"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"net"
)

func LaunchJudgeService(c *cli.Context) error {
	// Load judge enviroment configuration
	err := client.LoadSystemConfiguration()
	if err != nil {
		return err
	}

	hostname := fmt.Sprintf("%s:%d", agentConfig.GRPCConfig.Host, agentConfig.GRPCConfig.Port)
	srv := grpc.NewServer()
	rpc.RegisterJudgementServiceServer(srv, &logic.JudgementServiceServerImpl{})
	listener, err := net.Listen("tcp", hostname)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	fmt.Printf("Deer-Executor judgement agent rpc service listen at: %s\n", hostname)

	err = srv.Serve(listener)
	if err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
