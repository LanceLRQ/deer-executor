package problems

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileNotFoundError struct {
	FileName string
}

func (e FileNotFoundError) Error() string {
	return fmt.Sprintf("file (%s) not found", e.FileName)
}

func IsFileNotFoundError(e error) bool {
	_, ok := e.(FileNotFoundError)
	return ok
}

func NewFileNotFoundError(fileName string) error {
	return FileNotFoundError{FileName: fileName}
}

func UnZip(zipArchive *zip.ReadCloser, destDir string) error {
	return WalkZip(zipArchive, func(f *zip.File) error {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func WalkZip(zipArchive *zip.ReadCloser, walkFunc func(file *zip.File) error) error {
	for _, f := range zipArchive.File {
		err := walkFunc(f)
		if err != nil {
			return err
		}
	}
	return nil
}

// 搜索zip内的文件并打开(精确匹配)
func FindInZip(zipArchive *zip.ReadCloser, fileName string) (*io.ReadCloser, *zip.File, error) {
	var fileResult io.ReadCloser
	finded := false
	var fileInfo *zip.File
	err := WalkZip(zipArchive, func(file *zip.File) error {
		fileInfo = file
		if fileName == file.Name {
			finded = true
			var err error
			fileResult, err = file.Open()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if !finded {
		return nil, nil, NewFileNotFoundError(fileName)
	}
	return &fileResult, fileInfo, err
}
