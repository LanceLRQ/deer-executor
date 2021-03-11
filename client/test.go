package client

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/sandbox/forkexec"
    "github.com/LanceLRQ/deer-common/sandbox/process"
    "github.com/urfave/cli/v2"
    "log"
    "os"
    "syscall"
)

func Test(c *cli.Context) error {
    var err error
    pid := 0
    ok := make(chan bool, 1)
    go func () {
        var infile *os.File
        var p *process.Process
        var ps *process.ProcessState
        infile, err = os.OpenFile("/Users/yiyiwukeji/github/deer-executor/data/problems/APlusB/0.in", os.O_RDONLY, 0)
        if err != nil {
            ok <- false
            return
        }
        p, err = process.StartProcess("/tmp/1dc8b300-821a-11eb-abb1-787b8ab7e6fa/cd19f388-9a22-41ac-98d0-2e1b17bba0bf", []string {
            "/tmp/1dc8b300-821a-11eb-abb1-787b8ab7e6fa/cd19f388-9a22-41ac-98d0-2e1b17bba0bf",
        }, &process.ProcAttr{
            Dir: "/tmp/1dc8b300-821a-11eb-abb1-787b8ab7e6fa",
            Env:   os.Environ(),
            Files: []interface{}{infile,  os.Stdout,  os.Stderr},
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
            ok <- false
            return
        }
        pid = p.Pid
        log.Println(pid)
        ps, err = p.Wait()
        if err != nil {
            ok <- false
            return
        }
        ws := ps.Sys().(syscall.WaitStatus)
        fmt.Printf("\n\nExitCode: %d, Signal: %d\n", ps.ExitCode(), ws.Signal())
        ok <- true
    }()
    select {
        case rel := <- ok:
            if rel {
                fmt.Println("OK")
            }
    }
    return err
}
