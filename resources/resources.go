package resources

import (
	"fmt"
	"os"
	"path/filepath"
)

func Get(filename string) *os.File {
	relativePath, err := GetPath(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	// Open the file
	file, err := os.Open(relativePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	return file
}

func GetPath(filename string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return "", err
	}

	return filepath.Join(currentDir, fmt.Sprintf("/resources/%s", filename)), nil
}
