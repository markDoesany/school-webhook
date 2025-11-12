package menu

import (
	"fmt"
	"log"
	"strings"

	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/models"
	"school-assistant-wh/internal/services/helpers"
	"school-assistant-wh/internal/state"
	"school-assistant-wh/internal/utils"
)

// AskSupport handles the support request flow
func (h *MenuHandler) AskSupport(senderID string) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return h.utils.SendResponseWithQuickReplies(senderID, "An error occurred. Please try again later.")
	}

	if !user.IsActive {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	currentProfileData, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil || currentProfileData == nil {
		return h.utils.SendResponseWithQuickReplies(senderID, "No active profile found. Please select View Profiles and choose a profile to continue.")
	}

	// Set state to expect a support message
	h.stateManager.SetState(senderID, state.StateAskSupport, nil)

	// Ask user to type their support message
	message := "You may now type your inquiry message below or view your support tickets."
	quickReplies := helpers.GetAskSupportReplies()
	return h.fbSvc.SendQuickReplies(senderID, message, quickReplies)
}

// AddSupportMessage adds a new message to an existing support thread
func (h *MenuHandler) AddSupportMessage(senderID, message, threadID string) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return h.utils.SendResponseWithQuickReplies(senderID, "An error occurred. Please try again later.")
	}

	if !user.IsActive {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	currentProfileData, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil || currentProfileData == nil {
		return h.utils.SendResponseWithQuickReplies(senderID, "No active profile found. Please select a profile to continue.")
	}

	if threadID == "" {
		subject := "From Messenger App"
		thread := &models.SupportThread{
			MobileNo:       currentProfileData.Student.MobileNumber,
			GKBorrowerID:   currentProfileData.Student.BorrowerID,
			GKBorrowerName: currentProfileData.Student.FirstName + " " + currentProfileData.Student.LastName,
			HelpTopic:      "General Inquiry",
			Subject:        &subject,
			Status:         utils.StringPtr("OPEN"),
		}

		var err error
		threadID, err = h.supportRepo.CreateThread(thread, currentProfileData.Student.School.SchoolID)
		if err != nil {
			log.Printf("Error creating support thread: %v", err)
			return h.utils.SendResponseWithQuickReplies(senderID, "Failed to create support ticket. Please try again later.")
		}

		h.stateManager.SetState(senderID, state.StateAskSupport, map[string]any{
			state.KeyThreadID: threadID,
		})

		// Send confirmation that a new ticket was created
		confirmationMsg := fmt.Sprintf("✅ Support ticket #%s has been created.\n\n", threadID)
		err = h.utils.SendResponseWithQuickReplies(senderID, confirmationMsg)
		if err != nil {
			log.Printf("Error sending confirmation message: %v", err)
		}
	}

	supportMessage := &models.SupportConversation{
		ThreadID:           threadID,
		ReplySupportUserID: currentProfileData.Student.StudentID,
		ReplySupportName:   currentProfileData.Student.FirstName + " " + currentProfileData.Student.LastName,
		ThreadType:         "0",
		Message:            message,
	}

	err = h.supportRepo.CreateMessage(supportMessage, currentProfileData.Student.School.SchoolID)
	if err != nil {
		log.Printf("Error saving support message: %v", err)
		return h.utils.SendResponseWithQuickReplies(senderID, "Failed to send your message. Please try again.")
	}

	h.stateManager.SetState(senderID, state.StateInitial, nil)
	return h.utils.SendResponseWithQuickReplies(senderID, "✅ Your message has been sent. Our support team will respond as soon as possible.")
}

