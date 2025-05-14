package main

import (
	"fmt"
	"os"

	"github.com/ddddami/bindle/fs"
)

func main() {
	tempPath := "./temp_test/nested/path"

	err := fs.CreateDirIfNotExistsWithPerm(tempPath, 0o755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
	} else {
		fmt.Printf("Created directory: %s\n", tempPath)
	}

	// Clean up at the end
	defer fs.SafeRemove("./temp_test")

	// Create a test file
	testFilePath := tempPath + "/test.txt"
	err = os.WriteFile(testFilePath, []byte("Hello, this is a test file!"), 0o644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	} else {
		fmt.Printf("Created test file: %s\n", testFilePath)
	}

	// Check if file exists
	if fs.PathExists(testFilePath) {
		fmt.Printf("Confirmed file exists: %s\n", testFilePath)
	}

	files, err := fs.ListFilesWithExt(tempPath, ".txt")
	if err != nil {
		fmt.Printf("Error listing files: %v\n", err)
	} else {
		fmt.Printf("Found %d .txt files in %s\n", len(files), tempPath)
		for i, file := range files {
			fmt.Printf("  %d: %s\n", i+1, file)
		}
	}

	destFileName := tempPath + "/copy_test.txt"
	err = fs.CopyFile(testFilePath, destFileName)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
	} else {
		fmt.Printf("Copied file to: %s\n", destFileName)
		if fs.PathExists(destFileName) {
			fmt.Printf("Confirmed copied file exists: %s\n", destFileName)
		}
	}
}
