package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"school-assistant-wh/internal/config"
	"school-assistant-wh/internal/handlers/account"
	"school-assistant-wh/internal/handlers/menu"
	"school-assistant-wh/internal/handlers/utils"
	"school-assistant-wh/internal/repositories"
	"school-assistant-wh/internal/services/facebook"
	"school-assistant-wh/internal/services/helpers"
	"school-assistant-wh/internal/state"
)

type Handler struct {
	repo           repositories.UserRepository
	linkRepo       repositories.UserLinkRepository
	gradeRepo      repositories.GradeRepository
	bulletinRepo   *repositories.BulletinRepository
	payableRepo    *repositories.StudentPayableRepository
	paymentLogRepo *repositories.PaymentLogRepository
	dtrRepo        *repositories.DTRRepository
	supportRepo    *repositories.SupportRepository
	fbSvc          *facebook.Service
	accountHdlr    *account.AccountHandler
	menuHdlr       *menu.MenuHandler
	utils          *utils.ResponseUtils
	stateManager   *state.StateManager
}

func NewHandler(db *gorm.DB, fbSvc *facebook.Service) *Handler {
	repo := repositories.NewUserRepository(db, fbSvc)
	linkRepo := repositories.NewUserLinkRepository(db, repositories.NewStudentProfileRepository(db))
	gradeRepo := repositories.NewGradeRepository(db)
	bulletinRepo := repositories.NewBulletinRepository(db)
	payableRepo := repositories.NewStudentPayableRepository(db)
	paymentLogRepo := repositories.NewPaymentLogRepository(db)
	dtrRepo := repositories.NewDTRRepository(db)
	supportRepo := repositories.NewSupportRepository(db)
	stateManager := state.NewStateManager()

	// Create account handler with state manager
	accountHdlr := account.NewAccountHandler(*repo, *linkRepo, fbSvc, stateManager)
	menuHdlr := menu.NewMenuHandler(*repo, *linkRepo, *gradeRepo, *bulletinRepo, *payableRepo, *paymentLogRepo, *dtrRepo, *supportRepo, fbSvc, stateManager)

	// Preload active users into cache
	if err := repo.PreloadActiveUsers(); err != nil {
		log.Printf("Failed to preload active users: %v", err)
	}

	// Start state cleanup goroutine
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			stateManager.CleanupInactive(24 * time.Hour)
		}
	}()

	return &Handler{
		repo:         *repo,
		linkRepo:     *linkRepo,
		fbSvc:        fbSvc,
		accountHdlr:  accountHdlr,
		menuHdlr:     menuHdlr,
		utils:        utils.NewResponseUtils(*repo, *linkRepo, fbSvc),
		stateManager: stateManager,
	}
}

func (h *Handler) VerifyWebhook(c *gin.Context) {
	verifyToken := c.Query("hub.verify_token")
	if verifyToken == config.LoadFacebookConfig().VerifyToken {
		challenge := c.Query("hub.challenge")
		c.String(http.StatusOK, challenge)
		return
	}
	c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verification token"})
}

