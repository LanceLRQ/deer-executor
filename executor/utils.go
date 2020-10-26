/* Deer executor
 * (C) 2019-Now LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package executor

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"runtime"
	"syscall"
	"unsafe"
)

// 定义ITimer的常量，命名规则遵循Linux的原始设定
const (
	ITimerReal 		= 0
	ITimerVirtual 	= 1
	ITimerVProf 	= 2
)

// 定义公共环境变量
var CommonEnvs = []string{ "PYTHONIOENCODING=utf-8" }

type ITimerVal struct  {
	ItInterval TimeVal
	ItValue TimeVal
}

type TimeVal struct {
	TvSec uint64
	TvUsec uint64
}

type RLimit struct {
	Which 	int
	RLim	syscall.Rlimit
}

// 打开并获取文件的描述符
func getFileDescriptor(path string, flag int, perm uint32) (fd int, err error) {
	var filed = 0
	_, errMsg := os.Stat(path)
	if errMsg != nil {
		if os.IsNotExist(err) {
			return 0, errMsg
		}
	}
	filed, errMsg = syscall.Open(path, flag, perm)
	return filed, nil
}

// 重映射文件描述符
func redirectFileDescriptor(to int, path string, flag int, perm uint32) (fd int, err error) {
	fd, errMsg := getFileDescriptor(path, flag, perm)
	if errMsg == nil {
		errMsg = syscall.Dup2(fd, to)
		if errMsg != nil {
			syscall.Close(fd)
			return -1, errMsg
		}
		return fd, nil
	} else {
		return -1, errMsg
	}
}

// fork调用
func forkProc() (pid uintptr, err error) {
	r1, r2, errMsg := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	darwin := runtime.GOOS == "darwin"
	if errMsg != 0 {
		return 0, fmt.Errorf("system call: fork(); error: %s", errMsg)
	}
	if darwin {
		if r2 == 1 {
			pid = 0
		} else {
			pid = r1
		}
	} else {
		if r1 == 0 && r2 == 0 {
			pid = 0
		} else {
			pid = r1
		}
	}
	return pid, nil
}

func getPipe() ([]int, error) {
	var pipe = []int{0, 0}
	err := syscall.Pipe(pipe)
	if err != nil {
		return nil, err
	}
	return pipe, nil
}

// 打开文件并获取描述符 (强制文件检查)
func OpenFile(filePath string, flag int, perm os.FileMode) (*os.File, error) {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file (%s) not exists", filePath)
		} else {
			return nil, fmt.Errorf("open file (%s) error: %s", filePath, err.Error())
		}
	} else {
		if fp, err := os.OpenFile(filePath, flag, perm); err != nil {
			return nil, fmt.Errorf("open file (%s) error: %s", filePath, err.Error())
		} else {
			return fp, nil
		}
	}
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func Max32(a, b int) int {
	if a > b { return a } else { return b }
}

func getRLimitEntity(cur, max uint64) syscall.Rlimit {
	return syscall.Rlimit{Cur: cur, Max: max}
}

// 硬件计时器
func setHardTimer(realTimeLimit int) error {
	var prealt ITimerVal
	prealt.ItInterval.TvSec = uint64(math.Floor(float64(realTimeLimit) / 1000.0))
	prealt.ItInterval.TvUsec = uint64(realTimeLimit % 1000 * 1000)
	prealt.ItValue.TvSec = prealt.ItInterval.TvSec
	prealt.ItValue.TvUsec = prealt.ItInterval.TvUsec
	_, _, err := syscall.RawSyscall(syscall.SYS_SETITIMER, ITimerReal, uintptr(unsafe.Pointer(&prealt)), 0)
	if err != 0 {
		return fmt.Errorf("system call setitimer() error: %s", err)
	}
	return nil
}

// 设置资源限制 (setrlimit)
func setLimit(timeLimit, memoryLimit , realTimeLimit, fileSizeLimit int) error {

	// Set stack limit
	stack := uint64(memoryLimit * 1024)
	if runtime.GOOS == "darwin" {  // WTF?! >= 65mb caused an operation not permitted!
		stack = uint64(65500 * 1024)
	}

	rlimits := []RLimit{
		// Set time limit: RLIMIT_CPU
		{
			Which: syscall.RLIMIT_CPU,
			RLim: getRLimitEntity(
				uint64(math.Ceil(float64(timeLimit) / 1000.0)),
				uint64(math.Ceil(float64(timeLimit) / 1000.0)),
			),
		},
		// Set memory limit: RLIMIT_DATA
		{
			Which: syscall.RLIMIT_DATA,
			RLim: getRLimitEntity(
				uint64(memoryLimit * 1024),
				uint64(memoryLimit * 1024),
			),
		},
		// Set memory limit: RLIMIT_AS
		{
			Which: syscall.RLIMIT_AS,
			RLim: getRLimitEntity(
				uint64(memoryLimit * 1024 * 2),
				uint64(memoryLimit * 1024 * 2 + 1024),
			),
		},
		// Set stack limit
		{
			Which: syscall.RLIMIT_STACK,
			RLim: getRLimitEntity(
				stack,
				stack,
			),
		},
		// Set file size limit: RLIMIT_FSIZE
		{
			Which: syscall.RLIMIT_FSIZE,
			RLim: getRLimitEntity(
				uint64(fileSizeLimit),
				uint64(fileSizeLimit),
			),
		},

	}

	for _, rlimit := range rlimits {
		err := syscall.Setrlimit(rlimit.Which, &rlimit.RLim)
		if err != nil {
			return fmt.Errorf("setrlimit(%d) error: %s", rlimit.Which, err)
		}
	}

	// Set time limit: setITimer
	if realTimeLimit > 0 {
		err := setHardTimer(realTimeLimit)
		if err != nil {
			return err
		}
	}

	return nil
}

// 文件读写(有重试次数，checker专用)
func readFileWithTry(filePath string, name string, tryOnFailed int) ([]byte, string, error) {
	errCnt, errText := 0, ""
	var err error
	for errCnt < tryOnFailed {
		fp, err := OpenFile(filePath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
		if err != nil {
			errText = err.Error()
			errCnt++
			continue
		}
		data, err := ioutil.ReadAll(fp)
		if err != nil {
			_ = fp.Close()
			errText = fmt.Sprintf("Read file(%s) i/o error: %s", name, err.Error())
			errCnt++
			continue
		}
		_ = fp.Close()
		return data, errText, nil
	}
	return nil, errText, err
}

func ObjectToJSONStringFormatted(conf interface{}) string {
	b, err := json.Marshal(conf)
	if err != nil {
		return fmt.Sprintf("%+v", conf)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", conf)
	}
	return out.String()
}

func ObjectToJSONByte(obj interface{}) []byte {
	b, err := json.Marshal(obj)
	if err != nil {
		return []byte("{}")
	}
	return b
}

func ObjectToJSONString(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return "{}"
	} else {
		return string(b)
	}
}

func JSONStringObject(jsonStr string, obj interface{}) bool {
	return JSONBytesObject([]byte(jsonStr), obj)
}

func JSONBytesObject(jsonBytes []byte, obj interface{}) bool {
	err := json.Unmarshal(jsonBytes, &obj)
	if err != nil {
		return false
	} else {
		return true
	}
}

func IsExecutableFile (filePath string) (bool, error) {
	fp, err := os.OpenFile(filePath, os.O_RDONLY | syscall.O_NONBLOCK, 0)
	if err != nil {
		return false, fmt.Errorf("open file error")
	}
	defer fp.Close()

	var magic uint32 = 0
	err = binary.Read(fp, binary.BigEndian, &magic)
	if err != nil {
		return false, err
	}

	isExec := false
	if runtime.GOOS == "darwin" {
		isExec = magic == 0xCFFAEDFE || magic == 0xCEFAEDFE || magic == 0xFEEDFACF || magic == 0xFEEDFACE
	} else if runtime.GOOS == "linux" {
		isExec = magic == 0x7F454C46
	}
	return isExec, nil
}

// 检查配置文件里的所有文件是否存在
func CheckRequireFilesExists(session *JudgeSession) error {
	var err error
	// 检查特判程序是否存在
	if session.SpecialJudge.Mode != 0 {
		_, err = os.Stat(path.Join(session.ConfigDir, session.SpecialJudge.Checker))
		if os.IsNotExist(err) {
			return fmt.Errorf("special judge checker file (%s) not exists", session.SpecialJudge.Checker)
		}
	}
	// 检查每个test case里的文件是否存在
	// 新版判题机要求无论有没有数据，都要有对应的输入输出文件。
	for i := 0; i < len(session.TestCases); i++ {
		tcase := session.TestCases[i]
		_, err = os.Stat(path.Join(session.ConfigDir, tcase.TestCaseIn))
		if os.IsNotExist(err) {
			return fmt.Errorf("test case (%s) input file (%s) not exists", tcase.Id, tcase.TestCaseIn)
		}
		_, err = os.Stat(path.Join(session.ConfigDir, tcase.TestCaseOut))
		if os.IsNotExist(err) {
			return fmt.Errorf("test case (%s) output file (%s) not exists", tcase.Id, tcase.TestCaseOut)
		}
 	}
	return nil
}