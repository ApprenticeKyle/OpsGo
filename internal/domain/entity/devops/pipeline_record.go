package devops

import "time"

type PipelineRecord struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigID      uint64     `json:"config_id"`
	RepoName      string     `gorm:"size:100" json:"repo_name"`
	Status        string     `gorm:"size:20;default:'pending'" json:"status"` // pending, running, success, failed, canceled
	Ref           string     `gorm:"size:100" json:"ref"`                     // branch or tag
	CommitSHA     string     `gorm:"size:40" json:"commit_sha"`
	CommitMsg     string     `gorm:"type:text" json:"commit_msg"`
	Author        string     `gorm:"size:100" json:"author"`
	TriggerSource string     `gorm:"size:20;default:'manual'" json:"trigger_source"` // manual, webhook
	Duration      int64      `json:"duration"`                                       // seconds
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (PipelineRecord) TableName() string {
	return "devops_pipeline_records"
}
