package account

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/handlers/utils"
	"school-assistant-wh/internal/models"
	"school-assistant-wh/internal/repositories"
	"school-assistant-wh/internal/services/facebook"
	"school-assistant-wh/internal/services/helpers"
	"school-assistant-wh/internal/state"
)

type AccountHandler struct {
	repo         repositories.UserRepository
	linkRepo     repositories.UserLinkRepository
	fbSvc        *facebook.Service
	utils        *utils.ResponseUtils
	stateManager *state.StateManager
}

func NewAccountHandler(repo repositories.UserRepository, linkRepo repositories.UserLinkRepository, fbSvc *facebook.Service, stateManager *state.StateManager) *AccountHandler {
	return &AccountHandler{
		repo:         repo,
		linkRepo:     linkRepo,
		fbSvc:        fbSvc,
		utils:        utils.NewResponseUtils(repo, linkRepo, fbSvc),
		stateManager: stateManager,
	}
}

func (h *AccountHandler) HandleRegistration(senderID string) error {
	user, err := h.repo.RegisterUser(senderID)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		return h.fbSvc.SendTextMessage(senderID, "Sorry, we couldn't complete your registration. Please try again.")
	}

	if err := h.fbSvc.SendTextMessage(senderID, "Congratulations! You are now registered to School Assistant"); err != nil {
		return fmt.Errorf("failed to send header message: %v", err)
	}

	welcomeMsg := fmt.Sprintf(constants.LinkAccountInstructions, *user.Code)
	return h.utils.SendResponseWithQuickReplies(senderID, welcomeMsg)
}

