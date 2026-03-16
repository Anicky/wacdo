package config

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUploader interface {
	Upload(ctx context.Context, file interface{}, uploadParams uploader.UploadParams) (*uploader.UploadResult, error)
}

var UploadAPI CloudinaryUploader

func ConnectCloudinary() {
	cloudinaryUrl := os.Getenv("CLOUDINARY_URL")

	cloudinaryInstance, err := cloudinary.NewFromURL(cloudinaryUrl)

	if err != nil {
		log.Fatal("Unable to connect to cloudinary: ", err)
	}

	UploadAPI = &cloudinaryInstance.Upload
}
