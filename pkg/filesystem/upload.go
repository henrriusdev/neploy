package filesystem

import (
	"archive/zip"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"neploy.dev/config"
	"neploy.dev/pkg/logger"
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

// UnzipFile extracts a .zip file to a temporary directory
func UnzipFile(zipFilePath, appName string) (string, error) {
	// Sanitize the application name
	appName = sanitizeAppName(appName)

	// Create a temporary directory outside the app structure
	tempDir := filepath.Join(os.TempDir(), "neploy_extract_"+appName)
	destination := tempDir

	// Ensure the destination directory exists
	if err := os.MkdirAll(destination, os.ModePerm); err != nil {
		logger.Error("error creating destination directory: %v", err)
		return "", err
	}

	// Open the zip file using Go's standard library
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		logger.Error("error opening zip file: %v", err)
		return "", err
	}
	defer reader.Close()

	// Find the common root directory in the zip
	rootDir := ""
	for _, file := range reader.File {
		if rootDir == "" {
			// Get the first directory component
			parts := strings.Split(file.Name, "/")
			if len(parts) > 1 {
				rootDir = parts[0] + "/"
			}
		}
	}

	// Extract files
	for _, file := range reader.File {
		// Strip the root directory if it exists
		relativePath := file.Name
		if rootDir != "" && strings.HasPrefix(file.Name, rootDir) {
			relativePath = strings.TrimPrefix(file.Name, rootDir)
		}
		
		// Skip if this is just the root directory itself
		if relativePath == "" {
			continue
		}

		// Construct the full path for the extracted file
		fullPath := filepath.Join(destination, relativePath)

		// Handle directories
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fullPath, file.FileInfo().Mode()); err != nil {
				logger.Error("error creating directory: %v", err)
				return "", err
			}
			continue
		}

		// Ensure the parent directory exists
		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			logger.Error("error creating parent directory: %v", err)
			return "", err
		}

		// Open the file from the archive
		rc, err := file.Open()
		if err != nil {
			logger.Error("error opening file from archive: %v", err)
			return "", err
		}

		// Create the extracted file
		outFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			rc.Close()
			logger.Error("error creating extracted file: %v", err)
			return "", err
		}

		// Copy the file contents
		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			logger.Error("error copying file contents: %v", err)
			return "", err
		}
	}

	return destination, nil
}

// sanitizeAppName cleans the application name for safe file system usage
func sanitizeAppName(appName string) string {
	appName = strings.ReplaceAll(appName, " ", "-")
	appName = regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(appName, "")
	return strings.ToLower(appName)
}

// CopyDir recursively copies a directory tree from src to dst
func CopyDir(src, dst string) error {
	// Get properties of source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		logger.Error("error getting source directory info: %v", err)
		return err
	}

	// Ensure source is a directory
	if !srcInfo.IsDir() {
		return errors.New("source is not a directory")
	}

	// Create the destination directory with proper permissions
	if err = os.MkdirAll(dst, 0755); err != nil {
		logger.Error("error creating destination directory: %v", err)
		return err
	}

	// Read the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		logger.Error("error reading source directory: %v", err)
		return err
	}

	// Process each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// Skip hidden files and system files that might cause issues
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// If the entry is a directory, recursively copy it
		if entry.IsDir() {
			if err = CopyDir(srcPath, dstPath); err != nil {
				logger.Error("error copying directory %s: %v", srcPath, err)
				return err
			}
		} else {
			// Otherwise, copy the file
			if err = CopyFile(srcPath, dstPath); err != nil {
				logger.Error("error copying file %s: %v", srcPath, err)
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a single file from src to dst
func CopyFile(src, dst string) error {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		logger.Error("error opening source file: %v", err)
		return err
	}
	defer srcFile.Close()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		logger.Error("error getting source file info: %v", err)
		return err
	}

	// Create the destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		logger.Error("error creating destination file: %v", err)
		return err
	}
	defer dstFile.Close()

	// Copy the contents
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		logger.Error("error copying file contents: %v", err)
		return err
	}

	return nil
}
