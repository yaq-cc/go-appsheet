package logevent

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var encoder = base64.StdEncoding

type PubSubMessage struct {
	Message      Message `json:"message"`
	Subscription string  `json:"subscription"`
}

type Message struct {
	Attributes  map[string]interface{} `json:"attributes"`
	Data        *LoggingEvent          `json:"data"`
	MessageID   string                 `json:"messageId"`
	PublishTime time.Time              `json:"publishTime"`
}

// Proxy Object
type MarshalledMessage struct {
	Attributes  map[string]interface{} `json:"attributes"`
	Data        *Data                  `json:"data"`
	MessageID   string                 `json:"messageId"`
	PublishTime time.Time              `json:"publishTime"`
}

func (mm *MarshalledMessage) Transfer(m *Message) {
	m.Attributes = mm.Attributes
	m.MessageID = mm.MessageID
	m.PublishTime = mm.PublishTime
}

// Uses a proxy Object to accept "data" and then decodes it into a Message.
func (m *Message) UnmarshalJSON(in []byte) error {
	mm := &MarshalledMessage{}
	err := json.Unmarshal(in, mm)
	if err != nil {
		log.Fatal(err)
	}

	event := &LoggingEvent{}
	err = json.Unmarshal(*mm.Data, event)
	if err != nil {
		log.Fatal(err)
	}
	mm.Transfer(m)
	m.Data = event
	return nil
}

type Data []byte

func (d Data) String() string {
	return string(d)
}

func (d *Data) UnmarshalJSON(in []byte) error {
	in = in[1 : len(in)-1]
	out := make([]byte, encoder.DecodedLen(len(in)))
	n, err := encoder.Decode(out, in)
	if err != nil {
		return err
	}
	*d = out[:n]
	return nil
}

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
	PrincipalEmail        string `json:"principalEmail"`
	ServiceAccountKeyName string `json:"serviceAccountKeyName,omitempty"`
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
