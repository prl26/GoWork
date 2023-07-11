/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-11 17:43:43
 * @LastEditTime: 2023-03-11 17:45:47
 */
package util

import (
	"io"
	"os"
	"path/filepath"
)

func CopyDir(src string, dst string) error {
	// 获取源文件夹信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标文件夹
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	// 获取源文件夹下的所有文件和文件夹信息
	files, err := filepath.Glob(filepath.Join(src, "*"))
	if err != nil {
		return err
	}

	// 遍历所有文件和文件夹信息
	for _, file := range files {
		// 获取文件信息
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}

		// 如果是文件夹则递归调用该函数复制该文件夹
		if fileInfo.IsDir() {
			err = CopyDir(file, filepath.Join(dst, fileInfo.Name()))
			if err != nil {
				return err
			}
		} else {
			// 如果是文件则复制该文件到目标文件夹下
			err = copyFile(file, filepath.Join(dst, fileInfo.Name()), fileInfo.Mode())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src string, dst string, mode os.FileMode) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// 设置目标文件的权限
	err = os.Chmod(dst, mode)
	if err != nil {
		return err
	}

	return nil
}
