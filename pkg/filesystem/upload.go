package filesystem

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"neploy.dev/config"
)

func UploadFile(file *multipart.FileHeader, appName string) (string, error) {
	if filepath.Ext(file.Filename) != ".zip" {
		return "", errors.New("invalid file format (only .zip files are allowed}")
	}

	appNameWithoutSpace := strings.ReplaceAll(appName, " ", "-")
	appNameWithoutSpecialChars := regexp.MustCompile(`[^a-zA-Z0-9-]`).ReplaceAllString(appNameWithoutSpace, "")
	appName = strings.ToLower(appNameWithoutSpecialChars)

	filePath := filepath.Join(config.Env.UploadPath, file.Filename)
	// Create the directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return "", err
	}

	// Create the file
	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Copy the uploaded file to the created file on the server
	_, err = file.Open()
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	_, err = io.Copy(f, src)
	if err != nil {
		return "", err
	}

	return file.Filename, nil
}
