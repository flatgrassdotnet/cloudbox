/*
	cloudbox - the toybox server emulator
	Copyright (C) 2024-2025  patapancakes <patapancakes@pagefault.games>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

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
