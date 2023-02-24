package server_config

import "github.com/gookit/config/v2"

type GRPCConfigDefinition struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JudgementConfigDefinition struct {
	ProblemRoot       string `mapstructure:"problem_root"`        // 题目根目录
	SystemLibraryRoot string `mapstructure:"system_library_root"` // testlib.h等的库目录
	SessionRoot       string `mapstructure:"session_root"`        // 评测会话目录
}

var GRPCConfig GRPCConfigDefinition
var JudgementConfig JudgementConfigDefinition

func LoadGlobalConf() error {
	err := config.BindStruct("grpc", &GRPCConfig)
	if err != nil {
		return err
	}
	err = config.BindStruct("judgement", &JudgementConfig)
	if err != nil {
		return err
	}
	return nil
}
