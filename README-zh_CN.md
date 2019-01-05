<h1 align="center">Deer Executor</h1>
<p align="center">一个基于Go语言开发的代码判题内核</p>

[English](./README.md) | 简体中文

## ✨ 特性
 - 支持多种编程语言程序的判定。
 - 支持独立地工作，不依赖OJ系统。
 - 支持Linux和MacOS平台，未来将会支持Windows平台。
 
## 📦 安装

```
go get github.com/LanceLRQ/deer-executor
```
安装环境：强烈建议使用Go 1.11以上版本

## 🔨 使用

```
import (
    "fmt"
    "github.com/LanceLRQ/deer-executor"
)

// 构造一个编译提供程序，这里以C语言为例
compiler := new(compile.GnucCompileProvider)
compiler.Init("#include<stdio.h>\nint main(){ return 0; }", "/tmp")    // 第二个参数为工作目录

// 执行编译
success, ceinfo := gnuc.Compile()
if !success {
    fmt.Println("Compile Error: " + ceinfo)
}

// 获取编译后目标程序的执行参数
cmds := compiler.GetRunArgs()

judgeOptions := deer.JudgeOption {

    // 目标可执行程序的命令
    // Commands:      []string{ "/tmp/a.out", "-a", "123" },      // 参数列表以数组的形式存在，上述可以视作调用了 /tmp/a.out -a 123
    Commands:      cmds, 
    
    // 资源限制部分
    TimeLimit:     1000,                     // 最大运行时间限制 (ms)
    MemoryLimit:   32768,                    // 最大内存使用限制 (Kbytes)
    FileSizeLimit: 100 * 1024 * 1024,        // 最大文件输出限制 (Kbytes)
    
    // Test Cases
    TestCaseIn:    "/data/1.in",             // 测试数据输入文件位置
    TestCaseOut:   "/data/1.out",            // 测试数据输出文件位置
    ProgramOut:    "/tmp/user.out",          // 目标程序输出文件位置(stdout)
    ProgramError:  "/tmp/user.err",          // 目标程序错误信息输出文件位置(stderr)
    
    // Special Judge
    SpecialJudge struct {
        Mode int                    // 特殊评测模式: 0-禁用, 1-结果检查模式, 2-交互模式
        Checker string				// 特殊评测裁判程序路径, 它必须是一个可执行程序
        RedirectStd bool 			// 重定向目标程序的输出文件到裁判的stdin （交互模式无效）
        TimeLimit int				// 时间限制，值为0时，默认10秒超时。(ms)
        MemoryLimit int				// 内存限制，值为0时，默认256MB可用。(kb)
        Stdout string				// 裁判程序输出文件位置(stdout)
        Stderr string				// 裁判程序错误信息输出文件位置(stderr)
    }
    // Other
    Uid:    -1,                              // 执行时的Linux用户ID，通常它是可选的
}

judgeResult, err := deer.Judge(judgeOptions)
```
评测结果信息的结构体定义如下:
```
type JudgeResult struct {
	JudgeResult int 			// 评测结果
	TimeUsed int				// 运行时最大时间使用 (ms)
	MemoryUsed int				// 运行时最大内存使用 (Kbytes)
	ReSignum int				// RE时系统信号
	SameLines int				// WA时返回数据正确行数（测试）
	TotalLines int				// WA时返回数据总行数（测试）
	SeInfo string				// SE时的信息
	CeInfo string				// CE时的信息
}
```

## 💡 特殊评测
特殊评测通常是对一般黑盒评测的一种补充，裁判程序通常是出题人针对于题目编写的一段判题代码编译后的程序。通过运行特殊评测，出题者可以灵活控制判题流程和对数据准确性的要求。Deer判题核心支持以下两种特判模式:

 - 结果检查模式
 - 交互判题模式
 
**结果检查模式** Deer判题核心将会先启动正常的目标程序运行流程，当目标程序运行结束并且没有任何错误的时候，将会启动特判程序。特判程序可以通过运行参数上得到的目标程序的输出文件位置，进行文件读写操作，判定内容等。

特判程序的运行参数（命令行）定义如下:
```
./checker [1] [2] [3]
```
[1]: 测试数据输入文件位置; [2]: 测试数据输出文件位置; [3]: 目标程序输出文件位置


**交互判题模式** Deer判题核心会同时启动特判程序和目标程序，并将特判程序的_stdout_重定向到目标程序的_stdin_，同时将特判程序的_stdin_重定向到目标程序 _stdout_。判题核心将以特判程序的退出时间为准，完成判题流程。

