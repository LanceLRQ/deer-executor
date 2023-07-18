package logic

import (
	"context"
	"encoding/json"
	"github.com/LanceLRQ/deer-executor/v3/agent/rpc"
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

func convertToJSON(a interface{}) string {
	rel, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(rel)
}

func (s *JudgementServiceServerImpl) StartJudgement(context context.Context, request *rpc.JudgementRequest) (*rpc.JudgementResponse, error) {
	// check requset
	if err := checkJudgeRequsetArgs(request); err != nil {
		return nil, err
	}

	judgeResult, presistFile, err := runRpcJudge(request)
	if err != nil {
		return nil, err
	}

	return &rpc.JudgementResponse{
		JudgeFlag:         rpc.JudgeFlag(judgeResult.JudgeResult),
		SessionId:         judgeResult.SessionID,
		ResultData:        convertToJSON(judgeResult),
		ResultPackageFile: presistFile,
	}, nil
}
