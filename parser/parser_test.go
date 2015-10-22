package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	languages = []string{"java", "asp", "csharp"}
)

// getTestFiles returns a list of files from the test folder
func getTestFiles() ([]string, error) {
	fil, err := os.Open("test")
	if err != nil {
		return nil, fmt.Errorf("Unable to open test folder: %s", err)
	}
	defer fil.Close()

	var info []os.FileInfo
	result := make([]string, 0)

	info, err = fil.Readdir(256)

	for err == nil {
		for i := range info {
			name := info[i].Name()
			if strings.HasSuffix(name, ".babel") {
				result = append(result, filepath.Join("test", name))
			}
		}
		info, err = fil.Readdir(256)
	}

	return result, nil
}

func TestBabelFiles(t *testing.T) {
	files, err := getTestFiles()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, lang := range languages {
		for _, file := range files {
			_, err := ParseIdl(file, lang)
			if strings.HasSuffix(file, "_bad.babel") {
				if err == nil {
					t.Errorf("The parser allowed file \"%s\" but it is supposed to fail.", file)
				}
			} else {
				if err != nil {
					t.Errorf("The parser failed file \"%s\" which should have succeeded: %s", file, err)
				}
			}
		}
	}
}
