package entity

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src string, dest string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the content from the source file to the destination file
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Ensure that the content is flushed to the destination file
	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// CopyDir copies the contents of the source directory to the destination directory
func CopyDir(src string, dest string) error {
	// Walk through the source directory
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Determine the path of the target item
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dest, relPath)

		// Check if it's a directory
		if info.IsDir() {
			// Create the directory in the destination
			return os.MkdirAll(destPath, info.Mode())
		}

		// If it's a file, copy the file
		return copyFile(path, destPath)
	})
}

// copyFile copies a file from src to dest
func copyFile(src, dest string) error {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// Copy file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dest, srcInfo.Mode())
}
