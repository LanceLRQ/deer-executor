package run

import (
    "github.com/LanceLRQ/deer-common/persistence/judge_result"
    commonStructs "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-executor/v2/executor"
)

// Start Judgement
func StartJudgement(options *RunOption) (*executor.JudgeSession, *commonStructs.JudgeResult, error) {
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
        err = judge_result.PersistentJudgeResult(judgeResult, options.Persistence)
        if err != nil {
            return nil, nil, err
        }
    }
    return judgeSession, judgeResult, nil
}
