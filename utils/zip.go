package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 示例 zip.Zip("MCSManager 9.4.5_win64_x86", "./test.zip") 可使用相对路径和绝对路径
func Zip(filePath, zipPath string) error {
	os.RemoveAll(zipPath)
	var err error
	// 转化相对路径
	if isAbs := filepath.IsAbs(filePath); !isAbs {
		filePath, err = filepath.Abs(filePath)
		if err != nil {
			return err
		}
	}
	zipfile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := zipfile.Close(); err != nil {
			fmt.Printf("*File close error: %v, file: %s", err, zipfile.Name())
		}
	}()
	//创建zip.Writer
	zw := zip.NewWriter(zipfile)
	defer func() {
		if err := zw.Close(); err != nil {
			fmt.Printf("zipwriter close error: %v", err)
		}
	}()
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(filePath)
	}

	err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filePath == path {
			return nil
		}
		var zipfile io.Writer
		if info.IsDir() {
			zipfile, _ = zw.Create(baseDir + strings.TrimPrefix(path, filePath) + `/`)
			if err != nil {
				panic(err)
			}
		} else {
			zipfile, err = zw.Create(baseDir + strings.TrimPrefix(path, filePath))
			if err != nil {
				panic(err)
			}
		}
		f1, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		io.Copy(zipfile, f1)
		f1.Close()
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk err:%v", err)
		return err
	}
	return nil
}
