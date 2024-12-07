package models

const (
	Pending    = "pending"
	Processing = "processing"
	Completed  = "completed"
)

type JobCreateDto struct {
	Payload string `json:"payload"`
}

type JobIDDto struct {
	JobID int64 `json:"job_id"`
}

type JobStatusDto struct {
	JobID  int64  `json:"job_id"`
	Status string `json:"status"`
}

type JobTask struct {
	JobID   int64
	Payload string
}

type Job struct {
	Payload string
}
