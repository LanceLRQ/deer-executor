package executor

import (
	"encoding/json"
	"github.com/LanceLRQ/deer-executor/provider"
	"io/ioutil"
	"path"
	"path/filepath"
	"syscall"
)

type ProcessInfo struct {
	Pid uintptr							`json:"pid"`
	Status syscall.WaitStatus			`json:"status"`
	Rusage syscall.Rusage				`json:"rusage"`
}

type SpecialJudgeOptions struct {
	Mode 				int				`json:"mode"`					// Mode；0-Disabled；1-Normal；2-Interactor
	Checker 			string			`json:"checker"`				// Checker file path
	RedirectProgramOut 	bool 			`json:"redirect_program_out"`	// Redirect target program's STDOUT to checker's STDIN (checker mode). if not, redirect testcase-in file to checker's STDIN
	TimeLimit 			int				`json:"time_limit"`				// Time limit (ms)
	MemoryLimit 		int				`json:"memory_limit"`			// Memory limit (kb)
}

type TestCase struct {
	Id 				string				`json:"id"`						// Identifier
	Name 			string				`json:"name"`					// Testcase name
	TestCaseIn 		string				`json:"test_case_in"`			// Testcase input file path
	TestCaseOut		string				`json:"test_case_out"`			// Testcase output file path
}

type TestCaseResult struct {
	Id 				string				`json:"id"`						// Identifier

	TestCaseIn 		string				`json:"-"`						// Testcase input file path (internal)
	TestCaseOut		string				`json:"-"`						// Testcase output file path (internal)
	ProgramOut 		string				`json:"program_out"`			// Program-stdout file path
	ProgramError 	string				`json:"program_error"`			// Program-stderr file path

	JudgerOut 		string				`json:"judger_out"`				// Special judger checker's stdout
	JudgerError 	string				`json:"judger_error"`			// Special judger checker's stderr
	JudgerReport	string				`json:"judger_report"`			// Special judger checker's report file

	JudgeResult 	int 				`json:"judge_result"`			// Judge result flag number
	TextDiffLog		string				`json:"text_diff_log"`			// Text Checkup Log
	TimeUsed 		int					`json:"time_used"`				// Maximum time used
	MemoryUsed 		int					`json:"memory_used"`			// Maximum memory used
	ReSignum 		int					`json:"re_signal_num"`			// Runtime error signal number
	SameLines 		int					`json:"same_lines"`				// sameLines when WA
	TotalLines 		int					`json:"total_lines"`			// totalLines when WA
	ReInfo 			string				`json:"re_info"`				// ReInfo when Runtime Error or special judge Runtime Error
	SeInfo 			string				`json:"se_info"`				// SeInfo when System Error
	CeInfo 			string				`json:"ce_info"`				// CeInfo when Compile Error

	SPJExitCode  	int					`json:"spj_exit_code"`			// Special judge exit code
	SPJTimeUsed 	int					`json:"spj_time_used"`			// Special judge maximum time used
	SPJMemoryUsed 	int					`json:"spj_memory_used"`		// Special judge maximum memory used
	SPJReSignum 	int					`json:"spj_re_signal_num"`		// Special judge runtime error signal number
}

// Judge result
type JudgeResult struct {
	SessionId 		string				`json:"session_id"`				// Judge Session Id
	JudgeResult 	int 				`json:"judge_result"`			// Judge result flag number
	TimeUsed 		int					`json:"time_used"`				// Maximum time used
	MemoryUsed 		int					`json:"memory_used"`			// Maximum memory used
	TestCases		[]TestCaseResult	`json:"test_cases"`				// Testcase Results
	ReInfo 			string				`json:"re_info"`				// ReInfo when Runtime Error or special judge Runtime Error
	SeInfo 			string				`json:"se_info"`				// SeInfo when System Error
	CeInfo 			string				`json:"ce_info"`				// CeInfo when Compile Error
}

type ProblemIOSample struct {
	Input		string					`json:"input"`					// Input sample
	Output 		string					`json:"output"`					// Output sample
}

type ProblemContent struct {
	Author 			string				`json:"author"`					// Problem author
	Source 			string				`json:"source"`					// Problem source
	Description 	string				`json:"description"`			// Description
	Input			string				`json:"input"`					// Input requirements
	Output 			string				`json:"output"`					// Output requirements
	Sample			[]ProblemIOSample	`json:"sample"`					// Sample cases
	Tips 			string				`json:"tips"`					// Solution tips
}

type JudgeLimit struct {
	TimeLimit 		int					`json:"time_limit"`				// Time limit (ms)
	MemoryLimit 	int					`json:"memory_limit"`			// Memory limit (KB)
	RealTimeLimit 	int					`json:"real_time_limit"`		// Real Time Limit (ms) (optional)
	FileSizeLimit 	int					`json:"file_size_limit"`		// File Size Limit (bytes) (optional)
}

// Judge session
type JudgeSession struct {
	SessionId		string					`json:"-"`						// Judge Session Id
	SessionRoot		string					`json:"-"`						// Session Root Directory
	SessionDir		string					`json:"-"`						// Session Directory
	ConfigFile 		string					`json:"-"`						// Config file
	ConfigDir 		string					`json:"-"`						// Config file dir
	CodeLangName 	string					`json:"code_lang_name"`			// Code file language name
	CodeFile	 	string					`json:"-"`						// Code File Path
	Commands 		[]string				`json:"-"`						// Executable program commands
	TestCases		[]TestCase				`json:"test_cases"`				// Test cases
	TimeLimit 		int						`json:"time_limit"`				// Time limit (ms)
	MemoryLimit 	int						`json:"memory_limit"`			// Memory limit (KB)
	RealTimeLimit 	int						`json:"real_time_limit"`		// Real Time Limit (ms) (optional)
	FileSizeLimit 	int						`json:"file_size_limit"`		// File Size Limit (bytes) (optional)
	Uid 			int						`json:"uid"`					// User id (optional)
	StrictMode 		bool					`json:"strict_mode"`			// Strict Mode (if close, PE will be ignore)
	SpecialJudge  	SpecialJudgeOptions 	`json:"special_judge"`			// Special Judge Options

	Limitation		map[string]JudgeLimit	`json:"limitation"`				// Limitation
	Problem			ProblemContent			`json:"problem"`				// Problem Info

	compiler		provider.CodeCompileProviderInterface				// Compiler entity
}

func NewSession(configFile string) (*JudgeSession, error) {
	session := JudgeSession{}
	session.SessionRoot = "/tmp"
	session.CodeLangName = "auto"
	session.Uid = -1
	session.TimeLimit = 1000
	session.MemoryLimit = 65535
	session.StrictMode = true
	session.FileSizeLimit = 50 * 1024 * 1024
	session.SpecialJudge.Mode = 0
	session.SpecialJudge.RedirectProgramOut = true
	session.SpecialJudge.TimeLimit = 1000
	session.SpecialJudge.MemoryLimit = 65535
	if configFile != "" {
		configFileAbsPath, err := filepath.Abs(configFile)
		if err != nil {
			return nil, err
		}
		session.ConfigFile = configFileAbsPath
		session.ConfigDir = path.Dir(configFileAbsPath)
		cbody, err := ioutil.ReadFile(configFileAbsPath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(cbody, &session)
		if err != nil {
			return nil, err
		}
	}
	return &session, nil
}

