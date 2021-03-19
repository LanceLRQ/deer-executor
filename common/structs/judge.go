package structs

import (
	"encoding/xml"
	"github.com/LanceLRQ/deer-executor/v2/common/logger"
)

// JudgeConfiguration 评测配置信息
type JudgeConfiguration struct {
	TestCases     []TestCase                    `json:"test_cases"`      // Test cases
	TimeLimit     int                           `json:"time_limit"`      // Time limit (ms)
	MemoryLimit   int                           `json:"memory_limit"`    // Memory limit (KB)
	RealTimeLimit int                           `json:"real_time_limit"` // Real Time Limit (ms) (optional)
	FileSizeLimit int                           `json:"file_size_limit"` // File Size Limit (bytes) (optional)
	UID           int                           `json:"uid"`             // User id (optional)
	StrictMode    bool                          `json:"strict_mode"`     // Strict Mode (if close, PE will be ignore)
	SpecialJudge  SpecialJudgeOptions           `json:"special_judge"`   // Special Judge Options
	Limitation    map[string]JudgeResourceLimit `json:"limitation"`      // Limitation
	Problem       ProblemContent                `json:"problem"`         // Problem Info
	TestLib       TestlibOptions                `json:"testlib"`         // testlib设置
	AnswerCases   []AnswerCase                  `json:"answer_cases"`    // Answer cases (用于生成Output)
	ConfigDir     string                        `json:"-"`               // 内部字段：config文件所在目录绝对路径
}

// AnswerCase 答案代码样例
// 优先使用Content访问，其次使用FileName
type AnswerCase struct {
	Name     string `json:"name"`      // Case name
	FileName string `json:"file_name"` // code file name
	Language string `json:"language"`  // code language, default is 'auto'
	Content  string `json:"content"`   // code content (optional)
}

// TestCase 测试数据
type TestCase struct {
	Handle           string `json:"handle"`            // Identifier
	Order            int    `json:"order"`             // Order (ASC)
	Name             string `json:"name"`              // Testcase name
	Input            string `json:"input"`             // Testcase input file path
	Output           string `json:"output"`            // Testcase output file path
	Visible          bool   `json:"visible"`           // Is visible(for oj)
	Enabled          bool   `json:"enabled"`           // Is enabled
	UseGenerator     bool   `json:"use_genarator"`     // Use generator
	Generator        string `json:"generator"`         // Generator script
	ValidatorVerdict bool   `json:"validator_verdict"` // Testlib validator's result
	ValidatorComment string `json:"validator_comment"` // Testlib validator's output
}

// SpecialJudgeOptions 特殊评测设置
type SpecialJudgeOptions struct {
	Name               string                    `json:"name"`                 // Name, default is "checker"
	Mode               int                       `json:"mode"`                 // Mode；0-Disabled；1-Normal；2-Interactor
	CheckerLang        string                    `json:"checker_lang"`         // Checker languages, support gcc, g++(default) and golang, not support auto!
	Checker            string                    `json:"checker"`              // Checker file path (Use code file is better then compiled binary!)
	RedirectProgramOut bool                      `json:"redirect_program_out"` // Redirect target program's STDOUT to checker's STDIN (checker mode). if not, redirect testcase-in file to checker's STDIN
	TimeLimit          int                       `json:"time_limit"`           // Time limit (ms)
	MemoryLimit        int                       `json:"memory_limit"`         // Memory limit (kb)
	UseTestlib         bool                      `json:"use_testlib"`          // If use testlib, checker will only support c++
	CheckerCases       []SpecialJudgeCheckerCase `json:"checker_cases"`        // Special Judge checker cases (for Testlib, exclude interactor mode)
}

// SpecialJudgeCheckerCase 特判检查器样例
// Special Judge checker case item
type SpecialJudgeCheckerCase struct {
	Input           string `json:"input"`            // Input (1k Limit)
	Output          string `json:"output"`           // Output (1k Limit)
	Answer          string `json:"answer"`           // Answer (1k Limit)
	Verdict         bool   `json:"verdict"`          // Is verdict? (下边俩是否相同)
	ExpectedVerdict int    `json:"expected_verdict"` // Expected judge result (flag) (期望的判定结果)
	CheckerVerdict  int    `json:"checker_verdict"`  // (testlib/classical)checker's judge result (flag) (检查器的判定结果)
	CheckerComment  string `json:"checker_comment"`  // (testlib/classical) checker's output (检查器输出的信息)
}

