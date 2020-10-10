package executor

import (
	"fmt"
	"github.com/docker/docker/pkg/reexec"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func runTargetProgram(options JudgeOptions) (*exec.Cmd, error) {
	payload := ObjectToJSONString(options)
	target := reexec.Command("targetProgram", payload)

	tcin, err := OpenFile(options.TestCaseIn, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil { return target, err }
	target.Stdin = tcin

	pout, err := os.OpenFile(options.ProgramOut, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil { return target, err }
	target.Stdout = pout

	perr, err := os.OpenFile(options.ProgramError, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil { return target, err }
	target.Stderr = perr

	if err := target.Start(); err != nil {
		return target, fmt.Errorf("failed to run target program: %s", err)
	}
	if err := target.Wait(); err != nil {
		return target, fmt.Errorf("failed to wait command: %s", err)
	}
	return target, nil
}

// 基于JudgeOptions进行评测调度
func Run(options JudgeOptions, result *JudgeResult) error {
	target, err := runTargetProgram(options)
	if err != nil {
		return err
	}
	if target.ProcessState.Exited() {
		fmt.Println("Finish!")
	}
	return nil
}

// 目标程序子进程
func RunTargetProgramProcess() {
	payload := os.Args[1]
	options := JudgeOptions{}
	if !JSONStringObject(payload, &options) {
		log.Fatalf("parse judge options error")
		return
	}

	// Set UID
	if options.Uid > -1 {
		err := syscall.Setuid(options.Uid)
		if err != nil {
			log.Fatalf("set resource limit error: %s", err.Error())
			return
		}
	}
	// Set Resource Limit
	err := setLimit(options.TimeLimit, options.MemoryLimit, options.RealTimeLimit)
	if err != nil {
		log.Fatalf("set resource limit error: %s", err.Error())
		return
	}
	// Run Program
	if len(options.Commands) > 1 {
		_ = syscall.Exec(options.Commands[0], options.Commands[1:], CommonEnvs)
	} else {
		_ = syscall.Exec(options.Commands[0], nil, CommonEnvs)
	}
}

// 特判程序子进程
func RunSpecialJudgeProgramProcess() {

}