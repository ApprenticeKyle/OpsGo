package monitor

import (
	"OpsGo/internal/application/service/monitor"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MonitorHandler struct {
	service *monitor.MonitorService
}

func NewMonitorHandler(s *monitor.MonitorService) *MonitorHandler {
	return &MonitorHandler{service: s}
}

func (h *MonitorHandler) GetStats(c *gin.Context) {
	metrics, err := h.service.GetMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}
