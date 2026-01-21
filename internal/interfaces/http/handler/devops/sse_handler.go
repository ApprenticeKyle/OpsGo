package devops

import (
	"io"

	"github.com/gin-gonic/gin"
)

func (h *DevOpsHandler) StreamLogs(c *gin.Context) {
	clientChan := h.devopsService.Broadcaster.Register()
	defer h.devopsService.Broadcaster.Unregister(clientChan)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("X-Accel-Buffering", "no") // Disable Nginx buffering
	// CORS handled by middleware but explicit set here for SSE if needed
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Send initial ping to confirm connection
	c.SSEvent("ping", "connected")
	c.Writer.Flush()

	c.Stream(func(w io.Writer) bool {
		if event, ok := <-clientChan; ok {
			c.SSEvent("message", event)
			return true
		}
		return false
	})
}
