package utils

import (
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func UploadBase64Image(context *gin.Context, base64ImageData string) (*string, error) {
	cloudinaryInstance, err := getCloudinaryInstance(context)

	if cloudinaryInstance != nil {
		response, err := cloudinaryInstance.Upload.Upload(context, base64ImageData, uploader.UploadParams{})

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to upload image."})

			return nil, err
		}

		return &response.SecureURL, nil
	}

	return nil, err
}

func getCloudinaryInstance(context *gin.Context) (*cloudinary.Cloudinary, error) {
	cloudinaryInstance, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to contact image server."})

		return nil, err
	}

	return cloudinaryInstance, err
}
