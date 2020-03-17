package utils

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

var (
	letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// RandID ... 生成指定长度的随机ID
func RandID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// GetAllFiles ... 根据文件路径获取是所有的 markdown 文件
func GetAllFiles(dir string) ([]string, error) {
	dirPath, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string

	sep := string(os.PathSeparator)

	for _, fi := range dirPath {
		if fi.IsDir() { // 如果还是一个目录，则递归去遍历
			subFiles, err := GetAllFiles(dir + sep + fi.Name())
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			// 过滤指定格式的文件
			ok := strings.HasSuffix(fi.Name(), ".md")
			if ok {
				files = append(files, dir + sep + fi.Name())
			}
		}
	}
	return files, nil
}
