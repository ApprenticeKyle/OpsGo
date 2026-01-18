package devops

import (
	"OpsGo/internal/application/dto"
	"OpsGo/internal/application/service/devops"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DevOpsHandler struct {
	devopsService *devops.DevOpsService
}

func NewDevOpsHandler(devopsService *devops.DevOpsService) *DevOpsHandler {
	return &DevOpsHandler{
		devopsService: devopsService,
	}
}

func (h *DevOpsHandler) ConfigRepo(c *gin.Context) {
	var req dto.ConfigRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	resp, err := h.devopsService.ConfigRepo(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *DevOpsHandler) GetSummary(c *gin.Context) {
	resp, err := h.devopsService.GetSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *DevOpsHandler) HandleGitHubWebhook(c *gin.Context) {
	var ghPayload struct {
		Ref        string `json:"ref"`
		HeadCommit struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"head_commit"`
		Repository struct {
			HTMLURL string `json:"html_url"`
		} `json:"repository"`
	}

	if err := c.ShouldBindJSON(&ghPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid GitHub payload"})
		return
	}

	payload := dto.WebhookPayload{
		RepoURL:   ghPayload.Repository.HTMLURL,
		Ref:       ghPayload.Ref,
		CommitSHA: ghPayload.HeadCommit.ID,
		CommitMsg: ghPayload.HeadCommit.Message,
		Author:    ghPayload.HeadCommit.Author.Name,
		Status:    "success",
	}

	if err := h.devopsService.HandleWebhook(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *DevOpsHandler) TriggerDeployment(c *gin.Context) {
	var req struct {
		ConfigID uint64 `json:"config_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	if err := h.devopsService.TriggerDeployment(c.Request.Context(), req.ConfigID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deployment triggered"})
}
