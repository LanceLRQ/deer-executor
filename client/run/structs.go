package run

import "github.com/LanceLRQ/deer-executor/v2/common/persistence"

// JudgementRunOption options for StartJudgement
type JudgementRunOption struct {
	Persistence *persistence.JudgeResultPersisOptions
	Clean       bool
	ConfigFile  string
	WorkDir     string
	ShowLog     bool
	LogLevel    int
	Language    string
	LibraryDir  string
	CodePath    string
	SessionID   string
	SessionRoot string
}
