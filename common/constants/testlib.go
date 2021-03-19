package constants

var TestlibBinaryPrefixs = map[string]string{
	"generator":  "g_",
	"validator":  "",
	"checker":    "",
	"interactor": "",
}

var TestlibExitMsgMapping = []struct {
	ErrName     string
	JudgeResult int
	WithScore   bool
}{
	{ErrName: "ok", JudgeResult: JudgeFlagAC},
	{ErrName: "wrong answer", JudgeResult: JudgeFlagWA},
	{ErrName: "wrong output format", JudgeResult: JudgeFlagPE},
	{ErrName: "FAIL", JudgeResult: JudgeFlagSpecialJudgeError},
	{ErrName: "points", JudgeResult: JudgeFlagSpecialJudgeError}, // Unsupport
	{ErrName: "unexpected eof", JudgeResult: JudgeFlagPE},
	{ErrName: "partially correct", JudgeResult: JudgeFlagWA, WithScore: true},
	{ErrName: "What is the code", JudgeResult: JudgeFlagSpecialJudgeError},
}

var TestlibOutcomeMapping = map[string]int{
	"accepted":           JudgeFlagAC,
	"wrong-answer":       JudgeFlagWA,
	"presentation-error": JudgeFlagPE,
	"fail":               JudgeFlagSpecialJudgeError,
	"points":             JudgeFlagSpecialJudgeError, // Unsupport
	"relative-scoring":   JudgeFlagSpecialJudgeError, // Unsupport
	"unexpected-eof":     JudgeFlagPE,
	"partially-correct":  JudgeFlagWA,
	"reserved":           JudgeFlagSpecialJudgeError,
}
