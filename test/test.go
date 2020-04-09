package main

import (
	"fmt"
	deer_executor "github.com/LanceLRQ/deer-executor"
	"io"
	"os"
)

func main() {
	answer, err := os.Open("./test_cases/answer.out")
	if err != nil {
		panic(err)
	}
	defer answer.Close()
	userout, err := os.Open("./test_cases/test.out")
	if err != nil {
		panic(err)
	}
	defer userout.Close()

	useroutLen, _ := userout.Seek(0, io.SeekEnd)
	answerLen, _ := answer.Seek(0, io.SeekEnd)

	fmt.Println(deer_executor.CharDiffIoUtil(userout, answer, useroutLen, answerLen))
}
