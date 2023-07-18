package client

import (
	"github.com/LanceLRQ/deer-executor/v3/executor/constants"
	"github.com/LanceLRQ/deer-executor/v3/executor/provider"
)

func LoadSystemConfiguration() error {
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
