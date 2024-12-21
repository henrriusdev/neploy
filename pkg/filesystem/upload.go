package filesystem

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mholt/archives"
	"neploy.dev/config"
)

// UploadFile uploads a .zip file to the server
func UploadFile(file *multipart.FileHeader, appName string) (string, error) {
	if filepath.Ext(file.Filename) != ".zip" {
		return "", errors.New("invalid file format (only .zip files are allowed)")
	}

	// Sanitize the application name
	appName = sanitizeAppName(appName)

	// Define the file path
	filePath := filepath.Join(config.Env.UploadPath, appName+".zip")

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return "", err
	}

	// Create the destination file
	destFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	// Open the uploaded file
	srcFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	// Copy the uploaded file to the server
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return "", err
	}

	return filePath, nil
}

// UnzipFile extracts a .zip file to the specified destination directory
func UnzipFile(zipFilePath, appName string) (string, error) {
	// Sanitize the application name
	appName = sanitizeAppName(appName)

	// Define the destination directory
	destination := filepath.Join(config.Env.UploadPath, appName)

	// Ensure the destination directory exists
	if err := os.MkdirAll(destination, os.ModePerm); err != nil {
		return "", err
	}

	// Open the .zip file
	zipFile, err := os.Open(zipFilePath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	// Create a context
	ctx := context.Background()

	// Identify the archive format
	format, stream, err := archives.Identify(ctx, zipFilePath, zipFile)
	if err != nil {
		return "", err
	}

	// Ensure the format is a zip archive
	zipFormat, ok := format.(archives.Zip)
	if !ok {
		return "", errors.New("the file is not a valid .zip archive")
	}

	// Extract the archive
	err = zipFormat.Extract(ctx, stream, func(ctx context.Context, f archives.FileInfo) error {
		// Construct the full path for the extracted file

		// Handle directories
		if f.IsDir() {
			return os.MkdirAll(destination, os.ModePerm)
		}

		// Ensure the parent directory exists
		if err := os.MkdirAll(filepath.Dir(destination), os.ModePerm); err != nil {
			return err
		}

		// Create the extracted file
		outFile, err := os.Create(destination)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Open the file from the archive
		inFile, err := f.Open()
		if err != nil {
			return err
		}
		defer inFile.Close()

		// Copy the file contents
		_, err = io.Copy(outFile, inFile)
		return err
	})

	return destination, err
}

// sanitizeAppName cleans the application name for safe file system usage
func sanitizeAppName(appName string) string {
	appName = strings.ReplaceAll(appName, " ", "-")
	appName = regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(appName, "")
	return strings.ToLower(appName)
}
