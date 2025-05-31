package model

type ErrorLog struct {
	Line              string `json:"line,omitempty"`
	Filename          string `json:"filename,omitempty"`
	Function          string `json:"function,omitempty"`
	Message           string `json:"message,omitempty"`
	SystemMessage     string `json:"system_message,omitempty"`
	Url               string `json:"url,omitempty"`
	ConsumerTopic     string `json:"consumer_topic,omitempty"`
	ConsumerPartition int    `json:"consumer_partition,omitempty"`
	ConsumerName      string `json:"consumer_name,omitempty"`
	ConsumerOffset    int64  `json:"consumer_offset,omitempty"`
	ConsumerKey       string `json:"consumer_key,omitempty"`
	Err               error  `json:"-"`
	StatusCode        int    `json:"-"`
}
