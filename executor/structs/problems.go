package structs

// ProblemIOSample 题目Input/Output样例信息
type ProblemIOSample struct {
	Input  string `json:"input"`  // Input sample
	Output string `json:"output"` // Output sample
}

// ProblemContent 题目正文信息  (for oj)
type ProblemContent struct {
	Author      string                   `json:"author"`       // Problem author
	Source      string                   `json:"source"`       // Problem source
	Description string                   `json:"description"`  // Description
	Input       string                   `json:"input"`        // Input requirements
	Output      string                   `json:"output"`       // Output requirements
	Sample      []ProblemIOSample        `json:"sample"`       // Sample cases
	Tips        string                   `json:"tips"`         // Solution tips
	ProblemType int                      `json:"problem_type"` // 题目类型
	DemoCases   map[string]JudgeDemoCase `json:"demo_cases"`   // 代码填空样例数据
}

// JudgeDemoCase 代码填空样例 (for oj)
type JudgeDemoCase struct {
	Handle  string            `json:"handle"`  // handle
	Name    string            `json:"name"`    // 代码区域名称
	Answers map[string]string `json:"answers"` // 回答信息
	Demo    string            `json:"demo"`    // 样例代码（预设用）
	Line    int               `json:"line"`    // 插入位置
}
