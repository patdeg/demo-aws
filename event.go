package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

/**************************************************************
	Define /event Handler called from S3 Event Notification
 **************************************************************/
func eventHandler(w http.ResponseWriter, r *http.Request) {
	Info(">>>>> eventHandler")
	DebugInfo(r)
	PrintMemUsage()

	// Read event from http.Request
	event, err := ReadS3Event(r)
	if err != nil {
		Error("Error reading event: %v", err)
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	event.Print()

	// In case of "Notification", process each S3 files in Records
	if event.Type == "Notification" {
		for _, e := range event.MessageObject.Records {
			if e.EventName != "ObjectRemoved:Delete" {
				// Read S3 file
				data, err := ReadS3File(e.S3.Bucket.Name, e.S3.Object.Key)
				if err != nil {
					Error("Error with file s3://%v/%v", e.S3.Bucket.Name, e.S3.Object.Key)
				} else {
					// Process file content
					err = doWork(data, e.S3.Bucket.Name, e.S3.Object.Key)
					if err != nil {
						Error("Error doing work with file s3://%v/%v", e.S3.Bucket.Name, e.S3.Object.Key)
					}
				}
			}
		}
	}

	fmt.Fprintf(w, "OK")
}

/**************************************************************
	Do the work on the JSON content from s3://bucket/item
 **************************************************************/
func doWork(content []byte, bucket, item string) error {

	Debug("Working on data: %v", B2S(content))

	// Marshal content to a Go object
	var object DataObjectType
	err := json.Unmarshal(content, &object)
	if err != nil {
		Error("Error reading JSON data %v: %v", B2S(content), err)
		return err
	}

	// Execute the work, here add total = a + b, and set Timestamp to now in milliseconds
	for i := range object {
		object[i].Total = object[i].A + object[i].B
		object[i].Timestamp = time.Now().UnixNano() / 1000000 // TIMESTAMP_MILLIS
	}

	// Translate s3 item (e.g. data/test.json) into parquet item (e.g. processed/test.parquet)
	itemParquet := strings.Replace(item, "data", "processed", 1)
	itemParquet = strings.Replace(itemParquet, filepath.Ext(item), "", 1)
	itemParquet += ".parquet"

	// Translate s3 item (e.g. data/test.json) into error parking lot item (e.g. error/test.json)
	itemError := strings.Replace(item, "data", "error", 1)

	Debug("Raw filename s3://%v/%v", bucket, item)
	Debug("Processed filename s3://%v/%v", bucket, itemParquet)
	Debug("Error filename s3://%v/%v", bucket, itemError)

	// Write content to parquet file s3://bucket/itemParquet
	err = WriteToParquet(object, bucket, itemParquet)
	if err != nil {
		Error("Error processing file s3://%v/%v: %v", bucket, item, err)
		// In case of an error processing file, copy file to the error folder
		err2 := CopyS3File(bucket, item, bucket, itemError)
		if err2 != nil {
			Error("Error copying file s3://%v/%v to s3://%v/%v: %v", bucket, item, bucket, itemError, err2)
		}
		return err
	}

	Info("Parquet file s3://%v/%v ready", bucket, itemParquet)

	return nil
}

/**************************************************************
	Read S3 Event Notification from a http.Request call
 **************************************************************/
func ReadS3Event(r *http.Request) (*EventType, error) {

	// Read content for Body element of http.Request
	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(r.Body); err != nil {
		Debug("Error while dumping request: %v", err)
		return nil, err
	}
	body := buffer.Bytes()
	Debug("Response: %s", body)

	// interprate content with structure EventType
	var event EventType
	if err := json.Unmarshal(body, &event); err != nil {
		Error("Error decoding event: %v", err)
		return nil, err
	}

	// For a "Notification" event, interprate Message string element
	if event.Type == "Notification" {
		if err := json.Unmarshal([]byte(event.Message), &(event.MessageObject)); err != nil {
			Error("Error decoding event's message %v: %v", event.MessageObject, err)
			return nil, err
		}
	}

	return &event, nil
}

// Pretty Print an EventType
func (event *EventType) Print() {
	Debug("EVENT:")
	Debug("Type: %v", event.Type)
	Debug("SubscribeURL: %v", event.SubscribeURL)
	Debug("UnsubscribeURL: %v", event.UnsubscribeURL)
	Debug("Subject: %v", event.Subject)
	for _, e := range event.MessageObject.Records {
		Debug("%v on s3://%v/%v (%v bytes)", e.EventName, e.S3.Bucket.Name, e.S3.Object.Key, e.S3.Object.Size)
	}
}
