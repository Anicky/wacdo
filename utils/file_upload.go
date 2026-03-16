package utils

import (
	"net/http"
	"wacdo/config"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func UploadBase64Image(context *gin.Context, base64ImageData string) (*string, error) {
	response, err := config.UploadAPI.Upload(context, base64ImageData, uploader.UploadParams{})

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to upload image."})

		return nil, err
	}

	return &response.SecureURL, nil
}
