package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// processPaths iterates through the given paths and walks them to process files.
func processPaths(paths []string) {
	for _, path := range paths {
		err := filepath.Walk(path, processFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking path %s: %v\n", path, err)
		}
	}
}

// processFile is the core function called by filepath.Walk for each file and directory.
func processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		ext := filepath.Ext(path)
		switch ext {
		case extHTML, extCSS, extJS, extTS, extVue, extGo:
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", path, err)
				return nil
			}

			var newContent []byte
			switch ext {
			case extHTML:
				newContent, err = removeHTMLComments(content)
			case extCSS:
				newContent = removeCSSComments(content)
			case extJS, extTS, extGo:
				newContent = removeJSComments(content)
			case extVue:
				newContent = removeVueComments(content)
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing file %s: %v\n", path, err)
				return nil
			}

			if !bytes.Equal(content, newContent) {
				err = os.WriteFile(path, newContent, info.Mode())
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing file %s: %v\n", path, err)
				} else {
					fmt.Printf("Removed comments from %s\n", path)
				}
			}
		}
	}
	return nil
}
