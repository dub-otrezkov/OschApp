package api

type Submission struct {
	TaskId    string `json:"TaskId"`
	Answer    string `json:"Answer"`
	SessionId string `json:"SessionId"`
}
