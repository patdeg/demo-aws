package main

import (
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

/**************************************************************
	Write a parquet file in s3://s3_bucket/s3_item from
	a DataObjectType
 **************************************************************/
func WriteToParquet(object DataObjectType, s3_bucket, s3_item string) error {

	Debug("Preparing file s3://%v/%v", s3_bucket, s3_item)
	Debug("Object:%v", object)

	// Create temp folder
	exec_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error("Error getting executable folder: %v\n", err)
		return err
	}

	timeNano := strconv.FormatInt(time.Now().UnixNano(), 10)
	folder := exec_dir + "/data_" + timeNano
	if err := os.MkdirAll(folder, 0700); err != nil {
		Error("Error creating temp folder %v: %v\n", folder, err)
		return err
	}
	defer func() {
		_ = RemoveDirectory(folder)
	}()
	Debug("Working folder: %v", folder)

	// Create Parquet File Writer
	localParquetFilename := folder + "/" + filepath.Base(s3_item)
	Debug("Creating NewLocalFileWriter on local temp file %v", localParquetFilename)
	fw, err := local.NewLocalFileWriter(localParquetFilename)
	if err != nil {
		Error("Error: Can't create parquet file: %v", err)
		return err
	}
	defer fw.Close()

	pw, err := writer.NewParquetWriter(fw, new(DataObjectElement), 4)
	if err != nil {
		Error("Can't create json writer", err)
		return err
	}

	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	// Write data to Parquet with JSON content
	for _, element := range object {
		if err = pw.Write(element); err != nil {
			Error("Write error", err)
			return err
		}
	}

	// Stop Writer
	if err = pw.WriteStop(); err != nil {
		Error("WriteStop error", err)
		return err
	}

	Debug("Parquet file %v written", localParquetFilename)

	// Upload file to S3
	err = AddFileToS3(s3_bucket, s3_item, localParquetFilename)
	if err != nil {
		Error("Error adding file to S3", err)
		return err
	}

	Info("Parquet file s3://%v/%v ready", s3_bucket, s3_item)

	// Exiting will automatically RemoveDirectory(folder)
	return nil
}
