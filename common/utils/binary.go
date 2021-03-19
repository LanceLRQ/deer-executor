package utils

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/LanceLRQ/deer-executor/v2/common/constants"
	"github.com/LanceLRQ/deer-executor/v2/common/structs"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

func IsExecutableFile(filePath string) (bool, error) {
	fp, err := os.OpenFile(filePath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return false, errors.Errorf("open file error")
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

func GetCompiledBinaryFileName(typeName, moduleName string) string {
	prefix, ok := constants.TestlibBinaryPrefixs[typeName]
	if !ok {
		prefix = ""
	}
	return prefix + moduleName
}

// 根据配置文件将对应预编译文件转换成绝对路径
func GetCompiledBinaryFileAbsPath(typeName, moduleName, configDir string) (string, error) {
	targetName := GetCompiledBinaryFileName(typeName, moduleName)
	return filepath.Abs(path.Join(path.Join(configDir, "bin"), targetName))
}

// 解析generator脚本
func ParseGeneratorScript(script string) (string, []string, error) {
	vals := strings.Split(script, " ")
	if len(vals) <= 1 {
		return "", nil, errors.Errorf("generator calling script error")
	}
	return vals[0], vals[1:], nil
}

// 运行UnixShell，支持context
func RunUnixShell(options *structs.ShellOptions) (*structs.ShellResult, error) {
	fpath, err := exec.LookPath(options.Name)
	if err != nil {
		return nil, err
	}
	result := structs.ShellResult{}
	proc := exec.CommandContext(options.Context, fpath, options.Args...)
	var stderr, stdout bytes.Buffer

	if options.StdWriter != nil && options.StdWriter.Output != nil {
		proc.Stdout = options.StdWriter.Output
	} else {
		proc.Stdout = &stdout
	}
	if options.StdWriter != nil && options.StdWriter.Error != nil {
		proc.Stderr = options.StdWriter.Error
	} else {
		proc.Stderr = &stderr
	}

	if options.StdWriter != nil && options.StdWriter.Input != nil {
		proc.Stdin = options.StdWriter.Input
	} else {
		stdin, err := proc.StdinPipe()
		if err != nil {
			return nil, err
		}
		if options.OnStart != nil {
			err = options.OnStart(stdin)
			if err != nil {
				return nil, err
			}
		}
		_ = stdin.Close()
	}

	//err = proc.Run()
	if err := proc.Start(); err != nil {
		return nil, err
	}

	err = proc.Wait()

	if options.StdWriter == nil || options.StdWriter.Output == nil {
		result.Stdout = stdout.String()
	}
	if options.StdWriter == nil || options.StdWriter.Error == nil {
		result.Stderr = stderr.String()
	}
	result.ExitCode = proc.ProcessState.ExitCode()
	result.Signal = int(proc.ProcessState.Sys().(syscall.WaitStatus).Signal())
	if err != nil {
		result.Success = false
		result.ErrorMessage = err.Error()
		if serr := result.Stderr; serr == "" {
			result.Stderr += err.Error()
		}
		return &result, nil
	}
	result.Success = true
	return &result, nil
}

func CallGenerator(ctx context.Context, tc *structs.TestCase, configDir string) ([]byte, error) {
	name, args, err := ParseGeneratorScript(tc.Generator)
	if err != nil {
		return nil, err
	}
	gBin, err := GetCompiledBinaryFileAbsPath("generator", name, configDir)
	if err != nil {
		return nil, err
	}
	rel, err := RunUnixShell(&structs.ShellOptions{
		Context:   ctx,
		Name:      gBin,
		Args:      args,
		StdWriter: nil,
		OnStart:   nil,
	})
	if err != nil {
		return nil, err
	}
	if rel.Success {
		return []byte(rel.Stdout), nil
	} else {
		return nil, errors.Errorf("generator error")
	}
}

// 判断是否是题目包
func IsZipFile(filePath string) (bool, error) {
	fp, err := os.OpenFile(filePath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return false, errors.Errorf("open file error")
	}
	defer fp.Close()

	var magic uint32 = 0
	err = binary.Read(fp, binary.BigEndian, &magic)
	if err != nil {
		return false, err
	}
	return magic == constants.ZipArchiveMagicCode, nil
}

// 判断是否是题目包
func IsProblemPackage(filePath string) (bool, error) {
	fp, err := os.OpenFile(filePath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		return false, errors.Errorf("open file error")
	}
	defer fp.Close()

	var magic uint16 = 0
	err = binary.Read(fp, binary.BigEndian, &magic)
	if err != nil {
		return false, err
	}

	return magic == constants.ProblemPackageMagicCode, nil
}