// ListSupportTickets displays the user's support tickets and allows selection
func (h *MenuHandler) ListSupportTickets(senderID string) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return h.utils.SendResponseWithQuickReplies(senderID, "An error occurred. Please try again later.")
	}

	if !user.IsActive {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	currentProfileData, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil || currentProfileData == nil {
		return h.utils.SendResponseWithQuickReplies(senderID, "No active profile found. Please select a profile to continue.")
	}

	tickets, err := h.supportRepo.GetThreadsByBorrowerID(currentProfileData.Student.BorrowerID, currentProfileData.Student.School.SchoolID)
	if err != nil {
		log.Printf("Error fetching support tickets: %v", err)
		return h.utils.SendResponseWithQuickReplies(senderID, "Failed to fetch your support tickets. Please try again later.")
	}

	if len(tickets) == 0 {
		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(senderID, "You don't have any support tickets yet.", quickReplies)
	}

	var message strings.Builder
	message.WriteString("📋 Your Support Tickets:\n\n")

	ticketMap := make(map[string]string)

	for i, ticket := range tickets {
		status := ""
		if ticket.Status != nil && *ticket.Status == "CLOSED" {
			status = " (Closed)"
		}
		ticketNumber := fmt.Sprint(i + 1)
		ticketMap[ticketNumber] = ticket.ThreadID
		message.WriteString(fmt.Sprintf("[%s] #%s%s\n", ticketNumber, ticket.ThreadID, status))

		if i >= 9 {
			break
		}
	}

	message.WriteString("\nPlease reply with the number of the ticket you want to view.")

	h.stateManager.SetState(senderID, state.StateSelectSupportTicket, map[string]any{
		state.KeyTicketMap: ticketMap,
	})

	quickReplies := helpers.GetBack()
	return h.fbSvc.SendQuickReplies(senderID, message.String(), quickReplies)
}

// HandleSupportTicketSelection processes the user's ticket selection
func (h *MenuHandler) HandleSupportTicketSelection(senderID, ticketID string) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return h.utils.SendResponseWithQuickReplies(senderID, "An error occurred. Please try again later.")
	}

	if !user.IsActive {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	currentProfileData, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil || currentProfileData == nil {
		return h.utils.SendResponseWithQuickReplies(senderID, "No active profile found. Please select a profile to continue.")
	}

	// Get the ticket details
	ticket, err := h.supportRepo.GetThread(ticketID, currentProfileData.Student.School.SchoolID)
	if err != nil {
		log.Printf("Error fetching ticket %s: %v", ticketID, err)
		return h.utils.SendResponseWithQuickReplies(senderID, "Could not find the selected ticket. Please try again.")
	}

	// Get ticket messages
	messages, err := h.supportRepo.GetMessages(ticketID, currentProfileData.Student.School.SchoolID)
	if err != nil {
		log.Printf("Error fetching messages for ticket %s: %v", ticketID, err)
		return h.utils.SendResponseWithQuickReplies(senderID, "Failed to load ticket messages. Please try again.")
	}

	var message strings.Builder
	status := "OPEN"
	if ticket.Status != nil && *ticket.Status == "CLOSED" {
		status = "CLOSED"
	}

	message.WriteString(fmt.Sprintf("📋 Ticket #%s (%s)\n\n", ticket.ThreadID, status))

	if ticket.HelpTopic != "" {
		message.WriteString(fmt.Sprintf("Topic: %s\n", ticket.HelpTopic))
	}

	message.WriteString("\n💬 Conversation:\n\n")

	// Add messages
	for _, msg := range messages {
		timestamp := ""
		if !msg.DateTimeIN.IsZero() {
			timestamp = msg.DateTimeIN.Format("Jan 2, 3:04 PM") + "\n"
		}

		if msg.ThreadType == "0" {
			message.WriteString(fmt.Sprintf("%s\n", timestamp, "You:"))
			message.WriteString(fmt.Sprintf("\n", msg.Message))
		} else {
			message.WriteString(fmt.Sprintf("%s\n", timestamp, "Support:"))
			message.WriteString(fmt.Sprintf("%s\n", msg.Message))
		}
	}

	if ticket.Status != nil && *ticket.Status == "CLOSED" {
		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(senderID, message.String()+"\n\nThis ticket is closed. You can view other tickets or go back.", quickReplies)
	}

	err = h.fbSvc.SendTextMessage(senderID, message.String()+"\n\nYou can select another ticket to view its details or go back to the main menu.")
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return h.utils.SendResponseWithQuickReplies(senderID, "Failed to load the conversation. Please try again later.")
	}

	quickReplies := helpers.GetBack()
	return h.fbSvc.SendQuickReplies(senderID, "What would you like to do next?", quickReplies)
}
