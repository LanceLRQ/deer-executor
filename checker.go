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
	"io/ioutil"
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

func lineDiff (options *JudgeOption) (sameLines int, totalLines int) {
	answer, err := os.Open(options.TestCaseOut)
	if err != nil {
		return 0, 0
	}
	defer answer.Close()
	userout, err := os.Open(options.ProgramOut)
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
	reader := bufio.NewReader(fp)
	crc := crc32.NewIEEE()
	if  _, err := io.Copy(crc, reader); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", crc.Sum32()), nil
}

func compareCRC(options *JudgeOption) (bool, string) {
	answer, err := os.Open(options.TestCaseOut)
	if err != nil {
		return false, fmt.Sprintf("open answer file error: %s", err.Error())
	}
	defer answer.Close()
	userout, err := os.Open(options.ProgramOut)
	if err != nil {
		return false, fmt.Sprintf("open userout file error: %s", err.Error())
	}
	defer userout.Close()

	crcuser, _ := checkCRC(userout)
	crcout, _ := checkCRC(answer)
	if crcuser != crcout {
		return false, "PE: CRC not match"
	}
	return true, "AC!"
}


// 使用ioutil.ReadAll来读
func CharDiffIoUtil (options *JudgeOption, useroutLen int64, answerLen int64) (rel int, logtext string) {
	answer, err := os.Open(options.TestCaseOut)
	if err != nil {
		return JUDGE_FLAG_SE, fmt.Sprintf("open answer file error: %s", err.Error())
	}
	defer answer.Close()
	userout, err := os.Open(options.ProgramOut)
	if err != nil {
		return JUDGE_FLAG_SE, fmt.Sprintf("open userout file error: %s", err.Error())
	}
	defer userout.Close()

	useroutBuffer, useroutError := ioutil.ReadAll(userout)
	answerBuffer, answerError := ioutil.ReadAll(answer)

	if useroutError != nil {
		return JUDGE_FLAG_SE, fmt.Sprintf("read userout file io error: %s", useroutError.Error())
	}

	if answerError != nil {
		return JUDGE_FLAG_SE, fmt.Sprintf("read answer file io error: %s", answerError.Error())
	}

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
			return JUDGE_FLAG_WA, fmt.Sprintf(
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

	if leftPos - 1 == useroutLen && rightPos - 1 == answerLen && leftPos == rightPos {
		return JUDGE_FLAG_AC, "AC!"
	} else {
		return JUDGE_FLAG_PE, fmt.Sprintf(
			"PE: leftPos=%d, rightPos=%d, leftLen=%d, rightLen=%d",
			leftPos,
			rightPos,
			useroutLen,
			answerLen,
		)
	}
}

func DiffText(options JudgeOption, result *JudgeResult) (err error, logtext string) {
	answer, err := os.Stat(options.TestCaseOut)
	if err != nil {
		result.JudgeResult = JUDGE_FLAG_SE
		return err, fmt.Sprintf("get answer file info failed: %s", err.Error())
	}
	userout, err := os.Stat(options.ProgramOut)
	if err != nil {
		result.JudgeResult = JUDGE_FLAG_SE
		return err, fmt.Sprintf("get userout file info failed: %s", err.Error())
	}

	useroutLen := userout.Size()
	answerLen := answer.Size()

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

	rel, logText := CharDiffIoUtil(&options, useroutLen ,answerLen)
	result.JudgeResult = rel

	if rel != JUDGE_FLAG_WA {
		// PE or AC or SE
		if rel == JUDGE_FLAG_AC {
			ret := false
			ret, logText = compareCRC(&options)
			if !ret {
				result.JudgeResult = JUDGE_FLAG_PE
			}
		}
		return nil, sizeText + "; " + logText
	} else {
		// WA
		sameLines, totalLines := lineDiff(&options)
		result.SameLines = sameLines
		result.TotalLines = totalLines
		return nil, sizeText + "; " + logText
	}
}

//func CharDiff (userout *os.File, answer *os.File, useroutLen int64, answerLen int64) (rel int, logtext string) {
//	_, _ = userout.Seek(0, io.SeekStart)
//	_, _ = answer.Seek(0, io.SeekStart)
//
//	useroutBuffer := bufio.NewReader(userout)
//	answerBuffer := bufio.NewReader(answer)
//
//	var (
//		leftPos, rightPos int64 = 0, 0
//		maxLength = Max(useroutLen, answerLen)
//		leftErr, rightErr error = nil, nil
//		leftByte, rightByte byte
//	)
//
//	// Lo-runner 源代码中对于格式错误的判断只是判断长度不同，没有判断字符不同的格式错误。
//	// 这边很暴力的直接用CRC32去测试了
//	for (leftPos < maxLength) && (rightPos < maxLength) {
//		leftByte, leftErr = useroutBuffer.ReadByte()
//		rightByte, rightErr = answerBuffer.ReadByte()
//
//		if (leftErr != nil) && (rightErr != nil) {
//			break
//		}
//
//		for leftErr == nil && isSpaceChar(leftByte) {
//			leftByte, leftErr = useroutBuffer.ReadByte(); leftPos++
//		}
//		for  rightErr == nil && isSpaceChar(rightByte) {
//			rightByte, rightErr = answerBuffer.ReadByte(); rightPos++
//		}
//
//		if leftByte != rightByte {
//			return JUDGE_FLAG_WA, fmt.Sprintf(
//				"WA: at leftPos=%d, rightPos=%d, leftByte=%d, rightByte=%d",
//				leftPos,
//				rightPos,
//				leftByte,
//				rightByte,
//			)
//		}
//		if leftErr == nil { leftPos++ }
//		if rightErr == nil { rightPos++ }
//	}
//
//	if leftPos == useroutLen && rightPos == answerLen && leftPos == rightPos {
//		crcuser, _ := checkCRC(userout)
//		crcout, _ := checkCRC(answer)
//		if crcuser != crcout {
//			return JUDGE_FLAG_PE, "PE: CRC not match"
//		}
//		return JUDGE_FLAG_AC, "AC!"
//	} else {
//		return JUDGE_FLAG_PE, fmt.Sprintf(
//			"PE: leftPos=%d, rightPos=%d, leftLen=%d, rightLen=%d",
//			leftPos,
//			rightPos,
//			useroutLen,
//			answerLen,
//		)
//	}
//}