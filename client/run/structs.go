package run

import "github.com/LanceLRQ/deer-common/persistence"

type RunOption struct {
    Persistence     *persistence.JudgeResultPersisOptions
    Clean           bool
    ConfigFile      string
    WorkDir         string
    ShowLog         bool
    LogLevel        int
    Language        string
    LibraryDir      string
    CodePath        string
    SessionId       string
    SessionRoot     string
}
