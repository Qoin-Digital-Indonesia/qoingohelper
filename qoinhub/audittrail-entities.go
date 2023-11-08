package qoingohelper

const (
	AUDITTRAIL_PROCESS = "audit-trail-process" // command rabbit mq audittrail process
	AUDITTRAIL_DATA    = "audit-trail-data"    // command rabbit mq audittrail data
	TIME_FORMAT        = "2006-01-02 15:04:05" // time format
)

type MessagePayloadAudit struct {
	Id       int         `json:"Id"`
	Command  string      `json:"Command"`
	Time     string      `json:"Time"`
	ModuleId string      `json:"ModuleId"`
	Data     interface{} `json:"Data"`
}

type AuditTrialProcess struct {
	Subject     string                `json:"Subject,omitempty"`     // app name or service name (required)
	Function    string                `json:"Function,omitempty"`    // name of function the audittrail will be implemented (required)
	Description string                `json:"Description,omitempty"` // what purpose of the function (required)
	Key         []interface{}         `json:"Key"`                   // like user id, merchant code, trx code .etc (required)
	Data        DataAudittrailProcess `json:"Data"`                  // DataAudittrailProcess struct (required)
}

type DataAudittrailProcess struct {
	Time string `json:"Time"` // time will be handle in library
	Info string `json:"Info"` // message from service/app want to print in log (required)
}

type AuditTrialData struct {
	Subject           string              `json:"Subject,omitempty"`     // app name or service name (required)
	Function          string              `json:"Function,omitempty"`    // name of function the audittrail will be implemented (required)
	Description       string              `json:"Description,omitempty"` // what purpose of the function (required)
	Key               []interface{}       `json:"Key"`                   // // like user id, merchant code, trx code .etc (required)
	Source            string              `json:"Source"`                // internal or external (required)
	CommunicationType string              `json:"CommunicationType"`     // like rabbit mq, grcp or Rest API (required)
	Data              *RequestAndResponse `json:"Data"`                  // request outgoing or service/app receive from client and response from target or return from our service/app to client (required)
}

type RequestAndResponse struct {
	Request  Request       `json:"Request"`  // request outgoing or service/app receive from client
	Response ResponseAudit `json:"Response"` // response from target or return from our service/app to client
}

type Request struct {
	Time        string      `json:"Time"`                  // time will be handle in library
	Path        string      `json:"Path,omitempty"`        // ex: /merchantpg/onboarding(resp api), <funcName>(grpc), <queue/exchange>|command(rabbitmq) (required)
	QueryString interface{} `json:"QueryString,omitempty"` // ex:  (optional)
	Header      interface{} `json:"Header,omitempty"`      // optional except RestApi incoming or outgoing
	Param       interface{} `json:"Param,omitempty"`       // ex:{ "id" : 1234} optional
	Body        interface{} `json:"Body,omitempty"`        // request body from client or our service to target (optional)
	IpAddress   string      `json:"IpAddress,omitempty"`   // optional except RestApi Incoming
	BrowserId   int         `json:"BrowserId,omitempty"`   // optional
	Latitude    string      `json:"Latitude,omitempty"`    // optional
	Longitude   string      `json:"Longitude,omitempty"`   // optional
}

type ResponseAudit struct {
	Time   string `json:"Time"` // time will be handle in library
	Detail Detail `json:"Detail,omitempty"`
}

type Detail struct {
	StatusCode int         `json:"StatusCode"`
	Message    string      `json:"Message"`
	Data       interface{} `json:"Data,omitempty"`
}
