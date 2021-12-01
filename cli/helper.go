package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func (c cli) writeFile(content []byte) {
	fName := fmt.Sprintf("%s.go", c.clientName)
	filePathName := filepath.Join(c.dst, fName)
	err := os.MkdirAll(c.dst, 0777)
	if err != nil {
		panic(err)
	}
	openFile, err := os.OpenFile(filePathName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	_, err = openFile.Write(content)
	if err != nil {
		panic(err)
	}
	if err := openFile.Close(); err != nil {
		panic(err)
	}
}
