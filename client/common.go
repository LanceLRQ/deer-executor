package client

import (
	"github.com/LanceLRQ/deer-executor/v2/common/constants"
	"github.com/LanceLRQ/deer-executor/v2/common/provider"
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
