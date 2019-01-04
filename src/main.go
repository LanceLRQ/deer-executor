package deer

import (
	"bytes"
	"encoding/json"
	"fmt"
)

const (
	JUDGE_FLAG_AC int = 0   	//0 Accepted
	JUDGE_FLAG_PE int = 1	    //1 Presentation Error
	JUDGE_FLAG_TLE int = 2		//2 Time Limit Exceeded
	JUDGE_FLAG_MLE int = 3		//3 Memory Limit Exceeded
	JUDGE_FLAG_WA int = 4	    //4 Wrong Answer
	JUDGE_FLAG_RE int = 5	    //5 Runtime Error
	JUDGE_FLAG_OLE int = 6		//6 Output Limit Exceeded
	JUDGE_FLAG_CE int = 7	    //7 Compile Error
	JUDGE_FLAG_SE int = 8     	//8 System Error

	JUDGE_FLAG_SPJ_TIME_OUT int = 10    	// 10 Special Judger Time OUT
	JUDGE_FLAG_SPJ_ERROR int = 11    		// 11 Special Judger ERROR
	JUDGE_FLAG_SPJ_REQUIRE_CHECK int = 12 	// 12 Special Judger Finish, Need Standard Checkup

	JUDGE_FILE_SIZE_LIMIT = 100 * 1024  // kb
)

const (
	SPECIAL_JUDGE_MODE_DISABLED = 0
	SPECIAL_JUDGE_MODE_CHECKER = 1
	SPECIAL_JUDGE_MODE_INTERACTIVE = 2

	SPECIAL_JUDGE_TIME_LIMIT = 10 * 1000		// ms
	SPECIAL_JUDGE_MEMORY_LIMIT = 256 * 1024		// kb

)


type JudgeResult struct {
	JudgeResult int 			// Judge result flag number
	TimeUsed int				// Maximum time used
	MemoryUsed int				// Maximum memory used
	ReSignum int				// Runtime error signal number
	SameLines int				// sameLines when WA
	TotalLines int				// totalLines when WA
	SeInfo string				// SeInfo When System Error
	CeInfo string				// CeInfo When CeInfo
}


type JudgeOption struct {
	Commands [] string				// Executable program commands
	TestCaseIn string
	TestCaseOut string
	ProgramOut string
	ProgramError string
	TimeLimit int					// Time limit (ms)
	MemoryLimit int					// Memory limit (kb)
	FileSizeLimit int				// File Size Limit (kb)
	Uid int							// User id (optional)
	SpecialJudge struct {
		Mode int					// Mode
		Checker string				// Checker file path
		RedirectStd bool 			// Redirect target program's Stdout to checker's Stdin
		TimeLimit int				// Time limit (ms)
		MemoryLimit int				// Memory limit (kb)
		Stdout string				// checker's stdout
		Stderr string				// checker's stderr
	}
}

func (conf *JudgeResult) String() string {
	b, err := json.Marshal(*conf)
	if err != nil {
		return fmt.Sprintf("%+v", *conf)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *conf)
	}
	return out.String()
}


func Judge(options JudgeOption) (*JudgeResult, error) {
	judgeResult :=  new(JudgeResult)

	if options.SpecialJudge.Mode == SPECIAL_JUDGE_MODE_INTERACTIVE {
		err := InteractiveChecker(options, judgeResult)
		if err != nil {
			return nil, err
		}
		// Text Diff
		if judgeResult.JudgeResult == JUDGE_FLAG_SPJ_REQUIRE_CHECK {
			err = DiffText(options, judgeResult)
			if err != nil {
				return nil, err
			}
		}
	} else {
		// Run Program
		err := RunProgram(options, judgeResult)
		if err != nil {
			return nil, err
		}
		if options.SpecialJudge.Mode == SPECIAL_JUDGE_MODE_CHECKER && judgeResult.JudgeResult == JUDGE_FLAG_AC {
			err := CustomChecker(options, judgeResult)
			if err != nil {
				return nil, err
			}
			// Text Diff
			if judgeResult.JudgeResult == JUDGE_FLAG_SPJ_REQUIRE_CHECK {
				err = DiffText(options, judgeResult)
				if err != nil {
					return nil, err
				}
			}
		} else {
			// Text Diff
			if judgeResult.JudgeResult == JUDGE_FLAG_AC {
				err = DiffText(options, judgeResult)
				if err != nil {
					return nil, err
				}
			}

		}
	}
	return judgeResult, nil
}