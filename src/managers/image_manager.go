package managers

import (
	"bytes"
	"fmt"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/expense_errors"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/models"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type ImageMgr interface {
	UploadImage(file *multipart.FileHeader, userId *uuid.UUID) (string, *models.ExpenseServiceError)
	UploadDefaultProfilePicture(userId *uuid.UUID) (string, *models.ExpenseServiceError)
}

type ImageManager struct {
	Client *http.Client
}

func (im *ImageManager) UploadImage(header *multipart.FileHeader, userId *uuid.UUID) (string, *models.ExpenseServiceError) {
	var ENVIRONMENT = os.Getenv("ENVIRONMENT")
	var IMAGE_SERVICE_API = os.Getenv(fmt.Sprintf("%s_IMAGE_SERVICE_API", ENVIRONMENT))

	// Create multipart writer
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Create form field
	if err := writer.WriteField("userId", userId.String()); err != nil {
		return "", expense_errors.EXPENSE_BAD_REQUEST
	}

	// Create form file
	file, err := writer.CreateFormFile("image", header.Filename)
	if err != nil {
		log.Printf("Error creating form file: %v", err)
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Open image
	image, err := header.Open()
	if err != nil {
		log.Printf("Error opening image: %v", err)
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	fileExtension, err := utils.GetFileExtension(image)
	if err != nil {
		log.Printf("Error getting file extension: %v", err)
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Defer closing image and return expense error if closing fails
	defer func(image multipart.File) *models.ExpenseServiceError {
		err := image.Close()
		if err != nil {
			log.Printf("Error closing image: %v", err)
			return expense_errors.EXPENSE_INTERNAL_ERROR
		}

		return nil
	}(image)

	// Copy image to file
	if _, err = io.Copy(file, image); err != nil {
		log.Printf("Error copying image to file: %v", err)
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Close writer
	if err := writer.Close(); err != nil {
		log.Printf("Error closing writer: %v", err)
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	log.Printf("Uploading image to %s/images", IMAGE_SERVICE_API)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/images", IMAGE_SERVICE_API), &b)
	if err != nil {
		log.Printf("Error while uploading image: %s", err.Error())
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := im.Client.Do(req)
	if err != nil {
		log.Printf("Error while uploading image: %s", err.Error())
		return "", expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	if resp.StatusCode != http.StatusCreated {
		log.Print("Error while uploading image: status code not 201")
		return "", expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	return fileExtension, nil
}

func (im *ImageManager) UploadDefaultProfilePicture(userId *uuid.UUID) (string, *models.ExpenseServiceError) {
	var ENVIRONMENT = os.Getenv("ENVIRONMENT")
	var IMAGE_SERVICE_API = os.Getenv(fmt.Sprintf("%s_IMAGE_SERVICE_API", ENVIRONMENT))

	// Create multipart writer
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Create form field
	if err := writer.WriteField("userId", userId.String()); err != nil {
		log.Printf("Error creating form field: %v", err)
		return "", expense_errors.EXPENSE_BAD_REQUEST
	}

	// Get default profile picture from file system
	image, err := os.Open("static/default_avatar.png")
	if err != nil {
		log.Printf("Error opening default profile picture: %v", err)
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Defer closing image and return expense error if closing fails
	defer func(image multipart.File) *models.ExpenseServiceError {
		err := image.Close()
		if err != nil {
			log.Printf("Error closing image: %v", err)
			return expense_errors.EXPENSE_INTERNAL_ERROR
		}

		return nil
	}(image)

	fileExtension, err := utils.GetFileExtension(image)
	if err != nil {
		log.Printf("Error getting file extension: %v", err)
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Create form file
	file, err := writer.CreateFormFile("image", "default_profile_picture.png")
	if err != nil {
		log.Printf("Error creating form file: %v", err)
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Copy image to file
	if _, err = io.Copy(file, image); err != nil {
		log.Printf("Error copying image to file: %v", err)
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	// Close writer
	if err := writer.Close(); err != nil {
		log.Printf("Error closing writer: %v", err)
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/images", IMAGE_SERVICE_API), &b)
	if err != nil {
		log.Printf("Error while uploading image: %s", err.Error())
		return "", expense_errors.EXPENSE_INTERNAL_ERROR
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := im.Client.Do(req)
	if err != nil {
		log.Printf("Error while uploading image: %s", err.Error())
		return "", expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	if resp.StatusCode != http.StatusCreated {
		log.Print("Error while uploading image: status code not 201")
		return "", expense_errors.EXPENSE_UPSTREAM_ERROR
	}

	return fileExtension, nil
}
