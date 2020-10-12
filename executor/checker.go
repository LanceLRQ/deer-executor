/* Deer executor
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package executor

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
)

// 使用IOReader读写每一行， 并移除空白字符
func readLineAndRemoveSpaceChar(buf *bufio.Reader) (string, error) {
	line, isContinue, err := buf.ReadLine()
	for isContinue && err == nil {
		var next []byte
		next, isContinue, err = buf.ReadLine()
		line = append(line, next...)
	}
	return removeSpaceChar(string(line)), err
}

// 通常情况下，我们定义Tab、换行和空格字符是"空白字符"
// Usually, tab, line break and white space are the special blank words, called 'SpaceChar'

// 移除空白字符
// Remove the special blank words in a string
func removeSpaceChar (source string) string {
	source = strings.Replace(source, "\t", "", -1)
	source = strings.Replace(source, "\r", "", -1)
	source = strings.Replace(source, "\n", "", -1)
	source = strings.Replace(source, " ", "", -1)
	return source
}
// 判断是否是空白字符
func isSpaceChar (ch byte) bool {
	return ch == '\n' ||  ch == '\r' || ch == ' ' || ch == '\t'
}

// 逐行比较，获取错误行数
// Compare each line, to find out the number of wrong line
func lineDiff (session *JudgeSession) (sameLines int, totalLines int) {
	answer, err := os.OpenFile(session.TestCaseOut, os.O_RDONLY | syscall.O_NONBLOCK, 0)
	if err != nil {
		return 0, 0
	}
	defer answer.Close()
	userout, err := os.Open(session.ProgramOut)
	if err != nil {
		return 0, 0
	}
	defer userout.Close()

	useroutBuffer := bufio.NewReader(userout)
	answerBuffer := bufio.NewReader(answer)

	var (
		leftStr, rightStr = "", ""
		leftErr, rightErr error = nil, nil
		leftCnt, rightCnt = 0, 0
	)

	for leftErr == nil {
		leftStr, leftErr = readLineAndRemoveSpaceChar(answerBuffer)
		if leftStr == "" {
			continue
		}

		leftCnt++

		for rightStr == "" && rightErr == nil {
			rightStr, rightErr = readLineAndRemoveSpaceChar(useroutBuffer)
		}

		if rightStr == leftStr {
			rightCnt++
		}
		rightStr = ""
	}

	return rightCnt, leftCnt
}

// 严格比较每一个字符
// Compare each char in buffer strictly
func StrictDiff(useroutBuffer, answerBuffer []byte, useroutLen, answerLen int64) bool {
	if useroutLen != answerLen {
		return false
	}
	pos := int64(0)
	for ; pos < useroutLen; pos++ {
		leftByte, rightByte := useroutBuffer[pos], answerBuffer[pos]
		if leftByte != rightByte {
			return false
		}
	}
	return true
}

// 比较每一个字符，但是忽略空白
// Compare each char in buffer, but ignore the 'SpaceChar'
func CharDiffIoUtil (useroutBuffer, answerBuffer []byte, useroutLen, answerLen int64) (rel int, logtext string) {
	var (
		leftPos, rightPos int64 = 0, 0
		maxLength = Max(useroutLen, answerLen)
		leftByte, rightByte byte
	)
	for (leftPos < maxLength) && (rightPos < maxLength) && (leftPos < useroutLen) && (rightPos < answerLen) {
		if leftPos < useroutLen {
			leftByte = useroutBuffer[leftPos]
		}
		if rightPos < answerLen {
			rightByte = answerBuffer[rightPos]
		}

		for leftPos < useroutLen && isSpaceChar(leftByte) {
			leftPos++
			if leftPos < useroutLen {
				leftByte = useroutBuffer[leftPos]
			} else {
				leftByte = 0
			}
		}
		for rightPos < answerLen && isSpaceChar(rightByte) {
			rightPos++
			if rightPos < answerLen {
				rightByte = answerBuffer[rightPos]
			} else {
				rightByte = 0
			}
		}

		if leftByte != rightByte {
			return JudgeFlagWA, fmt.Sprintf(
				"WA: at leftPos=%d, rightPos=%d, leftByte=%d, rightByte=%d",
				leftPos,
				rightPos,
				leftByte,
				rightByte,
			)
		}
		leftPos++
		rightPos++
	}

	// 如果左游标没跑完
	for leftPos < useroutLen {
		leftByte = useroutBuffer[leftPos]
		if !isSpaceChar(leftByte) {
			return JudgeFlagWA, fmt.Sprintf(
				"WA: leftPos=%d, rightPos=%d, leftLen=%d, rightLen=%d",
				leftPos,
				rightPos,
				useroutLen,
				answerLen,
			)
		}
		leftPos++
	}
	// 如果右游标没跑完
	for rightPos < answerLen {
		rightByte = answerBuffer[rightPos]
		if !isSpaceChar(rightByte) {
			return JudgeFlagWA, fmt.Sprintf(
				"WA: leftPos=%d, rightPos=%d, leftLen=%d, rightLen=%d",
				leftPos,
				rightPos,
				useroutLen,
				answerLen,
			)
		}
		rightPos++
	}
	// 左右匹配，说明AC
	// if left cursor's position equals right cursor's, means Accepted.
	if leftPos == rightPos {
		return JudgeFlagAC, "AC!"
	} else {
		return JudgeFlagPE, fmt.Sprintf(
			"PE: leftPos=%d, rightPos=%d, leftLen=%d, rightLen=%d",
			leftPos,
			rightPos,
			useroutLen,
			answerLen,
		)
	}
}

// 进行结果文本比较（主要工具）
// Compare the text
func DiffText(session JudgeSession, result *JudgeResult) (err error, logtext string) {
	answerInfo, err := os.Stat(session.TestCaseOut)
	if err != nil {
		result.JudgeResult = JudgeFlagSE
		return err, fmt.Sprintf("get answer file info failed: %s", err.Error())
	}
	useroutInfo, err := os.Stat(session.ProgramOut)
	if err != nil {
		result.JudgeResult = JudgeFlagSE
		return err, fmt.Sprintf("get userout file info failed: %s", err.Error())
	}

	useroutLen := useroutInfo.Size()
	answerLen := answerInfo.Size()

	sizeText := fmt.Sprintf("tcLen=%d, ansLen=%d", answerLen, useroutLen)

	var useroutBuffer, answerBuffer []byte
	errText := ""

	answerBuffer, errText, err = readFileWithTry(session.TestCaseOut, "answer", 3)
	if err != nil {
		result.JudgeResult = JudgeFlagSE
		return err, errText
	}

	useroutBuffer, errText, err = readFileWithTry(session.ProgramOut, "userout", 3)
	if err != nil {
		result.JudgeResult = JudgeFlagSE
		return err, errText
	}

	if useroutLen == 0 && answerLen == 0 {
		// Empty File AC
		result.JudgeResult = JudgeFlagAC
		return nil, sizeText + "; AC=zero size."
	} else if useroutLen > 0 && answerLen > 0 {
		if (useroutLen > int64(session.FileSizeLimit)) || (useroutLen > answerLen * 2) {
			// OLE
			result.JudgeResult = JudgeFlagOLE
			if useroutLen > int64(session.FileSizeLimit) {
				return nil, sizeText + "; WA: larger then limitation."
			} else {
				return nil, sizeText + "; WA: larger then 2 times."
			}
		}
	} else {
		// WTF?
		result.JudgeResult = JudgeFlagWA
		return nil, sizeText + "; WA: less then zero size"
	}

	rel, logText := CharDiffIoUtil(useroutBuffer, answerBuffer, useroutLen ,answerLen)
	result.JudgeResult = rel

	if rel != JudgeFlagWA {
		// PE or AC or SE
		if rel == JudgeFlagAC {
			sret := StrictDiff(useroutBuffer, answerBuffer, useroutLen ,answerLen)
			if !sret {
				result.JudgeResult = JudgeFlagPE
				logText = "strict check: PE"
			}
		}
		return nil, sizeText + "; " + logText
	} else {
		// WA
		sameLines, totalLines := lineDiff(&session)
		result.SameLines = sameLines
		result.TotalLines = totalLines
		return nil, sizeText + "; " + logText
	}
}
