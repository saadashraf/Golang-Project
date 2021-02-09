package main

import (
	"fmt"
	"os"
)

func settingArguments(input []string) (string, string, string) {
	copyPath := ""
	sourcePath := ""
	presentIndex := 1
	buildName := "none"
	if input[presentIndex] == "-builddir" {
		copyPath = input[presentIndex+1]
		presentIndex += 2
	}
	if input[presentIndex] == "-exe" {
		buildName = input[presentIndex+1]
		presentIndex += 2
	}
	if input[presentIndex] == "-copydir" {
		sourcePath = input[presentIndex+1]
	}

	return copyPath, sourcePath, buildName
}

func main() {
	input := os.Args[1:]
	fmt.Println(input)
	fmt.Println(len(input))

	path, _ := os.Getwd()
	copyPath, sourcePath, buildName := settingArguments(input)
	copyPath = path + copyPath
	sourcePath = path + sourcePath

	fmt.Println(copyPath)
	fmt.Println(sourcePath)
	fmt.Println(buildName)
	fmt.Println(path)
}
