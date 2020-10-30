package problem

import (
    "fmt"
    "github.com/LanceLRQ/deer-common/structs"
    "os"
    "path"
)

func compileCodeFile () {

}


// 编译作业代码
func CompileWorkCodeFiles(config structs.JudgeConfiguration, configDir string) error {
    binRoot := path.Join(configDir, "bin")
    _, err := os.Stat(binRoot)
    if err != nil && os.IsNotExist(err) {
        err = os.MkdirAll(binRoot, 0775)
        if err != nil {
            return fmt.Errorf("cannot create binary work directory: %s", err.Error())
        }
    }
    //cpp := provider.GnucppCompileProvider{}
    return nil
}