特判程序的运行参数（命令行）定义如下:
```
./checker [1] [2] [3]
```
[1]: 测试数据输入文件位置; [2]: 测试数据输出文件位置; [3] 特判程序输出文件位置

特判程序输出文件: 通常这个文件用于特判程序记录和目标程序的交流内容，也可以作为特判日志使用。

**退出代码**

特判程序需要将判定结果将以退出代码的形式告知判题机。这里有个C语言的特判程序示例：
```
#define RESULT_AC 0
#define RESULT_PE 1
#define RESULT_WA 4
#define RESULT_OLE 6
#define REQUIRE_DEFAULT_CHECKER 12

int main(int argc,char **argv) {
    // 你的判题代码.
    
    return RESULT_AC;
}
```
_REQUIRE_DEFAULT_CHECKER_：请求默认的文本检查器

通常情况下，特判程序不一定直接给予判题结果，它也可以用于对目标程序的输出内容进行处理。例如，很多时候浮点类型保留两位小数，由于IEEE 754的问题导致输出内容会和实际结果偏差0.01之类的情况，导致单纯的文本比对失败。这时候通过特判程序处理输出内容，就可以忽略这个问题。 **REQUIRE_DEFAULT_CHECKER**这个退出代码被返回的时候，判题程序将继续调用标准的文本比对程序，来给出AC、PE或WA的判定

## 🧬 编译提供程序

Deer判题内核为一下语言提供了编译提供程序:
```
GCC、GNU C++、Java、Python2、Python3、Golang、NodeJS、PHP、Ruby
```
当然你也可以根据你的需要自己编写提供程序，继承**CodeCompileProviderInterface**接口即可
```
type CodeCompileProviderInterface interface {

    // 初始化编译提供程序
    Init(code string, workDir string) error
    
    // 编译程序（需要初始化后方可使用）
    Compile() (result bool, errmsg string)
    
    // 获取编译后目标程序的执行参数
    GetRunArgs() (args []string)

    // 判断目标程序的stderr的输出内容是否存在编译错误信息，通常用于脚本语言的判定。
    // 如Python语言不需要编译，在执行脚本的时候如果遇到编译错误会返回SyntaxError信息之类的
    IsCompileError(remsg string) bool
    
    // 是否为实时编译的语言
    IsRealTime() bool
    
    // 是否已经编译完毕
    IsReady() bool

    /** 
     ** 私有方法
     **/

    // 初始化文件信息
    initFiles(codeExt string, programExt string) error
    
	// 执行系统调用
	shell(commands string) (success bool, errout string)
	// 保存代码内容到文件
	saveCode() error
	// 检查工作目录是否存在
	checkWorkDir() error
}

type CodeCompileProvider struct {
	CodeCompileProviderInterface
	
    codeContent string		            // 代码内容
	realTime bool			            // 是否为实时编译的语言
	isReady bool			            // 是否已经编译完毕
	codeFileName string                 // 目标程序源文件名
	codeFilePath string			        // 目标程序源文件路径
	programFileName string              // 目标程序文件名
	programFilePath string		        // 目标程序文件路径
	workDir string			            // 工作目录
}
```
  

## 🤝 鸣谢

首先感谢 [Loco's runner](https://github.com/dojiong/Lo-runner) 的作者，为本程序提供了黑盒评测的实现思路。

另外，感谢Wolf Zheng和Tosh Qiu提出的交互式评测的需求和基本工作流程的描述。

最后，感谢[北京师范大学(珠海校区)](http://www.bnuz.edu.cn)[信息技术学院](http://itc.bnuz.edu.cn)对WeJudge项目的支持，感谢[北师珠ACM协会](http://acm.bnuz.edu.cn)，感谢WeJudge团队每一位成员的付出。

## 🔗 相关链接

📃 我的博客：[https://www.lanrongqi.com](https://www.lanrongqi.com)

🖥️ WeJudge程序设计课程在线判题辅助教学平台：

[https://www.wejudge.net](https://www.wejudge.net) 

[https://oj.bnuz.edu.cn](https://oj.bnuz.edu.cn)

[WeJudge 1.0开源代码](https://github.com/LanceLRQ/wejudge)


**欢迎各位开发者使用和开发本程序，只要遵守咱们的GPLv3协议即可，使用过程中如果遇到什么问题，欢迎发Issue一起讨论哦！**