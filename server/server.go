package server

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/v2/server/logic"
	"github.com/LanceLRQ/deer-executor/v2/server/rpc"
	"github.com/LanceLRQ/deer-executor/v2/server/server_config"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"net"
)

func LaunchJudgeService(c *cli.Context) error {
	hostname := fmt.Sprintf("%s:%d", server_config.GRPCConfig.Host, server_config.GRPCConfig.Port)
	srv := grpc.NewServer()
	rpc.RegisterJudgementServiceServer(srv, &logic.JudgementServiceServerImpl{})
	listener, err := net.Listen("tcp", hostname)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	fmt.Printf("Deer-Executor judgement rpc service listen at: %s\n", hostname)

	err = srv.Serve(listener)
	if err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