// TestlibOptions TestLib设置 (只支持c++版本的testlib)
// Testlib Options (we only support c++ verion)
type TestlibOptions struct {
	Version        string                 `json:"version"`        // Testlib version (预留，不太考虑实现)
	Validator      string                 `json:"validator"`      // Validator file
	ValidatorName  string                 `json:"validator_name"` // Validator name (compile target name)
	Generators     []TestlibGenerator     `json:"generators"`     // Validator cases
	ValidatorCases []TestlibValidatorCase `json:"validator_case"` // Validator cases
	// 未来这边可以加入对拍(stress)功能
}

// TestlibGenerator Testlib Generator
type TestlibGenerator struct {
	Name   string `json:"name"`   // Generator name
	Source string `json:"source"` // Source code file
}

// TestlibValidatorCase Testlib validator 样例
type TestlibValidatorCase struct {
	Input            string `json:"input"`             // Input (1k Limit)
	Verdict          bool   `json:"verdict"`           // Is verdict? (下边俩是否相同)
	ExpectedVerdict  bool   `json:"expected_verdict"`  // Expected result
	ValidatorVerdict bool   `json:"validator_verdict"` // Testlib validator's result
	ValidatorComment string `json:"validator_comment"` // Testlib validator's output
}

// JudgeResult 评测结果信息
type JudgeResult struct {
	SessionID   string                `json:"session_id"`   // Judge Session Id
	JudgeResult int                   `json:"judge_result"` // Judge result flag number
	TimeUsed    int                   `json:"time_used"`    // Maximum time used
	MemoryUsed  int                   `json:"memory_used"`  // Maximum memory used
	TestCases   []TestCaseResult      `json:"test_cases"`   // Testcase Results
	ReInfo      string                `json:"re_info"`      // ReInfo when Runtime Error or special judge Runtime Error
	SeInfo      string                `json:"se_info"`      // SeInfo when System Error
	CeInfo      string                `json:"ce_info"`      // CeInfo when Compile Error
	JudgeLogs   []logger.JudgeLogItem `json:"judge_logs"`   // Judge Logs
}

// TestCaseResult 测试数据运行结果
type TestCaseResult struct {
	Handle       string `json:"handle"`        // Identifier
	Input        string `json:"-"`             // Testcase input file path (internal)
	Output       string `json:"-"`             // Testcase output file path (internal)
	ProgramOut   string `json:"program_out"`   // Program-stdout file path
	ProgramError string `json:"program_error"` // Program-stderr file path

	CheckerOut    string `json:"checker_out"`    // Special judge checker's stdout
	CheckerError  string `json:"checker_error"`  // Special judge checker's stderr
	CheckerReport string `json:"checker_report"` // Special judge checker's report file

	JudgeResult    int `json:"judge_result"`    // Judge result flag number
	PartiallyScore int `json:"partially_score"` // Testlib Partially Score or Math.floor(SameLines / TotalLines)

	TextDiffLog string `json:"text_diff_log"` // Text Checkup Log
	TimeUsed    int    `json:"time_used"`     // Maximum time used
	MemoryUsed  int    `json:"memory_used"`   // Maximum memory used
	ReSignum    int    `json:"re_signal_num"` // Runtime error signal number
	SameLines   int    `json:"same_lines"`    // Same lines when WA
	TotalLines  int    `json:"total_lines"`   // Total lines when WA
	ReInfo      string `json:"re_info"`       // ReInfo when Runtime Error or special judge Runtime Error
	SeInfo      string `json:"se_info"`       // SeInfo when System Error
	CeInfo      string `json:"ce_info"`       // CeInfo when Compile Error

	SPJExitCode   int    `json:"spj_exit_code"`     // Special judge exit code
	SPJTimeUsed   int    `json:"spj_time_used"`     // Special judge maximum time used
	SPJMemoryUsed int    `json:"spj_memory_used"`   // Special judge maximum memory used
	SPJReSignum   int    `json:"spj_re_signal_num"` // Special judge runtime error signal number
	SPJMsg        string `json:"spj_msg"`           // Special judge checker  msg
}

// JudgeResourceLimit 评测资源限制信息
type JudgeResourceLimit struct {
	TimeLimit     int `json:"time_limit"`      // Time limit (ms)
	MemoryLimit   int `json:"memory_limit"`    // Memory limit (KB)
	RealTimeLimit int `json:"real_time_limit"` // Real Time Limit (ms) (optional)
	FileSizeLimit int `json:"file_size_limit"` // File Size Limit (bytes) (optional)
}

// TestlibCheckerResult Testlib检查器报告
type TestlibCheckerResult struct {
	XMLName     xml.Name `xml:"result"`
	Outcome     string   `xml:"outcome,attr"`
	PcType      string   `xml:"pctype,attr"`
	Description string   `xml:",innerxml"`
}
