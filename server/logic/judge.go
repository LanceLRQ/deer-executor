package logic

import (
	"context"
	"github.com/LanceLRQ/deer-executor/v2/server/rpc"
	"time"
)

type JudgementServiceServerImpl struct {
	*rpc.UnimplementedJudgementServiceServer
}

func (s *JudgementServiceServerImpl) Ping(context context.Context, request *rpc.PingRequest) (*rpc.PingResponse, error) {
	return &rpc.PingResponse{
		Ready: true,
		Time:  time.Now().Unix(),
	}, nil
}
func (s *JudgementServiceServerImpl) StartJudgement(context context.Context, request *rpc.JudgementRequest) (*rpc.JudgementResponse, error) {
	return &rpc.JudgementResponse{}, nil
}
