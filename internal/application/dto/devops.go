package dto

import "time"

type ConfigRepoRequest struct {
	RepoURL      string `json:"repo_url" binding:"required"`
	DeployScript string `json:"deploy_script" binding:"required"`
	Name         string `json:"name" binding:"required"`
}

type ConfigRepoResponse struct {
	ID           uint64 `json:"id"`
	RepoURL      string `json:"repo_url"`
	DeployScript string `json:"deploy_script"`
	Name         string `json:"name"`
}

type PipelineRecordResponse struct {
	ID            uint64     `json:"id"`
	RepoName      string     `json:"repo_name"`
	Status        string     `json:"status"`
	Ref           string     `json:"ref"`
	CommitSHA     string     `json:"commit_sha"`
	CommitMsg     string     `json:"commit_msg"`
	Author        string     `json:"author"`
	TriggerSource string     `json:"trigger_source"`
	Duration      int64      `json:"duration"`
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type DevOpsSummaryResponse struct {
	Services  []ConfigRepoResponse     `json:"services"`
	Pipelines []PipelineRecordResponse `json:"pipelines"`
}

type WebhookPayload struct {
	RepoURL   string `json:"repo_url"`
	Ref       string `json:"ref"`
	CommitSHA string `json:"commit_sha"`
	CommitMsg string `json:"commit_msg"`
	Author    string `json:"author"`
	Status    string `json:"status"`
}
