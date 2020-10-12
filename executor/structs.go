package executor

const (
	JudgeFlagAC 	int 	= 0   							// 0 Accepted
	JudgeFlagPE 	int 	= 1	    						// 1 Presentation Error
	JudgeFlagTLE 	int		= 2								// 2 Time Limit Exceeded
	JudgeFlagMLE 	int 	= 3								// 3 Memory Limit Exceeded
	JudgeFlagWA 	int 	= 4	    						// 4 Wrong Answer
	JudgeFlagRE 	int 	= 5	    						// 5 Runtime Error
	JudgeFlagOLE 	int 	= 6								// 6 Output Limit Exceeded
	JudgeFlagCE 	int 	= 7	    						// 7 Compile Error
	JudgeFlagSE 	int 	= 8     						// 8 System Error

	JudgeFlagSpecialJudgeTimeout 		int 	= 10    	// 10 Special Judger Time OUT
	JudgeFlagSpecialJudgeError 			int 	= 11    	// 11 Special Judger ERROR
	JudgeFlagSpecialJudgeRequireChecker int 	= 12 		// 12 Special Judger Finish, Need Standard Checkup
)

const (
	JudgeFileSizeLimit 				= 200 * 1024 * 1024  	// 200MB

	SpecialJudgeModeDisabled 		= 0
	SpecialJudgeModeChecker 		= 1
	SpecialJudgeModeInteractive 	= 2

	SpecialJudgeTimeLimit 			= 1 * 1000				// Unit: ms
	SpecialJudgeMemoryLimit 		= 256 * 1024			// Unit: kb
)

// Judge result
type JudgeResult struct {
	JudgeResult 	int 			`json:"judge_result"`			// Judge result flag number
	TimeUsed 		int				`json:"time_used"`				// Maximum time used
	MemoryUsed 		int				`json:"memory_used"`			// Maximum memory used
	ReSignum 		int				`json:"re_signal_num"`			// Runtime error signal number
	SameLines 		int				`json:"same_lines"`				// sameLines when WA
	TotalLines 		int				`json:"total_lines"`			// totalLines when WA
	SeInfo 			string			`json:"se_info"`				// SeInfo When System Error
	CeInfo 			string			`json:"ce_info"`				// CeInfo When CeInfo
}

type SpecialJudgeOptions struct {
	Mode 		int				`json:"mode"`				// Mode；0-Disabled；1-Normal；2-Interactor
	Checker 	string			`json:"checker"`			// Checker file path
	RedirectStd bool 			`json:"redirect_std"`		// Redirect target program's Stdout to checker's Stdin (checker mode)
	TimeLimit 	int				`json:"time_limit"`			// Time limit (ms)
	MemoryLimit int				`json:"memory_limit"`		// Memory limit (kb)
	Stdout 		string			`json:"stdout"`				// checker's stdout
	Stderr 		string			`json:"stderr"`				// checker's stderr
	Logfile		string			`json:"log_file"`			// checker's log file params
}

// Judge session
type JudgeSession struct {
	SessionId		string				`json:"session"`			// Judge Session Id
	SessionDir 		string				`json:"session_dir"`		// Session Directory
	CodeLangName 	string				`json:"code_lang_name"`		// Code file language name
	CodeFile	 	string				`json:"code_file"`			// Code File Path
	Commands [] 	string				`json:"commands"`			// Executable program commands
	TestCaseIn 		string				`json:"test_case_in"`		// Testcase input file path
	TestCaseOut		string				`json:"test_case_out"`		// Testcase output file path
	ProgramOut 		string				`json:"program_out"`		// Program-stdout file path
	ProgramError 	string				`json:"program_error"`		// Program-stderr file path
	ProgramLog 		string				`json:"program_log"`		// Program-log file path
	TimeLimit 		int					`json:"time_limit"`			// Time limit (ms)
	MemoryLimit 	int					`json:"memory_limit"`		// Memory limit (KB)
	RealTimeLimit 	int					`json:"real_time_limit"`	// Real Time Limit (ms) (optional)
	FileSizeLimit 	int					`json:"file_size_limit"`	// File Size Limit (bytes) (optional)
	Uid 			int					`json:"uid"`				// User id (optional)
	SpecialJudge  	SpecialJudgeOptions `json:"special_judge"`
}

