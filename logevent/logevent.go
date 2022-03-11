package logevent

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type LoggingEvent struct {
	InsertID         string       `json:"insertId"`
	LogName          string       `json:"logName"`
	ProtoPayload     ProtoPayload `json:"protoPayload"`
	ReceiveTimestamp time.Time    `json:"receiveTimestamp"`
	Resource         Resource     `json:"resource"`
	Severity         string       `json:"severity"`
	Timestamp        time.Time    `json:"timestamp"`
}

func FromRequest(r *http.Request) (*LoggingEvent, error) {
	var e LoggingEvent
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		return nil, err
	} else {
		return &e, nil
	}
}

func FromReader(r io.Reader) (*LoggingEvent, error) {
	var e LoggingEvent
	err := json.NewDecoder(r).Decode(&e)
	if err != nil {
		return nil, err
	} else {
		return &e, nil
	}
}

func (e *LoggingEvent) GetResourceName() string {
	return e.ProtoPayload.ResourceName
}


// Example: "projects/_/buckets/holy-diver-297719-appsheet-docai/objects/DocId_1fqXlSIif8F7aSdkVfEEff8NjWkdZqftBRZ7OZA7Z9EM/invoices_Files_/8F50B70D_78F2_4B25_9BA6_ADFD0A413D88.invoice_file.135228.pdf"
// Desired output: bucket, object
func (e *LoggingEvent) GetObjectData() (bkt string, obj string, key string) {
	resourceName := e.GetResourceName()
	resourceParts := strings.Split(resourceName, "/")
	bkt = resourceParts[3]
	obj = strings.Join(resourceParts[5:], "/")
	fileParts := strings.Split(resourceParts[7], ".")
	key = fileParts[0]
	return bkt, obj, key
}

type ProtoPayload struct {
	Type               string              `json:"@type"`
	AuthenticationInfo AuthenticationInfo  `json:"authenticationInfo"`
	AuthorizationInfo  []AuthorizationInfo `json:"authorizationInfo"`
	MethodName         string              `json:"methodName"`
	RequestMetadata    RequestMetadata     `json:"requestMetadata"`
	ResourceLocation   ResourceLocation    `json:"resourceLocation"`
	ResourceName       string              `json:"resourceName"`
	ServiceName        string              `json:"serviceName"`
	Status             Status              `json:"status"`
}

type AuthenticationInfo struct {
	PrincipalEmail string `json:"principalEmail"`
}

type AuthorizationInfo struct {
	Granted            bool               `json:"granted"`
	Permission         string             `json:"permission"`
	Resource           string             `json:"resource"`
	ResourceAttributes ResourceAttributes `json:"resourceAttributes"`
}

type ResourceAttributes struct {
}

type RequestMetadata struct {
	CallerIP                string                `json:"callerIp"`
	CallerSuppliedUserAgent string                `json:"callerSuppliedUserAgent"`
	DestinationAttributes   DestinationAttributes `json:"destinationAttributes"`
	RequestAttributes       RequestAttributes     `json:"requestAttributes"`
}

type DestinationAttributes struct {
}

type RequestAttributes struct {
	Auth Auth      `json:"auth"`
	Time time.Time `json:"time"`
}

type Auth struct {
}

type ResourceLocation struct {
	CurrentLocations []string `json:"currentLocations"`
}

type Status struct {
}

type Resource struct {
	Labels Labels `json:"labels"`
	Type   string `json:"type"`
}

type Labels struct {
	BucketName string `json:"bucket_name"`
	Location   string `json:"location"`
	ProjectID  string `json:"project_id"`
}
