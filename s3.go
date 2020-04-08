package main

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/http"
	"os"
)

/**************************************************************
	Read a file from s3://bucket/item
 **************************************************************/
func ReadS3File(bucket, item string) ([]byte, error) {

	// Start a session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	// Prepare a WriteAt Buffer for s3
	buf := aws.NewWriteAtBuffer([]byte{})

	// Download data from s3
	numBytes, err := s3manager.NewDownloader(sess).Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		Error("Unable to download item %q, %v", item, err)
		return []byte{}, err
	}

	Debug("Downloaded s3://%v/%v with %v bytes", bucket, item, numBytes)

	return buf.Bytes(), nil
}

/**************************************************************
	AddFileToS3 will upload a single file to S3, it will 
	require a pre-built aws session and will set file info 
	like content type and encryption on the uploaded file.
 **************************************************************/
func AddFileToS3(s3_bucket, s3_item, fileName string) error {

	// Start a session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	// Open the file for use
	file, err := os.Open(fileName)
	if err != nil {
		Error("Error opening file %v: %v", fileName, err)
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	if _,err = file.Read(buffer); err != nil {
		Error("Error reading file %v: %v", fileName, err)
		return err
	}

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s3_bucket),
		Key:                  aws.String(s3_item),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

/**************************************************************
	Copy a file s3://source_bucket/source_item to
	s3://bucket/item
 **************************************************************/
func CopyS3File(source_bucket, source_item, bucket, item string) error {

	// Start a session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	// Copy the item
	if _, err := s3.New(sess).CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(source_bucket + "/" + source_item),
		Key:        aws.String(item),
	}); err != nil {
		Error("Unable to copy item from bucket %v to bucket %v: %v", source_bucket, bucket, err)
		return err
	}

	// Wait to see if the item got copied
	if err = s3.New(sess).WaitUntilObjectExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	}); err != nil {
		Error("Error occurred while waiting for item %v to be copied to bucket %v: %v", source_item, bucket, err)
		return err
	}
	
	return nil
}
