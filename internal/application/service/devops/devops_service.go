package devops

import (
	"OpsGo/internal/application/dto"
	"OpsGo/internal/domain/entity/devops"
	"OpsGo/internal/domain/repository"
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type DevOpsService struct {
	repo        repository.DevOpsRepository
	Broadcaster *LogBroadcaster
}

func NewDevOpsService(repo repository.DevOpsRepository) *DevOpsService {
	return &DevOpsService{
		repo:        repo,
		Broadcaster: NewLogBroadcaster(),
	}
}

func (s *DevOpsService) ConfigRepo(ctx context.Context, req dto.ConfigRepoRequest) (*dto.ConfigRepoResponse, error) {
	config := &devops.RepoConfig{
		Name:         req.Name,
		RepoURL:      req.RepoURL,
		DeployScript: req.DeployScript,
		LogPath:      req.LogPath,
	}

	// Check if exists
	existing := s.repo.GetConfigByRepoURL(ctx, req.RepoURL)
	if existing != nil {
		config.ID = existing.ID
	}

	if err := s.repo.SaveConfig(ctx, config); err != nil {
		return nil, err
	}

	return &dto.ConfigRepoResponse{
		ID:           config.ID,
		RepoURL:      config.RepoURL,
		DeployScript: config.DeployScript,
	}, nil
}

func (s *DevOpsService) DeleteConfig(ctx context.Context, id uint64) error {
	// Check if exists
	config := s.repo.GetConfig(ctx, id)
	if config == nil {
		return fmt.Errorf("config not found")
	}
	return s.repo.DeleteConfig(ctx, id)
}

func (s *DevOpsService) GetServiceLog(ctx context.Context, configID uint64) (string, error) {
	config := s.repo.GetConfig(ctx, configID)
	if config == nil {
		return "", fmt.Errorf("config not found")
	}

	if config.LogPath == "" {
		return "", fmt.Errorf("log path not configured")
	}

	// Read last 5KB of log file
	file, err := os.Open(config.LogPath)
	if err != nil {
		return "", fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	fileSize := stat.Size()
	readSize := int64(5120) // 5KB
	if fileSize < readSize {
		readSize = fileSize
	}

	offset := fileSize - readSize
	if offset < 0 {
		offset = 0
	}

	buf := make([]byte, readSize)
	_, err = file.ReadAt(buf, offset)
	if err != nil && err != io.EOF {
		return "", err
	}

	return string(buf), nil
}

func (s *DevOpsService) GetSummary(ctx context.Context) (*dto.DevOpsSummaryResponse, error) {
	configs, err := s.repo.ListConfigs(ctx)
	if err != nil {
		return nil, err
	}

	records, err := s.repo.ListPipelineRecords(ctx, 10)
	if err != nil {
		return nil, err
	}

	var services []dto.ConfigRepoResponse
	for _, c := range configs {
		services = append(services, dto.ConfigRepoResponse{
			ID:           c.ID,
			Name:         c.Name,
			RepoURL:      c.RepoURL,
			DeployScript: c.DeployScript,
			LogPath:      c.LogPath,
		})
	}

	var pipelines []dto.PipelineRecordResponse
	for _, p := range records {
		pipelines = append(pipelines, dto.PipelineRecordResponse{
			ID:            p.ID,
			RepoName:      p.RepoName,
			Status:        p.Status,
			Ref:           p.Ref,
			CommitSHA:     p.CommitSHA,
			CommitMsg:     p.CommitMsg,
			Author:        p.Author,
			TriggerSource: p.TriggerSource,
			Duration:      p.Duration,
			StartedAt:     p.StartedAt,
			FinishedAt:    p.FinishedAt,
			CreatedAt:     p.CreatedAt,
		})
	}

	return &dto.DevOpsSummaryResponse{
		Services:  services,
		Pipelines: pipelines,
	}, nil
}


func (s *DevOpsService) TriggerDeployment(ctx context.Context, configID uint64) error {
	config := s.repo.GetConfig(ctx, configID)
	if config == nil {
		return fmt.Errorf("config not found")
	}

	record := &devops.PipelineRecord{
		ConfigID:      config.ID,
		RepoName:      config.Name,
		Status:        "pending",
		Ref:           "manual",
		TriggerSource: "manual",
		CreatedAt:     time.Now(),
	}

	if err := s.repo.CreatePipelineRecord(ctx, record); err != nil {
		return err
	}

	// Trigger async deployment
	go s.runDeployment(record.ID, config.DeployScript)

	return nil
}

func (s *DevOpsService) HandleCICallback(ctx context.Context, req dto.CICallbackRequest) error {
	if req.Status != "success" {
		return fmt.Errorf("CI build failed, skipping deployment")
	}

	config := s.repo.GetConfigByRepoURL(ctx, req.RepoURL)
	if config == nil {
		return fmt.Errorf("repository not configured: %s", req.RepoURL)
	}

	record := &devops.PipelineRecord{
		ConfigID:      config.ID,
		RepoName:      config.Name,
		Status:        "pending",
		Ref:           req.Tag,
		CommitSHA:     req.CommitSHA,
		TriggerSource: "ci_cd",
		CreatedAt:     time.Now(),
	}

	if err := s.repo.CreatePipelineRecord(ctx, record); err != nil {
		return err
	}

	// Trigger async deployment
	go s.runDeployment(record.ID, config.DeployScript)

	return nil
}

func (s *DevOpsService) runDeployment(recordID uint64, scriptPath string) {
	ctx := context.Background()
	startTime := time.Now()

	s.updateRecordStatus(ctx, recordID, "running", &startTime, nil, "")
	s.Broadcaster.BroadcastStatus(recordID, "running")

	cmd := exec.Command("/bin/bash", scriptPath)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	multi := io.MultiReader(stdout, stderr)

	if err := cmd.Start(); err != nil {
		now := time.Now()
		msg := fmt.Sprintf("Failed to start script: %v", err)
		s.updateRecordStatus(ctx, recordID, "failed", nil, &now, msg)
		s.Broadcaster.BroadcastLog(recordID, msg)
		s.Broadcaster.BroadcastStatus(recordID, "failed")
		return
	}

	// Stream logs
	reader := bufio.NewReader(multi)
	fullLog := ""
	for {
		line, err := reader.ReadString('\n')
		if line != "" {
			fullLog += line
			s.Broadcaster.BroadcastLog(recordID, line)
		}
		if err != nil {
			break
		}
	}

	err := cmd.Wait()
	finishTime := time.Now()
	status := "success"
	if err != nil {
		status = "failed"
		errMsg := fmt.Sprintf("\nCommand failed: %v\n", err)
		fullLog += errMsg
		s.Broadcaster.BroadcastLog(recordID, errMsg)
	}

	s.updateRecordStatus(ctx, recordID, status, nil, &finishTime, fullLog)
	s.Broadcaster.BroadcastStatus(recordID, status)
}

func (s *DevOpsService) updateRecordStatus(ctx context.Context, id uint64, status string, start *time.Time, finish *time.Time, logContent string) {
	record := s.repo.GetPipelineRecord(ctx, id)
	if record == nil {
		return
	}

	record.Status = status
	if start != nil {
		record.StartedAt = start
	}
	if finish != nil {
		record.FinishedAt = finish
		if record.StartedAt != nil {
			record.Duration = int64(finish.Sub(*record.StartedAt).Seconds())
		}
	}
	if logContent != "" {
		record.CommitMsg = logContent
	}

	s.repo.UpdatePipelineRecord(ctx, record)
}
