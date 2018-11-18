package main

import (
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

func processFiles(ctx *cli.Context, path string, callback func(string) string) error {
	if IsDir(path) {
		names := ScanDir(path)
		for _, name := range names {
			processFiles(ctx, path+string(os.PathSeparator)+name, callback)
		}

		return nil
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	content := callback(string(b))
	return ioutil.WriteFile(path, []byte(content), 0666)
}

// ScanDir 列出指定路径中的文件和目录
// 如果目录不存在，则返回空slice
func ScanDir(directory string) []string {
	file, err := os.Open(directory)
	if err != nil {
		return []string{}
	}
	names, err := file.Readdirnames(-1)
	if err != nil {
		return []string{}
	}
	return names
}

// IsDir 判断给定文件名是否是一个目录
// 如果文件名存在并且为目录则返回 true。如果 filename 是一个相对路径，则按照当前工作目录检查其相对路径。
func IsDir(filename string) bool {
	return isFileOrDir(filename, true)
}

// IsFile 判断给定文件名是否为一个正常的文件
// 如果文件存在且为正常的文件则返回 true
func IsFile(filename string) bool {
	return isFileOrDir(filename, false)
}

// isFileOrDir 判断是文件还是目录，根据decideDir为true表示判断是否为目录；否则判断是否为文件
func isFileOrDir(filename string, decideDir bool) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	isDir := fileInfo.IsDir()
	if decideDir {
		return isDir
	}
	return !isDir
}
