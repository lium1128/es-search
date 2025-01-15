package utils

import (
	"errors"
	"os"
)

// PathExists 判断文件是否存在
func PathExists(path string) (bool, error) {
	if path == "" {
		return false, errors.New("路径为空,请检查")
	}
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
