package run

import (
	"github.com/LanceLRQ/deer-executor/v2/common/constants"
	"github.com/LanceLRQ/deer-executor/v2/common/persistence/problems"
	"github.com/LanceLRQ/deer-executor/v2/common/provider"
	"github.com/LanceLRQ/deer-executor/v2/common/utils"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"os"
)

func loadSystemConfiguration() error {
	// 载入默认配置
	err := provider.PlaceCompilerCommands("./compilers.json")
	if err != nil {
		return err
	}
	err = constants.PlaceMemorySizeForJIT("./jit_memory.json")
	if err != nil {
		return err
	}
	return nil
}

func loadProblemConfiguration(configFile string, workDir string) (string, bool, string, error) {
	_, err := os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		return "", false, "", errors.Errorf("problem config file (%s) not found", configFile)
	}
	isDeerPack, err := utils.IsProblemPackage(configFile)
	if err != nil {
		return "", false, "", err
	}
	isZipPack, err := utils.IsZipFile(configFile)
	if err != nil {
		return "", false, "", err
	}
	autoRemoveWorkDir := false
	// 如果是题目包文件，进行解包，并返回解包后的配置文件位置
	if isDeerPack || isZipPack {
		if workDir == "" {
			workDir = "/tmp/" + uuid.NewV4().String()
			autoRemoveWorkDir = true
		}
		if info, err := os.Stat(workDir); os.IsNotExist(err) {
			err = os.MkdirAll(workDir, 0755)
			if err != nil {
				return "", false, "", err
			}
		} else if !info.IsDir() {
			return "", false, "", errors.Errorf("work dir path cannot be a file path")
		}
		if isDeerPack {
			_, configFile, err = problems.ReadProblemInfo(configFile, true, true, workDir)
		} else {
			_, configFile, err = problems.ReadProblemInfoZip(configFile, true, true, workDir)
		}
		if err != nil {
			return "", false, "", err
		}
	}
	// 如果不是题包文件，返回的是配置文件的路径和工作目录，
	return configFile, autoRemoveWorkDir, workDir, nil
}
