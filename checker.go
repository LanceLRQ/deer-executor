/* Deer executor
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package deer_executor

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"io"
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
		leftStr, rightStr = "", ""
		leftErr, rightErr error = nil, nil
		leftCnt, rightCnt = 0, 0
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

func isSpaceChar (ch byte) bool {
	return ch == '\n' ||  ch == '\r' || ch == ' ' || ch == '\t'
}

func checkCRC(fp *os.File) (string, error) {
	if _, err := fp.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	reader := bufio.NewReader(fp)
	crc := crc32.NewIEEE()
	if  _, err := io.Copy(crc, reader); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", crc.Sum32()), nil
}

func CharDiff (userout *os.File, answer *os.File, useroutLen int64, answerLen int64) (rel int, logtext string) {
	_, _ = userout.Seek(0, io.SeekStart)
	_, _ = answer.Seek(0, io.SeekStart)

	useroutBuffer := bufio.NewReader(userout)
	answerBuffer := bufio.NewReader(answer)

	var (
		leftPos, rightPos int64 = 0, 0
		maxLength = Max(useroutLen, answerLen)
		leftErr, rightErr error = nil, nil
		leftByte, rightByte byte
	)

	// Lo-runner 源代码中对于格式错误的判断只是判断长度不同，没有判断字符不同的格式错误。
	// 这边很暴力的直接用CRC32去测试了
	for (leftPos < maxLength) && (rightPos < maxLength) {
		leftByte, leftErr = useroutBuffer.ReadByte()
		rightByte, rightErr = answerBuffer.ReadByte()

		if (leftErr != nil) && (rightErr != nil) {
			break
		}

		for leftErr == nil && isSpaceChar(leftByte) {
			leftByte, leftErr = useroutBuffer.ReadByte(); leftPos++
		}
		for  rightErr == nil && isSpaceChar(rightByte) {
			rightByte, rightErr = answerBuffer.ReadByte(); rightPos++
		}

		if leftByte != rightByte {
			return -1, fmt.Sprintf("WA: at leftPos=%d, rightPos=%d", leftByte, rightByte)
		}
		if leftErr == nil { leftPos++ }
		if rightErr == nil { rightPos++ }
	}

	if leftPos == useroutLen && rightPos == answerLen && leftPos == rightPos {
		crcuser, _ := checkCRC(userout)
		crcout, _ := checkCRC(answer)
		if crcuser != crcout {
			return JUDGE_FLAG_PE, "PE: CRC not match"
		}
		return JUDGE_FLAG_AC, "AC!"
	} else {
		return JUDGE_FLAG_PE, fmt.Sprintf(
			"PE: leftPos=%d, rightPos=%d, leftLen=%d, rightLen=%d",
			leftByte,
			rightByte,
			useroutLen,
			answerLen,
		)
	}
}

func DiffText(options JudgeOption, result *JudgeResult) (err error, logtext string) {
	answer, err := os.Open(options.TestCaseOut)
	if err != nil {
		result.JudgeResult = JUDGE_FLAG_SE
		return err, "open testcase_out error"
	}
	defer answer.Close()
	userout, err := os.Open(options.ProgramOut)
	if err != nil {
		result.JudgeResult = JUDGE_FLAG_SE
		return err, "open answer_out error"
	}
	defer userout.Close()

	useroutLen, _ := userout.Seek(0, os.SEEK_END)
	answerLen, _ := answer.Seek(0, os.SEEK_END)

	sizeText := fmt.Sprintf("tcLen=%d, ansLen=%d", answerLen, useroutLen)

	if useroutLen == 0 && answerLen == 0 {
		// Empty File AC
		result.JudgeResult = JUDGE_FLAG_AC
		return nil, sizeText + "; AC=zero size."
	} else if useroutLen > 0 && answerLen > 0 {
		if (useroutLen > int64(options.FileSizeLimit)) || (useroutLen > answerLen * 2) {
			// OLE
			result.JudgeResult = JUDGE_FLAG_OLE
			if useroutLen > int64(options.FileSizeLimit) {
				return nil, sizeText + "; WA: larger then limitation."
			} else {
				return nil, sizeText + "; WA: larger then 2 times."
			}
		}
	} else {
		// WTF?
		result.JudgeResult = JUDGE_FLAG_WA
		return nil, sizeText + "; WA: less then zero size"
	}

	rel, logText := CharDiff(userout, answer, useroutLen ,answerLen)

	if rel != -1 {
		// PE or AC
		result.JudgeResult = rel
		return nil, sizeText + "; " + logText
	} else {
		// WA
		sameLines, totalLines := lineDiff(userout, answer)
		result.JudgeResult = JUDGE_FLAG_WA
		result.SameLines = sameLines
		result.TotalLines = totalLines
		return nil, sizeText + "; " + logText
	}
}
