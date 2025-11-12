package state

import (
	"sync"
	"time"
)

type State string

// Define all possible states
const (
	StateInitial              State = "Initial"
	StateMainMenu             State = "MainMenu"
	StateProfileView          State = "ProfileView"
	StateProfileSwitch        State = "ProfileSwitch"
	StateConfirmProfileSwitch State = "ConfirmProfileSwitch"
	StateViewGradesDetails    State = "ViewGradesDetails"
	StateViewGrades           State = "ViewGrades"
	StateViewBulletin         State = "ViewBulletin"
	StateViewDTR              State = "ViewDTR"
	StateViewPayables         State = "ViewPayables"
	StateProfileMenu          State = "ProfileMenu"
	StateViewSubjects         State = "ViewSubjects"
	StateSelectSubject        State = "SelectSubject"
	StateAskSupport           State = "AskSupport"
	StateViewTickets          State = "ViewTickets"
	StateSelectSupportTicket  State = "SelectSupportTicket"
)

// Key state
const (
	KeyProfileMap      string = "ProfileMap"
	KeySelectedProfile string = "SelectedProfile"
	KeySchoolYearMap   string = "KeySchoolYearMap"
	KeyThreadID        string = "KeyThreadID"
	KeyTicketMap       string = "KeyTicketMap"
	KeyPageOffset      string = "KeyPageOffset"
	KeyPaginationItems string = "PaginationItems"
	KeyPaginationTotal string = "PaginationTotal"
	KeyPaginationPage  string = "PaginationPage"
	KeyPaginationSize  string = "PaginationSize"
	KeyPaginationPages string = "PaginationPages"
)

type StateData struct {
	CurrentState State
	Data         map[string]interface{}
	LastActive   time.Time
}

type StateManager struct {
	mu     sync.RWMutex
	states map[string]*StateData
}

func NewStateManager() *StateManager {
	return &StateManager{
		states: make(map[string]*StateData),
	}
}

// SetState sets the current state for a user
func (sm *StateManager) SetState(userID string, state State, data map[string]any) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.states[userID]; !exists {
		sm.states[userID] = &StateData{
			Data: make(map[string]interface{}),
		}
	}

	sm.states[userID].CurrentState = state
	sm.states[userID].LastActive = time.Now()

	// Merge new data with existing data
	for k, v := range data {
		sm.states[userID].Data[k] = v
	}
	return nil
}

// GetState returns the current state and data for a user
func (sm *StateManager) GetState(userID string) (State, map[string]any) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if stateData, exists := sm.states[userID]; exists {
		dataCopy := make(map[string]any)
		for k, v := range stateData.Data {
			dataCopy[k] = v
		}
		return stateData.CurrentState, dataCopy
	}

	return StateInitial, nil
}

// ClearState clears the state for a user
func (sm *StateManager) ClearState(userID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.states, userID)
}

// CleanupInactive removes states that haven't been active for a duration
func (sm *StateManager) CleanupInactive(timeout time.Duration) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	for userID, stateData := range sm.states {
		if now.Sub(stateData.LastActive) > timeout {
			delete(sm.states, userID)
		}
	}
}
