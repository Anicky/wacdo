package tests

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryMock struct {
}

func (m *CloudinaryMock) Upload(ctx context.Context, file interface{}, uploadParams uploader.UploadParams) (*uploader.UploadResult, error) {
	return &uploader.UploadResult{
		SecureURL: "https://res.cloudinary.com/demo/image/upload/sample.jpg",
	}, nil
}
