package db

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func GetContentFile(id int) (*s3.GetObjectOutput, error) {
	o, err := s3client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("flatgrass-toybox-content"),
		Key:    aws.String(strconv.Itoa(id)),
	})
	if err != nil {
		return nil, err
	}

	return o, nil
}

func PutThumbnail(id int, data io.Reader) error {
	_, err := s3client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:         aws.String("flatgrass-toybox-image"),
		Key:            aws.String(fmt.Sprintf("%d_thumb_128.png", id)),
		ACL:            types.ObjectCannedACLPublicRead,
		Body:           data,
		ChecksumSHA256: aws.String("UNSIGNED-PAYLOAD"), // required otherwise OVH S3 rejects
	})
	if err != nil {
		return err
	}

	return nil
}
