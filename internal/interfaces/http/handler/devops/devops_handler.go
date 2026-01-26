package devops

import (
	"OpsGo/internal/application/dto"
	"OpsGo/internal/application/service/devops"
	"net/http"
	"strconv"

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

func (h *DevOpsHandler) DeleteConfig(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.devopsService.DeleteConfig(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

func (h *DevOpsHandler) GetSummary(c *gin.Context) {
	resp, err := h.devopsService.GetSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *DevOpsHandler) GetServiceLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	content, err := h.devopsService.GetServiceLog(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": content})
}


func (h *DevOpsHandler) HandleCICallback(c *gin.Context) {
	var req dto.CICallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	if err := h.devopsService.HandleCICallback(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deployment triggered from CI"})
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
