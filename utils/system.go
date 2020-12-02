package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func ReadFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func WriteFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func RunExecutable(path string) (string, error) {
	var out bytes.Buffer
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to find working directory: %w", err)
	}
	full := fmt.Sprintf("%s/%s", wd, path)
	cmd := exec.Command(full)
	cmd.Stdout = &out
	err = cmd.Run()
	return out.String(), err
}
