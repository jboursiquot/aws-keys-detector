package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	hits, err := detectKeys()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, h := range hits {
		fmt.Printf("::set-output name=KeysFoundIn::%s|%s\n", h.file, h.line)
	}
}

type hit struct {
	file string
	line string
}

func detectKeys() ([]hit, error) {
	var hits []hit
	err := filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if hitsInFile, err := detectKeysInFile(path); err == nil {
			hits = append(hits, hitsInFile...)
		}
		return err
	})
	return hits, err
}

const secretAccessKeyPattern string = `[A-Z0-9]{20}`
const accessKeyIDPattern string = `(^[A-Za-z0-9\/\+\=]{40})`

func detectKeysInFile(path string) ([]hit, error) {
	var hits []hit

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Access Key ID?
		matched, err := regexp.Match(accessKeyIDPattern, []byte(line))
		if err != nil {
			log.Printf("error: %s line: %s", err, line)
			continue
		}
		if matched {
			hits = append(hits, hit{file: path, line: line})
		}

		// Secret Access Key?
		matched, err = regexp.Match(secretAccessKeyPattern, []byte(line))
		if err != nil {
			log.Printf("error: %s line: %s", err, line)
			continue
		}
		if matched {
			hits = append(hits, hit{file: path, line: line})
		}
	}

	if err := scanner.Err(); err != nil {
		return hits, err
	}

	return hits, nil
}
