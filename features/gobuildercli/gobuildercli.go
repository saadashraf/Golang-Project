package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PowerShell struct {
	powerShell string
}

func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

func (p *PowerShell) Execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}

func build(buildName string) {
	posh := New()
	stdout, stderr, err := posh.Execute("go build -o " + buildName + " test.go")

	fmt.Println(stdout)
	fmt.Println(stderr)

	if err != nil {
		fmt.Println(err)
	}
}

func buildInPath(buildName string, path string) {
	posh := New()
	fmt.Println(path)
	stdout, stderr, err := posh.Execute("cd " + path + "\ngo build -o " + buildName + " test.go")

	fmt.Println(stdout)
	fmt.Println(stderr)

	if err != nil {
		fmt.Println(err)
	}
}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func CopyDir(src string, dst string, excludeTests bool) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if strings.Contains(entry.Name(), "test") && excludeTests == true {
			continue
		}
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath, excludeTests)
			if err != nil {
				return
			}
		} else {
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func settingArguments(input []string) (string, string, string, bool) {
	copyPath := ""
	sourcePath := ""
	presentIndex := 1
	buildName := "none"
	excludeTests := false
	if input[presentIndex] == "-builddir" {
		sourcePath = input[presentIndex+1]
		presentIndex += 2
	}
	if input[presentIndex] == "-exe" {
		buildName = input[presentIndex+1]
		presentIndex += 2
	}
	if input[presentIndex] == "-copydir" {
		copyPath = input[presentIndex+1]
		presentIndex += 2
	}

	/*	if input[presentIndex] == "-exclude-tests" {
		excludeTests = true
		presentIndex += 2
	}*/

	return copyPath, sourcePath, buildName, excludeTests
}

func main() {
	input := os.Args[1:]
	fmt.Println(input)
	fmt.Println(len(input))

	path, _ := os.Getwd()
	dstPath := "F:\\job\\evatix\\phase 2\\Golang-Project\\features"
	copyPath, sourcePath, buildName, excludeTests := settingArguments(input)

	if buildName != "none" {
		if sourcePath == "" {
			build(buildName)
		} else {
			buildInPath(buildName, sourcePath)
		}
	}

	destination := dstPath + "\\" + copyPath
	source := path + "\\" + sourcePath

	fmt.Println(excludeTests)

	if sourcePath != copyPath {
		err := CopyDir(source, destination, excludeTests)

		if err != nil {
			fmt.Println(err)
		}
	}
}
