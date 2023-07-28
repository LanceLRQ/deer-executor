package run

import (
	"github.com/LanceLRQ/deer-executor/v3/executor"
	"github.com/LanceLRQ/deer-executor/v3/executor/persistence"
	commonStructs "github.com/LanceLRQ/deer-executor/v3/executor/structs"
)

// StartJudgement to run a judge work.
func StartJudgement(options *JudgementRunOption) (*executor.JudgeSession, *commonStructs.JudgeResult, error) {
	judgeResult, judgeSession, err := runOnceJudge(options)
	if err != nil {
		return nil, nil, err
	}
	// Do clean (or benchmark on)
	if options.Clean {
		defer judgeSession.Clean()
	}

	// persistence
	if options.Persistence != nil {
		options.Persistence.SessionDir = judgeSession.SessionDir
		pack := persistence.NewJudgeResultPackage(judgeResult)
		err = pack.WritePackageFile(options.Persistence)
		if err != nil {
			return nil, nil, err
		}
	}
	return judgeSession, judgeResult, nil
}
