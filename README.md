<h1 align="center">Deer Executor</h1>
<p align="center">A program executor for online judge written by golang.</p>

English | [简体中文](./README-zh_CN.md)

## ✨ Features
 - Most languages supported.
 - You can  build and run it standalone.
 - Support Linux and Darwin(MacOS), maybe Windows in the future.
 
## 📦 Get && Install

```
go get github.com/LanceLRQ/deer-executor
```
**Environment:** Go 1.11+ is best!

## 🔨 Usage

```
import (
    "fmt"
    "github.com/LanceLRQ/deer-executor"
    "github.com/LanceLRQ/deer-executor/compile"
)

// Create a compiler provider
compiler := new(deer_compile.GnucCompileProvider)
compiler.Init("#include<stdio.h>\nint main(){ return 0; }", "/tmp")    // The second argument means the work directory.

// Do compile
success, ceinfo := compiler.Compile()
if !success {
    fmt.Println("Compile Error: " + ceinfo)
}

// Get compile result
cmds := compiler.GetRunArgs()

judgeOptions := deer_executor.JudgeOption {

    // Executable program commands
    // Commands:      []string{ "/tmp/a.out", "-a", "123" },      // It means: /tmp/a.out -a 123
    Commands:      cmds, 
    
    // Resource Limitation
    TimeLimit:     1000,                     // Maximum time limit (ms)
    MemoryLimit:   32768,                    // Maximum memory limit (Kbytes)
    FileSizeLimit: 100 * 1024 * 1024,        // Maximum file size output limit (Kbytes)
    
    // Test Cases
    TestCaseIn:    "/data/1.in",             // TestCase-In file path
    TestCaseOut:   "/data/1.out",            // TestCase-Out file path
    ProgramOut:    "/tmp/user.out",          // Program's stdout file path
    ProgramError:  "/tmp/user.err",          // Program's stderr file path
    
    // Special Judge
    SpecialJudge struct {
        Mode int                    // Mode
        Checker string				// Checker file path
        RedirectStd bool 			// Redirect target program's Stdout to checker's Stdin (only for checker mode)
        TimeLimit int				// Time limit (ms)
        MemoryLimit int				// Memory limit (kb)
        Stdout string				// checker's stdout
        Stderr string				// checker's stderr
    }
    // Other
    Uid:    -1,                              // Linux user id (optional)
}

judgeResult, err := deer_executor.Judge(judgeOptions)
```
judgeResult define like this:
```
type JudgeResult struct {
	JudgeResult int 			// Judge result flag number
	TimeUsed int				// Maximum time used (ms)
	MemoryUsed int				// Maximum memory used  (Kbytes)
	ReSignum int				// Runtime error signal number
	SameLines int				// Same Lines when WA
	TotalLines int				// Total Lines when WA
	SeInfo string				// SeInfo When System Error
	CeInfo string				// CeInfo When CeInfo
}
```

## 💡 Special Judge
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

int main(int argc,char **argv) {
    // Do anything you want.
    
    return RESULT_AC;
}
```
_REQUIRE_DEFAULT_CHECKER_

  Special judge checker sometimes not only a checker, but also a processor program. You can use the checker to process the target program's output, _e.p_ keep two decimal for type _double_ and so on. After that you can return **REQUIRE_DEFAULT_CHECKER** for calling the default text-diff checker supported from deer-executor.

## 🧬 Compile

To compile code, deer-executor supported:
```
GCC、GNU C++、Java、Python2、Python3、Golang、NodeJS、PHP、Ruby
```
Sure, you can add any compiler you like. Deer make an interface **CodeCompileProviderInterface** 
```
type CodeCompileProviderInterface interface {

    // Initial the provider, set code content and work directory.
    Init(code string, workDir string) error
    
    // Compile the code. it must be run after Init() called.
    Compile() (result bool, errmsg string)
    
    // Get compiled program's file path and run arguments.
    GetRunArgs() (args []string)

    // If your compiler is a real-time compiler, like python.
    // It should't compile first, and will output compile error when running.
    // So you can use it to check if VM output a compile error
    IsCompileError(remsg string) bool
    
    // Is it a realtime compiler?
    IsRealTime() bool
    
    // Is code compiled?
    IsReady() bool

	
    /** 
     ** Private Methods
     **/

    // Write code content to file before compile.
    initFiles(codeExt string, programExt string) error
    
	// Call the system shell
	shell(commands string) (success bool, errout string)
	// Save your code content to file
	saveCode() error
	// Check if work dir exists
	checkWorkDir() error
}

type CodeCompileProvider struct {
	CodeCompileProviderInterface
	
	codeContent string		        // Code content
	realTime bool			        // Is it a realtime compiler?
	isReady bool			        // Is code compiled?
	codeFileName string             // Target code file name
	codeFilePath string			    // Target code file  path
	programFileName string          // Target program file name
	programFilePath string		    // Target program file path
	workDir string			        // Work Directory
}
```
  

## 🤝 Thanks！

First, I'm really appreciate to the author of [Loco's runner](https://github.com/dojiong/Lo-runner). 

Then, my classmates Wolf Zheng and Tosh Qiu propose the _interactive judge_ and describe how it works.
 
Finally，thanks to my alma mater [Beijing Normal University (Zhuhai)](http://www.bnuz.edu.cn), [BNUZ IT college](http://itc.bnuz.edu.cn), [ACM association](http://acm.bnuz.edu.cn) and WeJudge team

## 🔗 Links

📃 My blog：[https://www.lanrongqi.com](https://www.lanrongqi.com)

🖥️ WeJudge：

[https://www.wejudge.net](https://www.wejudge.net) 

[https://oj.bnuz.edu.cn](https://oj.bnuz.edu.cn)

[WeJudge 1.0 Open Source](https://github.com/LanceLRQ/wejudge)


**_We welcome all contributions. You can submit any ideas as pull requests or as GitHub issues. have a good time! :)_**