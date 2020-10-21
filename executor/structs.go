package executor

import (
	"github.com/LanceLRQ/deer-executor/provider"
	"syscall"
)

const (
	JudgeFlagAC 	 						= 0   					// 0 Accepted
	JudgeFlagPE 	 						= 1	    				// 1 Presentation Error
	JudgeFlagTLE 							= 2						// 2 Time Limit Exceeded
	JudgeFlagMLE 	 						= 3						// 3 Memory Limit Exceeded
	JudgeFlagWA 	 						= 4	    				// 4 Wrong Answer
	JudgeFlagRE 	 						= 5	    				// 5 Runtime Error
	JudgeFlagOLE 	 						= 6						// 6 Output Limit Exceeded
	JudgeFlagCE 	 						= 7	    				// 7 Compile Error
	JudgeFlagSE 						 	= 8     				// 8 System Error

	JudgeFlagSpecialJudgeTimeout 		 	= 10    				// 10 Special Judger Time OUT
	JudgeFlagSpecialJudgeError 			 	= 11    				// 11 Special Judger ERROR
	JudgeFlagSpecialJudgeRequireChecker  	= 12 					// 12 Special Judger Finish, Need Standard Checkup
)

const (
	JudgeFileSizeLimit 						= 50 * 1024 * 1024  	// 50MB

	SpecialJudgeModeDisabled 				= 0
	SpecialJudgeModeChecker 				= 1
	SpecialJudgeModeInteractive 			= 2

	SpecialJudgeTimeLimit 					= 1 * 1000				// Unit: ms
	SpecialJudgeMemoryLimit 				= 256 * 1024			// Unit: kb
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
	TestCaseIn 		string				`json:"test_case_in"`			// Testcase input file path
	TestCaseOut		string				`json:"test_case_out"`			// Testcase output file path
}

type TestCaseResult struct {
	Id 				string				`json:"id"`						// Identifier

	TestCaseIn 		string				`json:"-"`						// Testcase input file path (internal)
	TestCaseOut		string				`json:"-"`						// Testcase output file path (internal)
	ProgramOut 		string				`json:"program_out"`			// Program-stdout file path
	ProgramError 	string				`json:"program_error"`			// Program-stderr file path
	ProgramLog 		string				`json:"program_log"`			// Program-log file path

	JudgerOut 		string				`json:"judger_out"`				// Special judger checker's stdout
	JudgerError 	string				`json:"judger_error"`			// Special judger checker's stderr
	JudgerLog		string				`json:"judger_log"`				// Special judger checker's log file
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

// Judge session
type JudgeSession struct {
	SessionId		string				`json:"session_id"`				// Judge Session Id
	SessionRoot		string				`json:"session_root"`			// Session Root Directory
	SessionDir		string				`json:"-"`						// Session Directory
	WorkDir			string				`json:"work_dir"`				// Working Directory
	CodeLangName 	string				`json:"code_lang_name"`			// Code file language name
	CodeFile	 	string				`json:"code_file"`				// Code File Path
	Commands 		[]string			`json:"commands"`				// Executable program commands
	TestCases		[]TestCase			`json:"test_cases"`				// Test cases
	TimeLimit 		int					`json:"time_limit"`				// Time limit (ms)
	MemoryLimit 	int					`json:"memory_limit"`			// Memory limit (KB)
	RealTimeLimit 	int					`json:"real_time_limit"`		// Real Time Limit (ms) (optional)
	FileSizeLimit 	int					`json:"file_size_limit"`		// File Size Limit (bytes) (optional)
	Uid 			int					`json:"uid"`					// User id (optional)
	StrictMode 		bool				`json:"strict_mode"`			// Strict Mode (if close, PE will be ignore)
	SpecialJudge  	SpecialJudgeOptions `json:"special_judge"`			// Special Judge Options

	compiler		provider.CodeCompileProviderInterface				// Compiler entity
}

