package devops

import (
	"OpsGo/internal/domain/entity/devops"
	"OpsGo/internal/domain/repository"
	"context"

	"gorm.io/gorm"
)

type devopsRepository struct {
	db *gorm.DB
}

func NewDevOpsRepository(db *gorm.DB) repository.DevOpsRepository {
	return &devopsRepository{db: db}
}

func (r *devopsRepository) SaveConfig(ctx context.Context, config *devops.RepoConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

func (r *devopsRepository) GetConfig(ctx context.Context, id uint64) *devops.RepoConfig {
	var config devops.RepoConfig
	if err := r.db.WithContext(ctx).First(&config, id).Error; err != nil {
		return nil
	}
	return &config
}

func (r *devopsRepository) GetConfigByRepoURL(ctx context.Context, url string) *devops.RepoConfig {
	var config devops.RepoConfig
	if err := r.db.WithContext(ctx).Where("repo_url = ?", url).First(&config).Error; err != nil {
		return nil
	}
	return &config
}

func (r *devopsRepository) ListConfigs(ctx context.Context) ([]devops.RepoConfig, error) {
	var configs []devops.RepoConfig
	err := r.db.WithContext(ctx).Order("id desc").Find(&configs).Error
	return configs, err
}

func (r *devopsRepository) DeleteConfig(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&devops.RepoConfig{}, id).Error
}

func (r *devopsRepository) CreatePipelineRecord(ctx context.Context, record *devops.PipelineRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *devopsRepository) UpdatePipelineRecord(ctx context.Context, record *devops.PipelineRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}

func (r *devopsRepository) GetPipelineRecord(ctx context.Context, id uint64) *devops.PipelineRecord {
	var record devops.PipelineRecord
	if err := r.db.WithContext(ctx).First(&record, id).Error; err != nil {
		return nil
	}
	return &record
}

func (r *devopsRepository) ListPipelineRecords(ctx context.Context, limit int) ([]devops.PipelineRecord, error) {
	var records []devops.PipelineRecord
	err := r.db.WithContext(ctx).Order("id desc").Limit(limit).Find(&records).Error
	return records, err
}
