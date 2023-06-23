package utils

import (
	"bufio"
	"os"
	"strings"
)

func FileLinesToSlice(filename string) ([]string, error) {
	lines := []string{}

	file, err := os.OpenFile(filename, os.O_RDONLY, 0640)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}
