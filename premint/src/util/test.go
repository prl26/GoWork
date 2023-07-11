/*
 * @Descripttion: description
 * @Author: jianlinwei
 * @Date: 2023-03-11 17:32:17
 * @LastEditTime: 2023-03-11 17:34:02
 */
package util

import (
	"io"
	"os"
	"path/filepath"
)

func main() {
	src := "/path/to/source/folder"      // 源文件夹路径
	dst := "/path/to/destination/folder" // 目标文件夹路径

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		newPath := filepath.Join(dst, info.Name())
		newFile, err := os.Create(newPath)
		if err != nil {
			return err
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
