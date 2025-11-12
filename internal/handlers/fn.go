package handlers

import (
	"fmt"
	"log"

	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/services/helpers"
	"school-assistant-wh/internal/state"
)

// handleProfileSelection handles the profile selection logic
func (h *Handler) handleProfileSelection(senderID, message string, stateData map[string]any) error {
	if profileMap, ok := stateData[state.KeyProfileMap].(map[string]map[string]any); ok {
		if selectedProfile, exists := profileMap[message]; exists {
			if err := h.accountHdlr.HandleProfileSelection(senderID, selectedProfile); err != nil {
				return err
			}
			if err := h.fbSvc.SendTextMessage(senderID, constants.WelcomeMessageAboard); err != nil {
				return fmt.Errorf("failed to send welcome message: %w", err)
			}
			if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			err := h.menuHdlr.ShowMainMenu(senderID)
			if err != nil {
				return err
			}
			return nil
		}
		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(senderID, "Invalid selection. Please try again.", quickReplies)
	}
	return nil
}

// handleContinue processes the CONTINUE action
func (h *Handler) handleContinue(senderID string) error {
	err := h.accountHdlr.HandleSelectCurrentProfile(senderID)
	if err != nil {
		return err
	}

	if err := h.fbSvc.SendTextMessage(senderID, constants.WelcomeMessageAboard); err != nil {
		return fmt.Errorf("failed to send welcome message: %w", err)
	}
	if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
		log.Printf("Error resetting state: %v", err)
	}
	err = h.menuHdlr.ShowMainMenu(senderID)
	if err != nil {
		return err
	}
	return nil
}

// handleNo processes the NO action
func (h *Handler) handleNo(senderID string) error {
	return h.utils.SendResponseWithQuickReplies(senderID, "No problem! You can always access your profile later by typing 'View Profile'.")
}

// handleAboutUs shows the about us information
func (h *Handler) handleAboutUs(senderID string) error {
	return h.utils.SendResponseWithQuickReplies(senderID, constants.AboutUsMessage)
}

// handleTalkToHuman provides contact information for human assistance
func (h *Handler) handleTalkToHuman(senderID string) error {
	return h.utils.SendResponseWithQuickReplies(senderID, constants.TalkToHumanMessage)
}

// handleGetStarted shows the initial welcome message
func (h *Handler) handleGetStarted(senderID string) error {
	return h.utils.SendResponseWithQuickReplies(senderID, constants.WelcomeMessage)
}

// handleDefault provides a fallback response for unrecognized messages
func (h *Handler) handleDefault(senderID string) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return h.sendDefaultMessage(senderID)
	}
	fmt.Println(user)
	if user.IsActive == false {
		quickReplies := helpers.GetQuickReplies(constants.AccountDeactivatedMessage)
		return h.fbSvc.SendQuickReplies(senderID, constants.AccountDeactivatedMessage, quickReplies)
	}

	_, err = h.linkRepo.GetPrimaryLink(int(user.ID))
	if err == nil {
		if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
			log.Printf("Error resetting state: %v", err)
		}
		return h.menuHdlr.ShowMainMenu(senderID)
	}

	return h.sendDefaultMessage(senderID)
}

// sendDefaultMessage sends the default "I'm not sure I understand" message
func (h *Handler) sendDefaultMessage(senderID string) error {
	defaultMsg := "I'm not sure I understand. Here are the available options:"
	if err := h.fbSvc.SendTextMessage(senderID, defaultMsg); err != nil {
		log.Printf("Error sending default message: %v", err)
	}
	return h.utils.SendResponseWithQuickReplies(senderID, "How can I help you today?")
}

// handleSchoolYearSelection handles the school year selection for viewing grades
func (h *Handler) handleSchoolYearSelection(senderID, message string, stateData map[string]any) error {
	if yearMap, ok := stateData[state.KeySchoolYearMap].(map[string]string); ok {
		if year, exists := yearMap[message]; exists {
			return h.menuHdlr.HandleViewGradesByYear(senderID, year)
		}
		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(senderID,
			" Invalid selection. Please choose a valid school year from the options above.",
			quickReplies,
		)
	}
	if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
		log.Printf("Error resetting state: %v", err)
	}
	return h.menuHdlr.ShowMainMenu(senderID)
}

// handleSubjectByYearSelection handles the school year selection for viewing subjects
func (h *Handler) handleSubjectByYearSelection(senderID, message string, stateData map[string]any) error {
	if yearMap, ok := stateData[state.KeySchoolYearMap].(map[string]string); ok {
		if year, exists := yearMap[message]; exists {
			return h.menuHdlr.HandleViewSubjectsByYear(senderID, year)
		}

		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(senderID,
			" Invalid selection. Please choose a valid school year from the options above.",
			quickReplies,
		)
	}

	// If we get here, there was an issue with the state data
	if err := h.stateManager.SetState(senderID, state.StateProfileMenu, nil); err != nil {
		log.Printf("Error resetting state: %v", err)
	}
	return h.menuHdlr.ShowProfileMenu(senderID)
}

// handleSupportTicketSelection handles the support ticket selection
func (h *Handler) handleSupportTicketSelection(senderID, message string, stateData map[string]any) error {
	if ticketMap, ok := stateData[state.KeyTicketMap].(map[string]string); ok {
		if ticketID, exists := ticketMap[message]; exists {
			return h.menuHdlr.HandleSupportTicketSelection(senderID, ticketID)
		}

		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(senderID,
			"Invalid selection. Please choose a valid ticket from the options above.",
			quickReplies,
		)
	}
	// If we get here, there was an issue with the state data
	if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
		log.Printf("Error resetting state: %v", err)
	}
	return h.menuHdlr.ShowMainMenu(senderID)
}
