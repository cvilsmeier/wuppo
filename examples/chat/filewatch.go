package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func watchFiles(millis time.Duration) {
	// stop if a *.go file changes
	lastCheck := ""
	for {
		check := ""
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".go") {
				check = check + info.ModTime().String() + "_"
			}
			return nil
		})
		if lastCheck != "" && lastCheck != check {
			fmt.Printf("some watched file changed, game over\n")
			os.Exit(0)
		}
		lastCheck = check
		time.Sleep(millis * time.Millisecond)
	}
}
