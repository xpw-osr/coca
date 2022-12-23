package javaapp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func writeJsonFile(filename string, appType string, json string) (string, bool) {
	curDir, _ := os.Getwd()
	reporterPath := curDir + "/coca_reporter/" + appType
	if _, err := os.Stat(reporterPath); os.IsNotExist(err) {
		mkdirErr := os.MkdirAll(reporterPath, os.ModePerm)
		if mkdirErr != nil {
			fmt.Println(mkdirErr)
		}
	}

	file := filepath.FromSlash(reporterPath + "/" + filename)
	if err := ioutil.WriteFile(file, []byte(json), os.ModePerm); err == nil {
		return file, true
	}
	return "", false
}

func ReadJsonFile(path string) ([]byte, bool) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("* Failed read file: %s \n", err)
		return nil, false
	}
	return contents, true
}
