/* Deer executor
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package obsolete

import (
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	"os"
	"runtime"
	"syscall"
)

func waitCustomChecker(options JudgeOption, pid uintptr, rst *JudgeResult, isInteractive bool) error {
	var (
		status syscall.WaitStatus
		ru syscall.Rusage
	)
	_, err := syscall.Wait4(int(pid), &status, syscall.WUNTRACED, &ru)
	if err != nil {
		return err
	}
	if status.Signaled() {
		sig := status.Signal()
		if !isInteractive {
			if sig == syscall.SIGXCPU || sig == syscall.SIGALRM {
				rst.JudgeResult = JudgeFlagSpecialJudgeTimeout
				return fmt.Errorf("special judger time limit exceed, unix singal: %d", sig)
			}
			rst.JudgeResult = JudgeFlagSpecialJudgeError
			return fmt.Errorf("special judger caused an error, unix singal: %d", sig)
		} else {
			rst.JudgeResult = JudgeFlagRE
		}
	} else {
		if status.Exited() {
			exitcode := status.ExitStatus()
			fmt.Printf("Special ExitCode: %d\n", exitcode)

			if exitcode == JudgeFlagAC || exitcode == JudgeFlagPE ||
				exitcode == JudgeFlagWA || exitcode == JudgeFlagOLE ||
				exitcode == JudgeFlagSpecialJudgeRequireChecker {
				rst.JudgeResult = exitcode
			} else {
				rst.JudgeResult = JudgeFlagSpecialJudgeError
				return fmt.Errorf("special judger return with a wrong exitcode: %d", exitcode)
			}
		}
	}
	return nil
}

func CustomChecker(options JudgeOption, result *JudgeResult, msg chan string) error {
	if runtime.GOOS != "linux" {
		result.JudgeResult = JudgeFlagSE
		result.SeInfo += "special judge can only be enable at linux.\n"
		return fmt.Errorf("special judge can only be enable at linux")
	}
	var (
		err, childErr error = nil, nil
		pid uintptr
		stdinFd, stdoutFd, stderrFd int
	)
	pid, err = executor.forkProc()
	if err != nil {
		result.JudgeResult = JudgeFlagSE
		result.SeInfo += err.Error() + "\n"
		return  err
	}

	if pid == 0 {
		if options.SpecialJudge.RedirectStd {
			// Redirect testCaseIn to STDIN
			stdinFd, childErr = executor.redirectFileDescriptor(syscall.Stdin, options.ProgramOut, os.O_RDONLY, 0)
			if childErr != nil {
				return childErr
			}
		}

		// Redirect userOut to STDOUT
		stdoutFd, childErr = executor.redirectFileDescriptor(syscall.Stdout, options.SpecialJudge.Stdout, os.O_WRONLY | os.O_CREATE, 0644)
		if childErr != nil {
			return childErr
		}

		// Redirect programError to STDERR
		stderrFd, childErr = executor.redirectFileDescriptor(syscall.Stderr, options.SpecialJudge.Stderr, os.O_WRONLY | os.O_CREATE, 0644)
		if childErr != nil {
			return childErr
		}

		tl, ml := SpecialJudgeTimeLimit, SpecialJudgeMemoryLimit
		if options.SpecialJudge.TimeLimit > 0 { tl = options.SpecialJudge.TimeLimit }
		if options.SpecialJudge.MemoryLimit > 0 { tl = options.SpecialJudge.MemoryLimit  }

		// Set resource limit
		childErr = executor.setLimit(tl, ml, tl)
		if childErr != nil {
			return childErr
		}

		// Run Checker
		args := []string{ options.SpecialJudge.Checker, options.TestCaseIn, options.TestCaseOut, options.ProgramOut }
		childErr = syscall.Exec(options.SpecialJudge.Checker, args, nil)

		os.Exit(0)

	} else {
		if msg != nil {
			msg <- fmt.Sprintf("pid:program:%d", pid)
		}
		err = waitCustomChecker(options, pid, result, false)
		if err != nil {
			result.JudgeResult = JudgeFlagSE
			result.SeInfo += err.Error() + "\n"
			return err
		}
		//if childErr != nil {
		//	result.JudgeResult = JudgeFlagSE
		//	result.SeInfo += childErr.Error() + "\n"
		//	return childErr
		//}
		_ = syscall.Close(stdinFd)
		_ = syscall.Close(stdoutFd)
		_ = syscall.Close(stderrFd)
		if msg != nil {
			msg <- "done"
		}
	}
	return err
}

func InteractiveChecker(options JudgeOption, result *JudgeResult, msg chan string) error {
	if runtime.GOOS != "linux" {
		result.JudgeResult = JudgeFlagSE
		result.SeInfo += "interactive special judge can only be enable at linux.\n"
		return fmt.Errorf("interactive special judge can only be enable at linux")
	}
	var (
		err, childErr, judgerErr error = nil, nil, nil
		pidJudger, pidProgram uintptr
		fdjudger, fdtarget []int = []int{0, 0}, []int{0, 0}
	)

	syscall.Pipe(fdjudger)
	if err != nil {
		result.JudgeResult = JudgeFlagSE
		result.SeInfo += err.Error() + "\n"
		return err
	}
	syscall.Pipe(fdtarget)
	if err != nil {
		result.JudgeResult = JudgeFlagSE
		result.SeInfo += err.Error() + "\n"
		return err
	}

	// Run Program
	pidProgram, err = executor.forkProc()
	if err != nil {
		result.JudgeResult = JudgeFlagSE
		result.SeInfo += err.Error() + "\n"
		return err
	}

	fmt.Println(fdtarget[0], fdtarget[1])
	fmt.Println(fdjudger[0], fdjudger[1])

	if pidProgram == 0 {

		// Direct Program's Pipe[Read] to Stdin
		childErr = syscall.Dup2(fdtarget[0], syscall.Stdin)
		if childErr != nil {
			return childErr
		}
		// Direct Judger's Pipe[Write] to Stdout
		childErr = syscall.Dup2(fdjudger[1], syscall.Stdout)
		if childErr != nil {
			return childErr
		}

		// Set resource limit
		childErr = executor.setLimit(options.TimeLimit, options.MemoryLimit, options.TimeLimit)
		if childErr != nil {
			return childErr
		}
		// Run Program
		if len(options.Commands) > 1 {
			childErr = syscall.Exec(options.Commands[0], options.Commands[1:], nil)
		} else {
			childErr = syscall.Exec(options.Commands[0], nil, nil)
		}
		return childErr

	} else {
		if msg != nil {
			msg <- fmt.Sprintf("pid:program:%d", pidProgram)
		}
		// Run Judger
		pidJudger, judgerErr = executor.forkProc()
		if judgerErr != nil {
			return judgerErr
		}

		if pidJudger == 0 {
			// Direct Judger's Pipe[Read] to Stdout
			judgerErr = syscall.Dup2(fdjudger[0], syscall.Stdin)
			if judgerErr != nil {
				return judgerErr
			}
			// Direct Program's Pipe[Write] to Stdin
			judgerErr = syscall.Dup2(fdtarget[1], syscall.Stdout)
			if judgerErr != nil {
				return judgerErr
			}

			tl, ml := SpecialJudgeTimeLimit, SpecialJudgeMemoryLimit
			if options.SpecialJudge.TimeLimit > 0 { tl = options.SpecialJudge.TimeLimit }
			if options.SpecialJudge.MemoryLimit > 0 { tl = options.SpecialJudge.MemoryLimit  }

			// Set resource limit
			childErr = executor.setLimit(tl, ml, options.TimeLimit)
			if childErr != nil {
				return childErr
			}

			// Run Judger
			args := []string{ options.SpecialJudge.Checker, options.TestCaseIn, options.TestCaseOut, options.ProgramOut }
			judgerErr = syscall.Exec(options.SpecialJudge.Checker, args, nil)

			return childErr

		} else {
			if msg != nil {
				msg <- fmt.Sprintf("pid:checker:%d", pidProgram)
			}
			err = waitCustomChecker(options, pidJudger, result, true)
			if err != nil {
				result.JudgeResult = JudgeFlagSE
				result.SeInfo += err.Error() + "\n"
				return err
			}
			//if judgerErr != nil {
			//	result.JudgeResult = JudgeFlagSE
			//	result.SeInfo += judgerErr.Error() + "\n"
			//	return judgerErr
			//}
			//if childErr != nil {
			//	result.JudgeResult = JudgeFlagSE
			//	result.SeInfo += childErr.Error() + "\n"
			//	return childErr
			//}
			if msg != nil {
				msg <- "done"
			}
		}
	}
	return nil
}