func (h *AccountHandler) HandleViewProfile(senderID string) error {
	// Get user's linked profiles
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	if user.IsActive == false {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	profiles, err := h.linkRepo.GetUserLinks(int(user.ID))
	if err != nil {
		return fmt.Errorf("failed to get user links: %v", err)
	}

	if len(profiles) == 0 {
		if err := h.fbSvc.SendTextMessage(senderID, "You don't have any linked accounts yet."); err != nil {
			return fmt.Errorf("failed to send header message: %v", err)
		}

		message := fmt.Sprintf(constants.LinkAccountInstructions, *user.Code)
		return h.utils.SendResponseWithQuickReplies(senderID, message)
	}

	// If there's only one profile and no primary set, ask for confirmation
	if len(profiles) == 1 && !profiles[0].IsPrimary {
		profile := profiles[0]
		message := fmt.Sprintf(
			constants.SingleProfileConfirmation,
			profile.Student.FirstName,
			safeString(&profile.Student.LastName),
			safeString(&profile.Student.StudentID),
			safeString(&profile.Student.School.SchoolName),
		)
		quickReplies := helpers.GetProfileConfirmationReplies()
		return h.fbSvc.SendQuickReplies(senderID, message, quickReplies)
	}

	// Check if there's a primary profile
	for _, profile := range profiles {
		if profile.IsPrimary && profile.Student != nil {
			message := fmt.Sprintf(
				constants.AccountPrimaryMessage,
				profile.Student.FirstName,
				profile.Student.LastName,
				profile.Student.StudentID,
				profile.Student.School.SchoolName,
			)

			var quickReplies []facebook.QuickReply
			if len(profiles) > 1 {
				quickReplies = helpers.GetProfileManagementReplies()
			} else {
				quickReplies = helpers.GetProfileConfirmationReplies()
			}
			return h.fbSvc.SendQuickReplies(senderID, message, quickReplies)
		}
	}
	profileMap := make(map[string]map[string]any)

	for i, profile := range profiles {
		key := strconv.Itoa(i + 1)
		profileMap[key] = map[string]any{
			"StudentID": profile.Student.StudentID,
			"SchoolID":  profile.Student.School.SchoolID,
		}
	}

	h.stateManager.SetState(senderID, state.StateProfileView, map[string]any{
		state.KeyProfileMap: profileMap,
	})

	return h.ShowAccountList(senderID, profiles, "Select a profile to switch to:")

}

func (h *AccountHandler) HandleSwitchProfile(senderID string) error {
	// Get user's linked profiles
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	if user.IsActive == false {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	profiles, err := h.linkRepo.GetUserLinks(int(user.ID))
	if err != nil {
		return fmt.Errorf("failed to get user links: %v", err)
	}

	if len(profiles) <= 1 {
		return h.utils.SendResponseWithQuickReplies(senderID,
			"You only have one linked profile. Contact your school admin to link more profiles.")
	}
	profileMap := make(map[string]map[string]any)

	for i, profile := range profiles {
		key := strconv.Itoa(i + 1)
		profileMap[key] = map[string]any{
			"StudentID": profile.Student.StudentID,
			"SchoolID":  profile.Student.School.SchoolID,
		}
	}

	h.stateManager.SetState(senderID, state.StateProfileSwitch, map[string]any{
		state.KeyProfileMap: profileMap,
	})

	return h.ShowAccountList(senderID, profiles, "Select a profile to switch to:")
}
func (h *AccountHandler) HandleViewSaID(senderID string) error {
	// Get user's linked profiles
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	if user.IsActive == false {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	if err := h.fbSvc.SendTextMessage(senderID, "Your School Assistant ID:"); err != nil {
		return fmt.Errorf("failed to send header message: %v", err)
	}

	return h.utils.SendResponseWithQuickReplies(senderID, fmt.Sprintf("%s", *user.Code))
}

// internal/handlers/account/handler.go
func (h *AccountHandler) HandleProfileSelection(senderID string, selectedProfile map[string]any) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	studentID := selectedProfile["StudentID"].(string)
	schoolID := selectedProfile["SchoolID"].(string)

	if err := h.linkRepo.UpdatePrimaryStatus(int(user.ID), studentID, schoolID); err != nil {
		return fmt.Errorf("failed to update primary profile: %w", err)
	}

	primaryProfile, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil {
		return fmt.Errorf("failed to get primary profile: %w", err)
	}

	if err := h.fbSvc.SendTextMessage(senderID,
		fmt.Sprintf(constants.AccountSwitchedMessage,
			primaryProfile.Student.FirstName,
			primaryProfile.Student.LastName,
		),
	); err != nil {
		log.Printf("Error sending account switched message: %v", err)
	}

	return nil
}

func (h *AccountHandler) HandleSelectCurrentProfile(senderID string) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return h.fbSvc.SendTextMessage(senderID, "Sorry, we couldn't complete your request. Please try again.")
	}

	if user.IsActive == false {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	profiles, err := h.linkRepo.GetUserLinks(int(user.ID))
	if err != nil || len(profiles) == 0 {
		log.Printf("Error getting user links or no profiles found: %v", err)
		return h.utils.SendResponseWithQuickReplies(senderID, "Sorry, we couldn't find any linked profiles. Please try again.")
	}

	if len(profiles) == 1 {
		err = h.linkRepo.UpdatePrimaryStatus(int(user.ID), profiles[0].StudentID, profiles[0].SchoolID)
		if err != nil {
			log.Printf("Error setting primary status: %v", err)
			return h.utils.SendResponseWithQuickReplies(senderID, "Sorry, we couldn't update your profile. Please try again.")
		}
	}

	var primaryProfile *models.UserLinkWithStudent
	for i := range profiles {
		if profiles[i].IsPrimary || len(profiles) == 1 {
			primaryProfile = &profiles[i]
			break
		}
	}

	if primaryProfile != nil && primaryProfile.Student != nil {
		message := fmt.Sprintf(constants.ProfileConfirmationMessage,
			primaryProfile.Student.FirstName,
			primaryProfile.Student.LastName,
			primaryProfile.Student.StudentID,
			primaryProfile.Student.School.SchoolName)
		err := h.fbSvc.SendTextMessage(senderID, message)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}

	return nil
}

func (h *AccountHandler) ShowAccountList(senderID string, profiles []models.UserLinkWithStudent, message string) error {
	var sb strings.Builder

	// Use message param if provided, otherwise fallback
	if message != "" {
		sb.WriteString(message + "\n\n")
	} else {
		sb.WriteString("Your Linked Profiles:\n")
	}

	for i, profile := range profiles {
		sb.WriteString(fmt.Sprintf("[%d] ", i+1))

		if profile.IsPrimary {
			sb.WriteString(fmt.Sprintf("%s %s (Primary)\n",
				safeString(&profile.Student.FirstName),
				safeString(&profile.Student.LastName)))
		} else {
			sb.WriteString(fmt.Sprintf("%s %s\n",
				safeString(&profile.Student.FirstName),
				safeString(&profile.Student.LastName)))
		}

		if profile.Student != nil {
			sb.WriteString(fmt.Sprintf("Student ID: %s\n", safeString(&profile.Student.StudentID)))

			if profile.Student.School != nil {
				sb.WriteString(fmt.Sprintf("School: %s\n", safeString(&profile.Student.School.SchoolName)))
			}

			if profile.Student.YearLevel != "" {
				sb.WriteString(fmt.Sprintf("Year: %s", safeString(&profile.Student.YearLevel)))
				if profile.Student.Course != "" {
					sb.WriteString(fmt.Sprintf(" - %s", safeString(&profile.Student.Course)))
				}
				sb.WriteString("\n")
			}
		}
		sb.WriteString("\n")
	}

	quickReplies := helpers.GetBack()
	return h.fbSvc.SendQuickReplies(senderID, sb.String(), quickReplies)
}

// Helper function to safely handle nil string pointers
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
