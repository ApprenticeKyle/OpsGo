package devops

import "time"

type RepoConfig struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	RepoURL      string    `gorm:"size:255;not null" json:"repo_url"`
	DeployScript string    `gorm:"size:255;not null" json:"deploy_script"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (RepoConfig) TableName() string {
	return "repo_configs"
}
