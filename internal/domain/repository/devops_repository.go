package repository

import (
	"OpsGo/internal/domain/entity/devops"
	"context"
)

type DevOpsRepository interface {
	SaveConfig(ctx context.Context, config *devops.RepoConfig) error
	GetConfig(ctx context.Context, id uint64) *devops.RepoConfig
	GetConfigByRepoURL(ctx context.Context, url string) *devops.RepoConfig
	ListConfigs(ctx context.Context) ([]devops.RepoConfig, error)

	CreatePipelineRecord(ctx context.Context, record *devops.PipelineRecord) error
	UpdatePipelineRecord(ctx context.Context, record *devops.PipelineRecord) error
	GetPipelineRecord(ctx context.Context, id uint64) *devops.PipelineRecord
	ListPipelineRecords(ctx context.Context, limit int) ([]devops.PipelineRecord, error)
}
