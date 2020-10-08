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
	JudgeResult 	int 			// Judge result flag number
	TimeUsed 		int				// Maximum time used
	MemoryUsed 		int				// Maximum memory used
	ReSignum 		int				// Runtime error signal number
	SameLines 		int				// sameLines when WA
	TotalLines 		int				// totalLines when WA
	SeInfo 			string			// SeInfo When System Error
	CeInfo 			string			// CeInfo When CeInfo
}

// Judge options
type JudgeOptions struct {
	Commands [] 	string			// Executable program commands
	TestCaseIn 		string			// Testcase input file path
	TestCaseOut		string			// Testcase output file path
	ProgramOut 		string			// Program-stdout file path
	ProgramError 	string			// Program-stderr file path
	TimeLimit 		int				// Time limit (ms)
	MemoryLimit 	int				// Memory limit (KB)
	RealTimeLimit 	int				// Real Time Limit (ms) (optional)
	FileSizeLimit 	int				// File Size Limit (bytes) (optional)
	Uid 			int				// User id (optional)
	SpecialJudge struct {
		Mode 		int				// Mode；0-Disabled；1-Normal；2-Interactor
		Checker 	string			// Checker file path
		RedirectStd bool 			// Redirect target program's Stdout to checker's Stdin
		TimeLimit 	int				// Time limit (ms)
		MemoryLimit int				// Memory limit (kb)
		Stdout 		string			// checker's stdout
		Stderr 		string			// checker's stderr
	}
}

