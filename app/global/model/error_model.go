package model

type ErrorLog struct {
	Line              string      `json:"line,omitempty"`
	Filename          string      `json:"filename,omitempty"`
	Function          string      `json:"function,omitempty"`
	Message           interface{} `json:"message,omitempty"`
	SystemMessage     interface{} `json:"system_message,omitempty"`
	Url               string      `json:"url,omitempty"`
	Method            string      `json:"method,omitempty"`
	Fields            interface{} `json:"fields,omitempty"`
	ConsumerTopic     string      `json:"consumer_topic,omitempty"`
	ConsumerPartition int         `json:"consumer_partition,omitempty"`
	ConsumerName      string      `json:"consumer_name,omitempty"`
	ConsumerOffset    int64       `json:"consumer_offset,omitempty"`
	ConsumerKey       string      `json:"consumer_key,omitempty"`
	Err               error       `json:"-"`
	StatusCode        int         `json:"-"`
}
