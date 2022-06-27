package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

var BIG5 = false

// 示例: zip.Unzip("./mcsm.zip", "./") 可使用相对路径和绝对路径
func Unzip(zipPath string, targetPath string) error {
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	if zipEncodeIsUtf8(zipReader.File) {
		fmt.Println("decode: utf8")
		err = decoderUtf8(zipReader.File, targetPath)
		if err != nil {
			fmt.Printf("decoderUtf8 err:%v", err)
			panic(err)
		}
	} else if zipEncodeIsGBK(zipReader.File) {
		if BIG5 {
			fmt.Println("decode: big5")
			err = decoderBIG5(zipReader.File, targetPath)
			if err != nil {
				fmt.Printf("decoderUtf8 err:%v", err)
				panic(err)
			}
		} else {
			fmt.Println("decode: gbk")
			err = decoderGBK(zipReader.File, targetPath)
			if err != nil {
				fmt.Printf("decoderUtf8 err:%v", err)
				panic(err)
			}
		}
	}
	return nil
}

func zipEncodeIsUtf8(f []*zip.File) bool {
	var i = 0
	var count = 0
	for _, v := range f {
		if i == 3 {
			break
		}
		if isUtf8([]byte(v.Name)) {
			count++
		}
		i++
	}
	fmt.Printf("count: %v\n", count)
	if count == i {
		return true
	} else {
		return false
	}
}

func zipEncodeIsGBK(f []*zip.File) bool {
	var i = 0
	var count = 0
	for _, v := range f {
		if i == 3 {
			break
		}
		if isGBK([]byte(v.Name)) {
			count++
		}
		i++
	}
	if count == i {
		return true
	} else {
		return false
	}
}

func decoderUtf8(files []*zip.File, targetPath string) error {
	var decodeName string
	var err error
	for _, f := range files {
		decodeName = f.Name
		fpath := filepath.Join(targetPath, decodeName)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
			inFile.Close()
			outFile.Close()
		}
	}
	return nil
}

func decoderGBK(files []*zip.File, targetPath string) error {
	var decodeName string
	var err error
	for _, f := range files {
		i := bytes.NewReader([]byte(f.Name))
		decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
		content, _ := ioutil.ReadAll(decoder)
		decodeName = string(content)
		fpath := filepath.Join(targetPath, decodeName)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
			inFile.Close()
			outFile.Close()
		}
	}
	return nil
}

func decoderBIG5(files []*zip.File, targetPath string) error {
	var decodeName string
	var err error
	for _, f := range files {
		i := bytes.NewReader([]byte(f.Name))
		decoder := transform.NewReader(i, traditionalchinese.Big5.NewDecoder())
		content, _ := ioutil.ReadAll(decoder)
		decodeName = string(content)
		fpath := filepath.Join(targetPath, decodeName)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
			inFile.Close()
			outFile.Close()
		}
	}
	return nil
}

// 先判断是否是UTF8再判断是否是其它编码才有意义
func isUtf8(data []byte) bool {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}

func isGBK(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		if data[i] <= 0x7f {
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0x7f {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func preNUm(data byte) int {
	var mask byte = 0x80
	var num int = 0
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}
