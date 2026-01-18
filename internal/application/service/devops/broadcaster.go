package devops

import (
	"fmt"
	"sync"
)

type LogEvent struct {
	Type       string `json:"type"` // log, status
	PipelineID uint64 `json:"pipeline_id"`
	Content    string `json:"content,omitempty"`
	Status     string `json:"status,omitempty"`
}

type LogBroadcaster struct {
	clients    map[chan LogEvent]bool
	register   chan chan LogEvent
	unregister chan chan LogEvent
	broadcast  chan LogEvent
	mu         sync.RWMutex
}

func NewLogBroadcaster() *LogBroadcaster {
	lb := &LogBroadcaster{
		clients:    make(map[chan LogEvent]bool),
		register:   make(chan chan LogEvent),
		unregister: make(chan chan LogEvent),
		broadcast:  make(chan LogEvent),
	}
	go lb.run()
	return lb
}

func (lb *LogBroadcaster) run() {
	for {
		select {
		case client := <-lb.register:
			lb.mu.Lock()
			lb.clients[client] = true
			lb.mu.Unlock()
			fmt.Println("New SSE client registered")
		case client := <-lb.unregister:
			lb.mu.Lock()
			if _, ok := lb.clients[client]; ok {
				delete(lb.clients, client)
				close(client)
			}
			lb.mu.Unlock()
			fmt.Println("SSE client unregistered")
		case event := <-lb.broadcast:
			lb.mu.RLock()
			for client := range lb.clients {
				select {
				case client <- event:
				default:
					// If client buffer is full, we might want to drop or handle it
				}
			}
			lb.mu.RUnlock()
		}
	}
}

func (lb *LogBroadcaster) Register() chan LogEvent {
	client := make(chan LogEvent, 100)
	lb.register <- client
	return client
}

func (lb *LogBroadcaster) Unregister(client chan LogEvent) {
	lb.unregister <- client
}

func (lb *LogBroadcaster) BroadcastLog(pipelineID uint64, content string) {
	lb.broadcast <- LogEvent{
		Type:       "log",
		PipelineID: pipelineID,
		Content:    content,
	}
}

func (lb *LogBroadcaster) BroadcastStatus(pipelineID uint64, status string) {
	lb.broadcast <- LogEvent{
		Type:       "status",
		PipelineID: pipelineID,
		Status:     status,
	}
}
