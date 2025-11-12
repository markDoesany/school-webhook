package handlers

import (
	"school-assistant-wh/internal/state"
	"time"

	"github.com/gin-gonic/gin"
)

type StateAwareHandler struct {
	stateManager *state.StateManager
	handler      func(c *gin.Context, userID string, currentState state.State, stateData map[string]interface{})
}

func NewStateAwareHandler(sm *state.StateManager, handler func(c *gin.Context, userID string, currentState state.State, stateData map[string]interface{})) *StateAwareHandler {
	return &StateAwareHandler{
		stateManager: sm,
		handler:      handler,
	}
}

func (h *StateAwareHandler) Handle(c *gin.Context) {
	// Get user ID from context (you'll need to set this in your auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userIDStr := userID.(string)
	currentState, stateData := h.stateManager.GetState(userIDStr)

	h.stateManager.SetState(userIDStr, currentState, nil)

	h.handler(c, userIDStr, currentState, stateData)
}

func StateCleanupMiddleware(sm *state.StateManager, cleanupInterval, stateTTL time.Duration) gin.HandlerFunc {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			sm.CleanupInactive(stateTTL)
		}
	}()

	return func(c *gin.Context) {
		c.Next()
	}
}
