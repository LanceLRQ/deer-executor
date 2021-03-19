# 公共功能库

这里存放公共函数和结构定义等代码，为后续Polygon的接入和开发做准备

## 目录结构
```
.
├── constants             常量库
│   ├── executor.go         判题机的常量定义，如判题结果、语言等
│   ├── persistence.go      持久化模块的文件魔数常量
│   └── testlib.go          Testlib所需要用到的常量定义
├── docs                  文档
│   └── testlib.md          Testlib说明
├── errors                错误定义（暂时没用到）
├── logger                评测日志工具
├── persistence           持久化
│   ├── judge_result        评测结果持久化功能
│   ├── problems            题目包功能
│   ├── struct.go           公共结构定义
│   └── utils.go            公共工具
├── provider              编译提供程序
│   ├── main.go             编译提供程序公共定义
│   ├── gcc.go              C语言Provider实现等
│   ├── ...
├── sandbox               沙箱工具（具体使用看子目录下的README）
│   ├── forkexec            原syscall包中关于forkExec和startProcess的内容
│   └── process             原os包中关于Process的内容
├── structs               公共结构体定义
│   ├── binary.go           Shell、CMD运行相关
│   ├── judge.go            评测配置相关
│   └── problems.go         题目设置相关
└── utils                 公共函数库
    ├── binary.go           Shell、CMD调度相关，用了os/exec包。还有一些二进制文件判断的函数。
    ├── common.go           一些公共函数
    ├── json.go             JSON处理相关函数
    ├── session.go          创建临时会话目录函数
    └── xml.go              XML处理相关函数
```
