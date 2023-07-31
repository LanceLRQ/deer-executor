package run

import (
	"github.com/LanceLRQ/deer-executor/v3/executor/persistence"
	"github.com/LanceLRQ/deer-executor/v3/executor/utils"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"os"
	"path"
)

func loadProblemConfiguration(configFile string, workDir string) (string, bool, string, error) {
	_, err := os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		return "", false, "", errors.Errorf("problem config file (%s) not found", configFile)
	}
	isDeerPack, packageType, err := utils.IsDeerPackage(configFile)
	if err != nil {
		return "", false, "", err
	}
	autoRemoveWorkDir := false
	// If it is a deer-package file, unpack it and return the location of the extracted configuration file.
	if isDeerPack && packageType == persistence.PackageTypeProblem {
		pack, err := persistence.ParseProblemPackageFile(configFile, true)
		if err != nil {
			return "", false, "", err
		}
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

		configFile = path.Join(workDir, "problem.json")
		// unpack
		err = pack.UnpackProblemProject(workDir)
		if err != nil {
			return "", false, "", err
		}
	}
	// If not, return the path to the configuration file and the working directory.
	return configFile, autoRemoveWorkDir, workDir, nil
}
