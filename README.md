<h1 align="center">Deer Executor</h1>
<p align="center">A program executor for online judge written by golang.</p>

English | [简体中文](./README-zh_CN.md)

## ✨ Features
 - Most languages supported.
 - You can standalone build and run it
 - Support Linux and Darwin(MacOS), maybe Windows in the future.
 
## 📦 Get && Install

```
go get github.com/LanceLRQ/deer-executor
```

## 🔨 Usage

```
import "github.com/LanceLRQ/deer-executor"

judgeOptions := deer.JudgeOption {__

    // Executable program commands
    Commands:      []string{ "/tmp/a.out", "-a", "123" },      // It means: /tmp/a.out -a 123
    
    // Resource Limitation
    TimeLimit:     1000,                     // Real-time limit (ms)
    MemoryLimit:   32768,                    // Maximum memory limit (Kbytes)
    FileSizeLimit: 100 * 1024 * 1024,        // Maximum file size output limit (Kbytes)
    
    // Test Cases
    TestCaseIn:    "/data/1.in",             // TestCase-In file path
    TestCaseOut:   "/data/1.out",            // TestCase-Out file path
    ProgramOut:    "/tmp/user.out",          // Program's stdout file path
    ProgramError:  "/tmp/user.err",          // Program's stderr file path
    
    // Special Judge
    SpecialJudge:	0,                      // Special judge mode: 0-disabled, 1-checker, 2-interactive
    SpecialJudgeChecker: "/data/judger.out",    // Special judger filepath, it must be a executable binary program
    SpecialJudgeOut: "/tmp/spj.out",            // Special judger's stdout file path
    SpecialJudgeError: "/tmp/spj.err",          // Special judger's stderr file path
    // Other
    Uid:    0,                              // Linux user id (optional)
}

judgeResult, err := deer.Judge(judgeOption)
```
judgeResult define like this:
```
type JudgeResult struct {
	JudgeResult int 			// Judge result flag number
	TimeUsed int				// Maximum time used (ms)
	MemoryUsed int				// Maximum memory used  (kb)
	ReSignum int				// Runtime error signal number
	SameLines int				// sameLines when WA
	TotalLines int				// totalLines when WA
	SeInfo string				// SeInfo When System Error
	CeInfo string				// CeInfo When CeInfo
}
```

## ⌨ Special Judge
Special Judge supported two modes:

 - Checker Mode
 - Interactive Mode
 
**Checker Mode** Deer-executor will run the target program first. When it finished without any error, deer-executor will call the special judge checker. The checker should check up the target program's answer text, and exit with a code to tell deer-executor finally result. 

The special judge checker's arguments is:
```
./checker [1] [2] [3]
```
[1]: TestCase-In File; [2]: TestCases-Out File; [3]: Answer File


**Interactive Mode** Deer-executor will run the target program and special judge checker at the same time, redirect checker's _stdout_ to programs's _stdin_ and checker's _stdin_ from program's _stdout_. Deer-executor will use checker's exit-code as the result.
The special judge checker's arguments is:
```
./checker [1] [2] [3]
```
[1]: TestCase-In File; [2]: TestCases-Out File; [3] Run-result File

Run-result: Maybe you can output your communication with program, it can be the special judge checker's logs.

**Exit Code**

Special judge checker report the judge result with it's exit code. like this (checker.c):
```
#define RESULT_AC 0
#define RESULT_PE 1
#define RESULT_WA 4
#define RESULT_OLE 6
#define REQUIRE_DEFAULT_CHECKER 12

int main() {
    // Do anythings you want.
    
    return RESULT_AC;
}
```
_REQUIRE_DEFAULT_CHECKER_

  Special judge checker sometimes not only a checker, but also a processor program. You can using the checker to process the target program's output, _e.p_ keep two decimal for type _double_ and so on. After that you can return **REQUIRE_DEFAULT_CHECKER** for calling the default text-diff checker supported from deer-executor.
  
## 🤝 Thanks！

First, I really appreciate to the author of [Loco's runner](https://github.com/dojiong/Lo-runner). 

Then, my classmates Wolf Zheng and Tosh Qiu propose the _interactive judge_ and describe how it works. 