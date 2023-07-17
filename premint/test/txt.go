package main

import (
	"fmt"
	"os"
)

func main() {
	fileName := "test.txt"

	dstFile, _ := os.Create(fileName)
	fmt.Println("hello")
	defer dstFile.Close()
	s := "jel\n"
	dstFile.WriteString(s)
	s1 := "hello World\n"
	dstFile.WriteString(s1)

}
