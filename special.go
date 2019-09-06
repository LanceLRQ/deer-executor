package deer_executor

import (
	"fmt"
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
				rst.JudgeResult = JUDGE_FLAG_SPJ_TIME_OUT
				return fmt.Errorf("special judger time limit exceed, unix singal: %d", sig)
			}
			rst.JudgeResult = JUDGE_FLAG_SPJ_ERROR
			return fmt.Errorf("special judger caused an error, unix singal: %d", sig)
		} else {
			rst.JudgeResult = JUDGE_FLAG_RE
		}
	} else {
		if status.Exited() {
			exitcode := status.ExitStatus()
			fmt.Printf("Special ExitCode: %d\n", exitcode)

			if exitcode == JUDGE_FLAG_AC || exitcode == JUDGE_FLAG_PE ||
				exitcode == JUDGE_FLAG_WA || exitcode == JUDGE_FLAG_OLE ||
				exitcode == JUDGE_FLAG_SPJ_REQUIRE_CHECK {
				rst.JudgeResult = exitcode
			} else {
				rst.JudgeResult = JUDGE_FLAG_SPJ_ERROR
				return fmt.Errorf("special judger return with a wrong exitcode: %d", exitcode)
			}
		}
	}
	return nil
}

func CustomChecker(options JudgeOption, result *JudgeResult, checkerPid chan uintptr) error {
	if runtime.GOOS != "linux" {
		result.JudgeResult = JUDGE_FLAG_SE
		result.SeInfo += "special judge can only be enable at linux.\n"
		return fmt.Errorf("special judge can only be enable at linux")
	}
	var (
		err, childErr error = nil, nil
		pid uintptr
		stdinFd, stdoutFd, stderrFd int
	)
	pid, err = forkProc()
	if err != nil {
		result.JudgeResult = JUDGE_FLAG_SE
		result.SeInfo += err.Error() + "\n"
		return  err
	}

	if pid == 0 {
		if options.SpecialJudge.RedirectStd {
			// Redirect testCaseIn to STDIN
			stdinFd, childErr = redirectFileDescriptor(syscall.Stdin, options.ProgramOut, os.O_RDONLY, 0)
			if childErr != nil {
				return childErr
			}
		}

		// Redirect userOut to STDOUT
		stdoutFd, childErr = redirectFileDescriptor(syscall.Stdout, options.SpecialJudge.Stdout, os.O_WRONLY | os.O_CREATE, 0644)
		if childErr != nil {
			return childErr
		}

		// Redirect programError to STDERR
		stderrFd, childErr = redirectFileDescriptor(syscall.Stderr, options.SpecialJudge.Stderr, os.O_WRONLY | os.O_CREATE, 0644)
		if childErr != nil {
			return childErr
		}

		tl, ml := SPECIAL_JUDGE_TIME_LIMIT, SPECIAL_JUDGE_MEMORY_LIMIT
		if options.SpecialJudge.TimeLimit > 0 { tl = options.SpecialJudge.TimeLimit }
		if options.SpecialJudge.MemoryLimit > 0 { tl = options.SpecialJudge.MemoryLimit  }

		// Set resource limit
		childErr = setLimit(tl, ml)
		if childErr != nil {
			return childErr
		}

		// Run Checker
		args := []string{ options.SpecialJudge.Checker, options.TestCaseIn, options.TestCaseOut, options.ProgramOut }
		childErr = syscall.Exec(options.SpecialJudge.Checker, args, nil)

		os.Exit(0)

	} else {
		if checkerPid != nil {
			checkerPid <- pid
		}
		err = waitCustomChecker(options, pid, result, false)
		if err != nil {
			result.JudgeResult = JUDGE_FLAG_SE
			result.SeInfo += err.Error() + "\n"
			return err
		}
		if childErr != nil {
			result.JudgeResult = JUDGE_FLAG_SE
			result.SeInfo += childErr.Error() + "\n"
			return childErr
		}
		syscall.Close(stdinFd)
		syscall.Close(stdoutFd)
		syscall.Close(stderrFd)
	}
	return err
}

func InteractiveChecker(options JudgeOption, result *JudgeResult, checkerPid, childPid chan uintptr) error {
	if runtime.GOOS != "linux" {
		result.JudgeResult = JUDGE_FLAG_SE
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
		result.JudgeResult = JUDGE_FLAG_SE
		result.SeInfo += err.Error() + "\n"
		return err
	}
	syscall.Pipe(fdtarget)
	if err != nil {
		result.JudgeResult = JUDGE_FLAG_SE
		result.SeInfo += err.Error() + "\n"
		return err
	}

	// Run Program
	pidProgram, err = forkProc()
	if err != nil {
		result.JudgeResult = JUDGE_FLAG_SE
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
		childErr = setLimit(options.TimeLimit, options.MemoryLimit)
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
		if childPid != nil {
			childPid <- pidProgram
		}
		// Run Judger
		pidJudger, judgerErr = forkProc()
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

			tl, ml := SPECIAL_JUDGE_TIME_LIMIT, SPECIAL_JUDGE_MEMORY_LIMIT
			if options.SpecialJudge.TimeLimit > 0 { tl = options.SpecialJudge.TimeLimit }
			if options.SpecialJudge.MemoryLimit > 0 { tl = options.SpecialJudge.MemoryLimit  }

			// Set resource limit
			childErr = setLimit(tl, ml)
			if childErr != nil {
				return childErr
			}

			// Run Judger
			args := []string{ options.SpecialJudge.Checker, options.TestCaseIn, options.TestCaseOut, options.ProgramOut }
			judgerErr = syscall.Exec(options.SpecialJudge.Checker, args, nil)

			return childErr

		} else {
			if checkerPid != nil {
				checkerPid <- pidJudger
			}
			err = waitCustomChecker(options, pidJudger, result, true)
			if err != nil {
				result.JudgeResult = JUDGE_FLAG_SE
				result.SeInfo += err.Error() + "\n"
				return err
			}
			if judgerErr != nil {
				result.JudgeResult = JUDGE_FLAG_SE
				result.SeInfo += judgerErr.Error() + "\n"
				return judgerErr
			}
			if childErr != nil {
				result.JudgeResult = JUDGE_FLAG_SE
				result.SeInfo += childErr.Error() + "\n"
				return childErr
			}
		}
	}
	return nil
}