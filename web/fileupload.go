package web

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type FileUpload struct {
	Data []byte
	Ext  string
	Mime string
}

func (u FileUpload) Save(s3client *s3.Client, parentDirectory string) (string, error) {

	if u.Data != nil {
		reader := bytes.NewReader(u.Data)
		newName := uuid.New().String() + u.Ext
		_, err := s3client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(GetEnv("AWS_BUCKET_NAME", "")),
			Key:         aws.String(parentDirectory + "/" + newName),
			Body:        reader,
			ContentType: aws.String(u.Mime),
		})

		if err != nil {
			return "", err
		}
		return parentDirectory + "/" + newName, nil
	}

	return "", nil
}