func (h *Handler) SetupMessengerProfile() error {
	if err := h.fbSvc.SetGetStartedButton("GET_STARTED"); err != nil {
		return err
	}

	greeting := "Hello! I'm your school assistant. How can I help yoqu today?"
	if err := h.fbSvc.SetGreetingText(greeting); err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleWebhook(c *gin.Context) {
	var request facebook.WebhookRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if request.Object != "page" {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	for _, entry := range request.Entry {
		for _, messaging := range entry.Messaging {
			senderID := messaging.Sender.ID

			if messaging.Postback != nil && messaging.Postback.Payload != "" {
				if err := h.handleMessagePayload(senderID, messaging.Postback.Payload); err != nil {
					log.Printf("Error handling message payload: %v", err)
				}
			}

			if messaging.Message != nil && messaging.Message.Text != "" {
				message := strings.TrimSpace(strings.ToUpper(messaging.Message.Text))
				currentState, stateData := h.stateManager.GetState(senderID)
				if currentState != "" && currentState != state.StateInitial && !helpers.IsQuickReplyPayload(message) {
					if err := h.handleStateMessage(senderID, message, currentState, stateData); err != nil {
						log.Printf("Error handling state message: %v", err)
					} else {
						return
					}
				}

				if err := h.handleMessage(senderID, message); err != nil {
					log.Printf("Error handling message: %v", err)
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) handleMessagePayload(senderID string, payload string) error {
	switch strings.ToUpper(payload) {
	case "GET_STARTED":
		return h.handleGetStarted(senderID)
	}

	return nil
}

func (h *Handler) handleStateMessage(senderID, message string, currentState state.State, stateData map[string]any) error {
	switch currentState {
	case state.StateProfileSwitch:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateInitial, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.accountHdlr.HandleViewProfile(senderID)
		}
		return h.handleProfileSelection(senderID, message, stateData)

	case state.StateProfileView:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateInitial, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.handleGetStarted(senderID)
		}
		return h.handleProfileSelection(senderID, message, stateData)
	case state.StateMainMenu:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}

			return h.menuHdlr.ShowMainMenu(senderID)
		}
		return h.menuHdlr.MenuHandler(senderID, message)
	case state.StateViewGrades:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.ShowMainMenu(senderID)
		}
		return h.handleSchoolYearSelection(senderID, message, stateData)
	case state.StateViewGradesDetails:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateViewGrades, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.HandleViewGrades(senderID)
		}
		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(senderID,
			"Invalid selection. Go back to view grades by School Year.",
			quickReplies,
		)

	case state.StateViewBulletin:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.ShowMainMenu(senderID)
		}

		if message == "VIEW MORE" {
			if p, ok := stateData[state.KeyPaginationPage].(int); ok {
				nextPageNum := p + 1
				return h.menuHdlr.HandleViewBulletin(senderID, nextPageNum)
			}
		}
		quickReplies := helpers.GetViewMoreReplies()
		return h.fbSvc.SendQuickReplies(senderID,
			" Invalid selection. Go back to main menu or view more.",
			quickReplies,
		)

	case state.StateViewPayables:
		switch message {
		case "BACK":
			if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.ShowMainMenu(senderID)
		case "PAYMENT LOGS":
			return h.menuHdlr.HandleViewPaymentLogs(senderID)
		default:
			quickReplies := helpers.GetPaymentReplies()
			return h.fbSvc.SendQuickReplies(senderID,
				"Invalid selection. Go back to main menu or view payment logs.",
				quickReplies,
			)
		}
	case state.StateViewDTR:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.ShowMainMenu(senderID)
		}
		if message == "VIEW MORE" {
			if p, ok := stateData[state.KeyPaginationPage].(int); ok {
				nextPageNum := p + 1
				return h.menuHdlr.HandleViewDTR(senderID, nextPageNum)
			}
		}
		quickReplies := helpers.GetViewMoreReplies()
		return h.fbSvc.SendQuickReplies(senderID,
			" Invalid selection. Go back to main menu or view more.",
			quickReplies,
		)
	case state.StateProfileMenu:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.ShowMainMenu(senderID)
		}
		return h.menuHdlr.HandleProfileMenuSelection(senderID, message)
	case state.StateConfirmProfileSwitch:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateProfileMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.ShowProfileMenu(senderID)
		}
		if message == "PROCEED" {
			return h.accountHdlr.HandleSwitchProfile(senderID)
		}
		quickReplies := helpers.GetConfirmProfileSwitch()
		return h.fbSvc.SendQuickReplies(senderID,
			"Invalid selection. Go back to profile menu or proceed.",
			quickReplies,
		)
	case state.StateViewSubjects:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateProfileMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.ShowProfileMenu(senderID)
		}
		return h.handleSubjectByYearSelection(senderID, message, stateData)
	case state.StateSelectSubject:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateProfileMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.ShowProfileMenu(senderID)
		}
		return h.handleSubjectByYearSelection(senderID, message, stateData)
	case state.StateAskSupport:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.ShowMainMenu(senderID)
		}
		if message == "VIEW TICKETS" {
			return h.menuHdlr.ListSupportTickets(senderID)
		}
		return h.menuHdlr.AddSupportMessage(senderID, message, "")
	case state.StateViewTickets:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.ShowMainMenu(senderID)
		}
		return h.menuHdlr.ListSupportTickets(senderID)
	case state.StateSelectSupportTicket:
		if message == "BACK" {
			if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
				log.Printf("Error resetting state: %v", err)
			}
			return h.menuHdlr.ShowMainMenu(senderID)
		}
		return h.handleSupportTicketSelection(senderID, message, stateData)
	default:
		return h.handleDefault(senderID)
	}
}

func (h *Handler) handleMessage(senderID string, message string) error {
	switch {
	case message == "REGISTER":
		return h.accountHdlr.HandleRegistration(senderID)
	case message == "MENU":
		if err := h.stateManager.SetState(senderID, state.StateMainMenu, nil); err != nil {
			log.Printf("Error resetting state: %v", err)
		}
		return h.menuHdlr.ShowMainMenu(senderID)
	case message == "MY SA-ID":
		return h.accountHdlr.HandleViewSaID(senderID)
	case message == "VIEW PROFILE":
		return h.accountHdlr.HandleViewProfile(senderID)
	case message == "SWITCH PROFILE":
		return h.accountHdlr.HandleSwitchProfile(senderID)
	case message == "CONTINUE":
		return h.handleContinue(senderID)
	case message == "NO":
		return h.handleNo(senderID)
	case message == "ABOUT US":
		return h.handleAboutUs(senderID)
	case message == "TALK TO HUMAN":
		return h.handleTalkToHuman(senderID)
	default:
		exists, err := h.repo.UserExists(senderID)
		if err != nil {
			return fmt.Errorf("failed to check user existence: %w", err)
		}
		if !exists {
			return h.handleGetStarted(senderID)
		}

		return h.handleDefault(senderID)
	}
}
