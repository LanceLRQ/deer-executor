# 沙箱工具

这是一个基于Go语言自带的syscall库和os库里搬运出来进程调度工具，未来它的目标是实现一个更加复杂的沙箱。

本工具基于最新的Go 1.16源代码提取并编译，经测试在1.14版本及以上可以编译通过，其他低版本没有测试。

本工具**仅在**Linux(amd64)和MacOS(amd64,arm64)操作系统下测试编译并通过，请注意平台兼容性。

## 官方文档
- `forkexec` 对应 `syscall` [https://pkg.go.dev/syscall](https://pkg.go.dev/syscall)
- `process` 对应 `os` [https://pkg.go.dev/os](https://pkg.go.dev/os)

## 改动
- 原`syscall`包里的 `ProcAttr` 和 `SysProcAttr`抽出到`forkexec`包，原`os`包里的`ProcAttr`抽出到`process`包
- `forkexec`包里的`SysProcAttr`增加了`ExecRLimit`定义。注意，除栈限制外，只要是值为0就表示不做限制。
```golang
type ExecRLimit struct {
    TimeLimit     int           // 时间限制 (ms)
    RealTimeLimit int           // 真实时间限制 (ms, 触发SIGALRM)
    MemoryLimit   int           // 内存限制 (KB)
    FileSizeLimit int           // 文件读写限制 (B)
    StackLimit    int           // 栈大小限制 (KB，0表示用内存限制的值，-1表示不限制)
}
```

- `process`包的`ProcAttr`修改了`Files`的类型为`interface{}`
```golang
type ProcAttr struct {
    Dir string
    Env []string
    Files []interface{}         // 支持传入*os.File或者是uintptr（即使用syscall.Pipe()获取的管道文件描述符号）
    Sys *forkexec.SysProcAttr   // 换成修改过的SysProcAttr，支持setrlimit操作
}
```

## 使用
- 简单的运行一个程序并把子程序的内容输出到stdout
```golang
func run() {
    p, err := process.StartProcess("./test", nil, &process.ProcAttr{
        Env:   os.Environ(),
        Files: []interface{}{os.Stdin,  os.Stdout,  os.Stderr},
        Sys: &forkexec.SysProcAttr {
            Rlimit: forkexec.ExecRLimit {
                TimeLimit: 1000,
                RealTimeLimit: 2000,
                MemoryLimit: 128000,
                StackLimit: 128000,
                FileSizeLimit: 1024 * 1024 * 50,
            },
        },
    })
    if err != nil {
        panic(err)
    }
    ps, err := p.Wait()
    if err != nil {
        panic(err)
    }
    ws := ps.Sys().(syscall.WaitStatus)
    fmt.Printf("\n\nExitCode: %d, Signal: %d\n", ps.ExitCode(), ws.Signal())
}
```
- 运行一个交互程序，此时它们的输入输出使用管道进行连接。
```golang
func program(stdin, stdout uintptr) {
    p, err := process.StartProcess("./target", nil, &process.ProcAttr{
        Env:   os.Environ(),
        Files: []interface{}{stdin,  stdout,  os.Stderr},
        Sys: &forkexec.SysProcAttr {
            Rlimit: forkexec.ExecRLimit {
                TimeLimit: 1000,
                RealTimeLimit: 2000,
                MemoryLimit: 128000,
                StackLimit: 65500,
                FileSizeLimit: 1024 * 1024 * 50,
            },
        },
    })
    if err != nil {
        panic(err)
    }
    ps, err := p.Wait()
    if err != nil {
        panic(err)
    }
    ws := ps.Sys().(syscall.WaitStatus)
    fmt.Printf("\n\nExitCode: %d, Signal: %d\n", ps.ExitCode(), ws.Signal())
}

func judgement(stdin, stdout uintptr) {
    p, err := process.StartProcess("./judgement", []string {
        "./judgement",
        "./0.in",
        "./0.out",
        "./0.err",
        "./0.log",
    }, &process.ProcAttr{
        Env:   os.Environ(),
        Files: []interface{}{stdin,  stdout,  os.Stderr},
        Sys: &forkexec.SysProcAttr {
            Rlimit: forkexec.ExecRLimit {
                TimeLimit: 1000,
                RealTimeLimit: 2000,
                MemoryLimit: 128000,
                StackLimit: 65500,
                FileSizeLimit: 1024 * 1024 * 50,
            },
        },
    })
    if err != nil {
        panic(err)
    }
    ps, err := p.Wait()
    if err != nil {
        panic(err)
    }
    ws := ps.Sys().(syscall.WaitStatus)
    fmt.Printf("\n\nExitCode: %d, Signal: %d\n", ps.ExitCode(), ws.Signal())
}

func run() {
    fdjudge, err := forkexec.GetPipe()
    if err != nil {
        panic(err)
    }
    
    fdtarget, err := forkexec.GetPipe()
    if err != nil {
        panic(err)
    }

    go program(fdtarget[0], fdjudge[1])
    go judgement(fdjudge[0], fdtarget[1])

    // 这里演示，只是简单的睡一下，一般会用管道去监听协程，并注意进行超时处理。
    time.Sleep(10 * time.Second)
}
```

## 未来计划

- ptrace + seccomp
- windows支持
