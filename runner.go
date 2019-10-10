/* Deer executor
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package deer_executor

import (
	"fmt"
	"os"
	"syscall"
)


func waitExit(options JudgeOption, pid uintptr, rst *JudgeResult) error {
	var (
		status syscall.WaitStatus
		ru syscall.Rusage
	)
	_, errMsg := syscall.Wait4(int(pid), &status, syscall.WUNTRACED, &ru)
	if errMsg != nil {
		return errMsg
	}

	rst.TimeUsed = int(ru.Utime.Sec * 1000 + int64(ru.Utime.Usec) / 1000 + ru.Stime.Sec * 1000 + int64(ru.Stime.Usec) / 1000)
	rst.MemoryUsed = int(ru.Minflt * int64(syscall.Getpagesize() / 1024 ))

	if status.Signaled() {
		sig := status.Signal()
		rst.ReSignum = int(sig)
		if sig == syscall.SIGSEGV {
			// MLE or RE can also get SIGSEGV signal.
			if rst.MemoryUsed > options.MemoryLimit {
				rst.JudgeResult = JUDGE_FLAG_MLE
			} else {
				rst.JudgeResult = JUDGE_FLAG_RE
			}
		} else if sig == syscall.SIGXFSZ {
			// SIGXFSZ signal means OLE
			rst.JudgeResult = JUDGE_FLAG_OLE
		} else if sig == syscall.SIGALRM || sig == syscall.SIGVTALRM || sig == syscall.SIGXCPU {
			// Normal TLE signal
			rst.JudgeResult = JUDGE_FLAG_TLE
		} else if sig == syscall.SIGKILL {
			// Sometimes MLE might get SIGKILL signal.
			// So if real time used lower than TIME_LIMIT - 100, it might be a TLE error.
			if rst.TimeUsed > (options.TimeLimit - 100) {
				rst.JudgeResult = JUDGE_FLAG_TLE
			} else {
				rst.JudgeResult = JUDGE_FLAG_MLE
			}
		} else {
			// Otherwise, called runtime error.
			rst.JudgeResult = JUDGE_FLAG_RE
		}
	} else {
		// Sometimes setrlimit doesn't work accurately.
		if rst.TimeUsed > options.TimeLimit {
			rst.JudgeResult = JUDGE_FLAG_MLE
		} else if rst.MemoryUsed > options.MemoryLimit {
			rst.JudgeResult = JUDGE_FLAG_MLE
		} else {
			rst.JudgeResult = JUDGE_FLAG_AC
		}
	}
	return nil
}

func RunProgram(options JudgeOption, result *JudgeResult, msg chan string) error {

	var (
		err, childErr error = nil, nil
		pid uintptr
		stdinFd ,stdoutFd, stderrFd int
	)

	// Fork a new process
	pid, err = forkProc()
	if err != nil {
		result.JudgeResult = JUDGE_FLAG_SE
		result.SeInfo += err.Error() + "\n"
		return err
	}

	if pid == 0 {
		// child process: set limit & execute target program.

		// Redirect testCaseIn to STDIN
		stdinFd, childErr = redirectFileDescriptor(syscall.Stdin, options.TestCaseIn, os.O_RDONLY, 0)
		if childErr != nil {
			return childErr
		}

		// Redirect userOut to STDOUT
		stdoutFd, childErr = redirectFileDescriptor(syscall.Stdout, options.ProgramOut, os.O_WRONLY | os.O_CREATE, 0644)
		if childErr != nil {
			return childErr
		}

		// Redirect programError to STDERR
		stderrFd, childErr = redirectFileDescriptor(syscall.Stderr, options.ProgramError, os.O_WRONLY | os.O_CREATE, 0644)
		if childErr != nil {
			return childErr
		}

		// Set UID
		if options.Uid > -1 {
			childErr = syscall.Setuid(options.Uid)
			if childErr != nil {
				return childErr
			}
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

		return childErr		// In general, it won't be run.

	} else {
		if msg != nil {
			msg <- fmt.Sprintf("pid:program:%d", pid)
		}
		// paren process: wait for child process end.
		err = waitExit(options, pid, result)
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
		if msg != nil {
			msg <- "done"
		}
	}

	return nil
}
