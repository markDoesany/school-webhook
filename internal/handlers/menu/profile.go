package menu

import (
	"fmt"
	"log"
	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/services/helpers"
	"school-assistant-wh/internal/state"
)

// HandleProfileDetails displays the user's profile details and options
func (h *MenuHandler) ShowProfileMenu(senderID string) error {
	// Get user profile
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}
	currentProfileData, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil {
		return fmt.Errorf("failed to get user links: %w", err)
	}

	// Format the profile details message
	status := "Active"
	if !currentProfileData.IsActive {
		status = "Inactive"
	}

	fullName := fmt.Sprintf("%s %s", currentProfileData.Student.FirstName, currentProfileData.Student.LastName)
	message := fmt.Sprintf(
		constants.ProfileDetailsTemplate,
		fullName,
		currentProfileData.Student.Course,
		currentProfileData.Student.YearLevel,
		status,
	)
	if err := h.stateManager.SetState(senderID, state.StateProfileMenu, nil); err != nil {
		log.Printf("Error setting state: %v", err)
	}

	quickReplies := helpers.GetBack()
	return h.fbSvc.SendQuickReplies(senderID, message, quickReplies)
}

// HandleProfileMenuSelection processes the user's selection from the profile menu
func (h *MenuHandler) HandleProfileMenuSelection(senderID, selection string) error {
	switch selection {
	case "1": // Subjects Enrolled
		return h.HandleViewSubjects(senderID)
	case "2": // Switch Accounts
		if err := h.stateManager.SetState(senderID, state.StateConfirmProfileSwitch, nil); err != nil {
			log.Printf("Error setting switch profile state: %v", err)
		}
		quickReplies := helpers.GetConfirmProfileSwitch()
		return h.fbSvc.SendQuickReplies(senderID, "Please confirm you want to switch accounts.", quickReplies)
	default:
		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(
			senderID,
			"⚠️ Invalid selection. Please choose a valid option.",
			quickReplies,
		)
	}
}
