package menu

import (
	"fmt"
	"log"
	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/handlers/utils"
	"school-assistant-wh/internal/repositories"
	"school-assistant-wh/internal/services/facebook"
	"school-assistant-wh/internal/services/helpers"
	"school-assistant-wh/internal/state"
)

type MenuHandler struct {
	repo           repositories.UserRepository
	linkRepo       repositories.UserLinkRepository
	gradeRepo      repositories.GradeRepository
	bulletinRepo   repositories.BulletinRepository
	payableRepo    repositories.StudentPayableRepository
	paymentLogRepo repositories.PaymentLogRepository
	dtrRepo        repositories.DTRRepository
	supportRepo    repositories.SupportRepository
	fbSvc          *facebook.Service
	utils          *utils.ResponseUtils
	stateManager   *state.StateManager
}

func NewMenuHandler(
	repo repositories.UserRepository,
	linkRepo repositories.UserLinkRepository,
	gradeRepo repositories.GradeRepository,
	bulletinRepo repositories.BulletinRepository,
	payableRepo repositories.StudentPayableRepository,
	paymentLogRepo repositories.PaymentLogRepository,
	dtrRepo repositories.DTRRepository,
	supportRepo repositories.SupportRepository,
	fbSvc *facebook.Service,
	stateManager *state.StateManager,
) *MenuHandler {
	return &MenuHandler{
		repo:           repo,
		linkRepo:       linkRepo,
		gradeRepo:      gradeRepo,
		bulletinRepo:   bulletinRepo,
		payableRepo:    payableRepo,
		paymentLogRepo: paymentLogRepo,
		dtrRepo:        dtrRepo,
		supportRepo:    supportRepo,
		fbSvc:          fbSvc,
		utils:          utils.NewResponseUtils(repo, linkRepo, fbSvc),
		stateManager:   stateManager,
	}
}

func (h *MenuHandler) ShowMainMenu(senderID string) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.IsActive == false {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	currentProfileData, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil || currentProfileData == nil {
		if err := h.stateManager.SetState(senderID, state.StateInitial, nil); err != nil {
			log.Printf("Error resetting state: %v", err)
		}
		return h.utils.SendResponseWithQuickReplies(senderID, "No active profile found. Please select a View Profiles and choose a profile to continue.")
	}

	if currentProfileData.Student == nil {
		return fmt.Errorf("student data not found for the primary profile")
	}

	fullName := fmt.Sprintf("%s %s", currentProfileData.Student.FirstName, currentProfileData.Student.LastName)
	menuMessage := fmt.Sprintf(constants.MainMenuTemplate,
		fullName,
		currentProfileData.Student.School.SchoolName,
	)
	quickreplies := helpers.GetMainMenuReplies()
	return h.fbSvc.SendQuickReplies(senderID, menuMessage, quickreplies)
}
