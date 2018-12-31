package deer

import (
	"bufio"
	"os"
	"strings"
)

func readLine(buf *bufio.Reader) (string, error) {
	line, isContinue, err := buf.ReadLine()
	for isContinue && err == nil {
		var next []byte
		next, isContinue, err = buf.ReadLine()
		line = append(line, next...)
	}
	return clearBlank(string(line)), err
}

func clearBlank (source string) string {
	source = strings.Replace(source, "\t", "", -1)
	source = strings.Replace(source, " ", "", -1)
	return source
}

func lineDiff (userout *os.File, answer *os.File) (sameLines int, totalLines int) {
	userout.Seek(0, os.SEEK_SET)
	answer.Seek(0, os.SEEK_SET)

	useroutBuffer := bufio.NewReader(userout)
	answerBuffer := bufio.NewReader(answer)

	var (
		leftStr, rightStr string = "", ""
		leftErr, rightErr error = nil, nil
		leftCnt, rightCnt int = 0, 0
	)

	for leftErr == nil {
		leftStr, leftErr = readLine(answerBuffer)
		if leftStr == "" {
			continue
		}

		leftCnt++

		for rightStr == "" && rightErr == nil {
			rightStr, rightErr = readLine(useroutBuffer)
		}

		if rightStr == leftStr {
			rightCnt++
		}
		rightStr = ""
	}

	return rightCnt, leftCnt
}

func charDiff (userout *os.File, answer *os.File, useroutLen int64, answerLen int64) (int) {
	userout.Seek(0, os.SEEK_SET)
	answer.Seek(0, os.SEEK_SET)

	useroutBuffer := bufio.NewReader(userout)
	answerBuffer := bufio.NewReader(answer)

	var (
		leftPos, rightPos int64 = 0, 0
		leftErr, rightErr error = nil, nil
		leftByte, rightByte byte
	)

	for (leftPos < useroutLen) && (rightPos < answerLen) && (leftErr == nil) && (rightErr == nil) {
		leftByte, leftErr = useroutBuffer.ReadByte()
		rightByte, rightErr = answerBuffer.ReadByte()
		for leftByte == '\n' || leftByte == '\r' || leftByte == ' ' || leftByte == '\t' && leftErr == nil {
			leftByte, leftErr = useroutBuffer.ReadByte(); leftPos++
		}
		for rightByte == '\n' ||  rightByte == '\r' || rightByte == ' ' || rightByte == '\t' && rightErr == nil {
			rightByte, rightErr = answerBuffer.ReadByte(); rightPos++
		}
		if leftByte != rightByte {
			return -1
		}
		if leftErr == nil { leftPos++ }
		if rightErr == nil { rightPos++ }
	}
	if leftPos == useroutLen && rightPos == answerLen && leftPos == rightPos {
		return JUDGE_FLAG_AC
	} else {
		return JUDGE_FLAG_PE
	}
}

func DiffText(options JudgeOption, result *JudgeResult) (err error) {
	answer, err := os.Open(options.TestCaseOut)
	defer answer.Close()
	if err != nil {
		result.JudgeResult = JUDGE_FLAG_SE
		return err
	}
	userout, err := os.Open(options.ProgramOut)
	defer userout.Close()
	if err != nil {
		result.JudgeResult = JUDGE_FLAG_SE
		return err
	}

	useroutLen, _ := userout.Seek(0, os.SEEK_END)
	answerLen, _ := answer.Seek(0, os.SEEK_END)

	if useroutLen == 0 && answerLen == 0 {
		result.JudgeResult = JUDGE_FLAG_AC
		return nil
	} else if useroutLen > 0 && answerLen > 0 {
		if (useroutLen > int64(options.FileSizeLimit)) || (useroutLen > answerLen * 2) {
			result.JudgeResult = JUDGE_FLAG_OLE
			return nil
		}
	} else {
		result.JudgeResult = JUDGE_FLAG_WA
		return nil
	}

	rel := charDiff(userout, answer, useroutLen ,answerLen)
	if rel != -1 {
		result.JudgeResult = rel
		return nil
	}

	sameLines, totalLines := lineDiff(userout, answer)

	result.JudgeResult = JUDGE_FLAG_WA
	result.SameLines = sameLines
	result.TotalLines = totalLines
	return nil
}
