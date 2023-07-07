package utils

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
)

var mimeExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/jpg":  ".jpg",
	"image/png":  ".png",
}

func GetFileExtension(file multipart.File) (string, error) {
	fileHeader := make([]byte, 512)
	_, err := file.Read(fileHeader)
	if err != nil {
		log.Printf("Error reading file header: %v", err)
		return "", err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Printf("Error seeking file: %v", err)
		return "", err
	}
	mime := http.DetectContentType(fileHeader)
	ext, ok := mimeExtensions[mime]
	if !ok {
		log.Printf("Unsupported file type: %s", mime)
		return "", fmt.Errorf("unsupported file type: %s", mime)
	}
	return ext, nil
}
