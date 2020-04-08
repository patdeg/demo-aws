package main

import (
	"time"
)

// Root structure of the JSON file being loaded to S3
type DataObjectType []DataObjectElement

// Line Element of the JSON file being loaded to S3
type DataObjectElement struct {
	A         float32 `json:"a,omitempty" parquet:"name=a, type=FLOAT"`
	B         float32 `json:"b,omitempty" parquet:"name=b, type=FLOAT"`
	Total     float32 `json:"total,omitempty" parquet:"name=total, type=FLOAT"`
	Timestamp int64   `json:"created_ts,omitempty" parquet:"name=created_ts, type=TIMESTAMP_MILLIS"`
}

// Root S3 Event Notification
type EventType struct {
	Type             string      `json:"Type,omitempty"`
	MessageId        string      `json:"MessageId,omitempty"`
	Token            string      `json:"Token,omitempty"`
	TopicArn         string      `json:"TopicArn,omitempty"`
	Subject          string      `json:"Subject,omitempty"`
	Message          string      `json:"Message,omitempty"`
	MessageObject    MessageType `json:"-"`
	SubscribeURL     string      `json:"SubscribeURL,omitempty"`
	UnsubscribeURL   string      `json:"UnsubscribeURL,omitempty"`
	Timestamp        time.Time   `json:"Timestamp,omitempty"`
	SignatureVersion string      `json:"SignatureVersion,omitempty"`
	Signature        string      `json:"Signature,omitempty"`
	SigningCertURL   string      `json:"SigningCertURL,omitempty"`
}

// Interprated Message element into a Go structure
type MessageType struct {
	Records []RecordType `json:"Records,omitempty"`
}

type RecordType struct {
	EventVersion      string                `json:"eventVersion,omitempty"` // 2.1
	EventSource       string                `json:"eventSource,omitempty"`  // aws:s3
	AwsRegion         string                `json:"awsRegion,omitempty"`    // us-west-1
	EventTime         time.Time             `json:"eventTime,omitempty"`    // 2020-04-06T21:05:44.149Z
	EventName         string                `json:"eventName,omitempty"`    // ObjectCreated:Put
	UserIdentity      UserIdentityType      `json:"userIdentity,omitempty"`
	RequestParameters RequestParametersType `json:"requestParameters,omitempty"`
	ResponseElements  ResponseElementsType  `json:"responseElements,omitempty"`
	S3                S3Type                `json:"s3,omitempty"`
}

type UserIdentityType struct {
	PrincipalId string `json:"principalId,omitempty"` // AJQQI9EMTKHP2
}

type RequestParametersType struct {
	SourceIPAddress string `json:"sourceIPAddress,omitempty"` // 70.179.8.27
}

type ResponseElementsType struct {
	XAmzRequestId string `json:"x-amz-request-id,omitempty"` // 6CE52B06A4136BCF
	XAmzId2       string `json:"x-amz-id-2,omitempty"`       // gE+sMwgy35BhKkE1Feq9GZkpG1D3AAmDZm7BB3eBr3H6Rr5F+dHJ7U76NetTWeuDugP3Fk8K1AUKsKztdh5CesK5FvR1nuF2
}

type S3Type struct {
	S3SchemaVersion string     `json:"s3SchemaVersion,omitempty"` // 1.0
	ConfigurationId string     `json:"configurationId,omitempty"` // MyEvent
	Bucket          BucketType `json:"bucket,omitempty"`
	Object          ObjectType `json:"object,omitempty"`
}

type BucketType struct {
	Name          string           `json:"name,omitempty"` // deglon
	OwnerIdentity UserIdentityType `json:"ownerIdentity,omitempty"`
	Arn           string           `json:"arn,omitempty"` // arn:aws:s3:::deglon
}

type ObjectType struct {
	Key       string `json:"key,omitempty"`       // data/test2.json
	Size      int64  `json:"size,omitempty"`      // 15
	ETag      string `json:"eTag,omitempty"`      // 7185811e96191f0ef5c6830643eaa3d0
	Sequencer string `json:"sequencer,omitempty"` // 005E8B99AA4CE3A3D2
}
