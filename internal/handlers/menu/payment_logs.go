package menu

import (
	"fmt"
	"log"
	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/services/helpers"
	"sort"
	"strings"
	"time"
)

// HandleViewPaymentLogs handles the request to view payment logs for a student
func (h *MenuHandler) HandleViewPaymentLogs(senderID string) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	currentProfileData, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil {
		return h.utils.SendResponseWithQuickReplies(senderID, "No active profile found. Please select View Profiles and choose a profile to continue.")
	}

	if currentProfileData.Student == nil {
		return fmt.Errorf("student data not found for the primary profile")
	}

	// Get current year for default log retrieval
	currentYear := time.Now().Year()

	// Get payment logs for the student
	logs, err := h.paymentLogRepo.GetPaymentLogsByStudentAndSchoolID(currentYear, currentProfileData.Student.StudentID, currentProfileData.Student.School.SchoolID)
	if err != nil {
		log.Printf("Error fetching payment logs: %v", err)
		return h.utils.SendResponseWithQuickReplies(senderID, "Failed to fetch payment history. Please try again later.")
	}

	if len(logs) == 0 {
		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(senderID, "You don't have any payment history at the moment.", quickReplies)
	}

	// Sort logs by payment date (newest first)
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].DateTimePaid.After(logs[j].DateTimePaid)
	})

	// Calculate total payments
	totalPaid := 0.0
	for _, log := range logs {
		totalPaid += log.Amount
	}

	// Send summary message
	summaryMsg := fmt.Sprintf(
		"💳 *Your Payment History*\n\n"+
			"👤 *%s %s*\n"+
			"📝 Student ID: %s\n"+
			"🏫 %s\n\n"+
			"*Total Payments: ₱%.2f*\n"+
			"*%d transaction(s) found*\n\n"+
			"Here are your payment records:",
		currentProfileData.Student.FirstName,
		currentProfileData.Student.LastName,
		currentProfileData.Student.StudentID,
		currentProfileData.Student.School.SchoolName,
		totalPaid,
		len(logs),
	)

	if err := h.fbSvc.SendTextMessage(senderID, summaryMsg); err != nil {
		log.Printf("Error sending summary message: %v", err)
	}

	const maxLogsPerMessage = 5
	var currentMessage strings.Builder
	var currentCount int

	// Function to send the current batch of logs
	sendCurrentBatch := func() error {
		if currentMessage.Len() == 0 {
			return nil
		}

		// Add progress indicator
		progress := fmt.Sprintf("\n\n_Showing %d of %d transactions_", currentCount, len(logs))
		currentMessage.WriteString(progress)

		if err := h.fbSvc.SendTextMessage(senderID, currentMessage.String()); err != nil {
			return fmt.Errorf("error sending payment logs batch: %w", err)
		}

		time.Sleep(300 * time.Millisecond) // Small delay between messages
		currentMessage.Reset()
		return nil
	}

	// Process logs
	for _, paymentLog := range logs {
		// Format payment log details
		details := fmt.Sprintf(
			"\n📅 *%s*\n"+
				"   Transaction ID: %s\n"+
				"   Amount: ₱%.2f\n"+
				"   Status: %s\n"+
				"   Reference: %s\n"+
				"   Payment Type: %s\n"+
				"   SOA ID: %s\n",
			paymentLog.DateTimePaid.Format("January 2, 2006 3:04 PM"),
			paymentLog.PaymentTxnID,
			paymentLog.Amount,
			paymentLog.Status,
			paymentLog.ProcessID,
			paymentLog.PaymentType,
			paymentLog.SOAID,
		)

		// Check if we need to send the current batch
		if currentMessage.Len()+len(details) > 1800 {
			if err := sendCurrentBatch(); err != nil {
				log.Printf("Error sending batch: %v", err)
			}
		}

		currentMessage.WriteString(details)
		currentCount++

		// Check if we've reached the max items per message
		if currentCount%maxLogsPerMessage == 0 {
			if err := sendCurrentBatch(); err != nil {
				log.Printf("Error sending batch: %v", err)
			}
		}
	}

	// Send any remaining logs
	if currentMessage.Len() > 0 {
		if err := sendCurrentBatch(); err != nil {
			log.Printf("Error sending final batch: %v", err)
		}
	}

	// Send final message with navigation options
	quickReplies := helpers.GetBack()
	return h.fbSvc.SendQuickReplies(senderID, "What would you like to do next?", quickReplies)
}
