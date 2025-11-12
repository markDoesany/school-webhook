package menu

import (
	"fmt"
	"log"
	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/models"
	"school-assistant-wh/internal/services/helpers"
	"school-assistant-wh/internal/state"
	"sort"
	"strings"
	"time"
)

// HandleViewPayables handles the request to view student payables
func (h *MenuHandler) HandleViewPayables(senderID string) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	currentProfileData, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil {
		return h.utils.SendResponseWithQuickReplies(senderID, "No active profile found. Please select a View Profiles and choose a profile to continue.")
	}

	if currentProfileData.Student == nil {
		return fmt.Errorf("student data not found for the primary profile")
	}

	err = h.stateManager.SetState(senderID, state.StateViewPayables, nil)
	if err != nil {
		return fmt.Errorf("failed to set state: %w", err)
	}

	// Get active payables for the student
	payables, err := h.payableRepo.GetActiveStudentPayables(
		currentProfileData.Student.School.SchoolID,
		currentProfileData.Student.StudentID,
	)
	if err != nil {
		log.Printf("Error fetching payables: %v", err)
		return h.utils.SendResponseWithQuickReplies(senderID, "Failed to fetch payables. Please try again later.")
	}

	if len(payables) == 0 {
		return h.fbSvc.SendTextMessage(senderID, "You don't have any active payables at the moment.")
	}

	// Group payables by school year and semester
	payablesByTerm := make(map[string]map[string][]models.StudentPayable)
	for _, p := range payables {
		// Use "Not Specified" as default for nil SchoolYear
		schoolYear := "Not Specified"
		if p.SchoolYear != nil && *p.SchoolYear != "" {
			schoolYear = *p.SchoolYear
		}

		// Use "Not Specified" as default for empty Semester
		semester := p.Semester
		if semester == "" || semester == "." {
			semester = "Not Specified"
		}

		// Initialize the maps if they don't exist
		if payablesByTerm[schoolYear] == nil {
			payablesByTerm[schoolYear] = make(map[string][]models.StudentPayable)
		}

		payablesByTerm[schoolYear][semester] = append(payablesByTerm[schoolYear][semester], p)
	}

	// Sort school years in descending order (newest first)
	var schoolYears []string
	for sy := range payablesByTerm {
		schoolYears = append(schoolYears, sy)
	}
	sort.Slice(schoolYears, func(i, j int) bool {
		// Put "Not Specified" at the end
		if schoolYears[i] == "Not Specified" {
			return false
		}
		if schoolYears[j] == "Not Specified" {
			return true
		}
		return schoolYears[i] > schoolYears[j] // Sort in descending order
	})

	// Calculate total balance
	totalBalance := 0.0
	for _, p := range payables {
		totalBalance += p.TotalAmountToPay
	}

	// Count total terms (school year + semester combinations)
	totalTerms := 0
	for _, semesters := range payablesByTerm {
		totalTerms += len(semesters)
	}

	// Send summary message
	summaryMsg := fmt.Sprintf(
		"📋 *Your Active Payables*\n\n"+
			"👤 *%s %s*\n"+
			"📝 Student ID: %s\n"+
			"🏫 %s\n\n"+
			"*Total Balance: ₱%.2f*\n"+
			"*%d payment(s) across %d term(s)*\n\n"+
			"Here are your payables, grouped by school term:",
		currentProfileData.Student.FirstName,
		currentProfileData.Student.LastName,
		currentProfileData.Student.StudentID,
		currentProfileData.Student.School.SchoolName,
		totalBalance,
		len(payables),
		totalTerms,
	)

	if err := h.fbSvc.SendTextMessage(senderID, summaryMsg); err != nil {
		log.Printf("Error sending summary message: %v", err)
	}

	const maxPayablesPerMessage = 5
	var currentMessage strings.Builder
	var currentCount int

	// Function to send the current batch of payables
	sendCurrentBatch := func() error {
		if currentMessage.Len() == 0 {
			return nil
		}

		// Add progress indicator
		progress := fmt.Sprintf("\n\n_Showing %d of %d payables_", currentCount, len(payables))
		currentMessage.WriteString(progress)

		if err := h.fbSvc.SendTextMessage(senderID, currentMessage.String()); err != nil {
			return fmt.Errorf("error sending payables batch: %w", err)
		}

		time.Sleep(300 * time.Millisecond) // Small delay between messages
		currentMessage.Reset()
		return nil
	}

	// Process payables by school year and semester
	for _, schoolYear := range schoolYears {
		semesters := payablesByTerm[schoolYear]

		// Sort semesters in a logical order
		var sortedSemesters []string
		for sem := range semesters {
			sortedSemesters = append(sortedSemesters, sem)
		}
		sort.Slice(sortedSemesters, func(i, j int) bool {
			// Custom sort order: 1st Semester, 2nd Semester, Summer, Not Specified, others
			sem1, sem2 := sortedSemesters[i], sortedSemesters[j]
			if sem1 == sem2 {
				return false
			}
			order := map[string]int{"1st Semester": 1, "2nd Semester": 2, "Summer": 3, "Not Specified": 99}
			i1, i2 := order[sem1], order[sem2]
			if i1 == 0 {
				i1 = 4
			}
			if i2 == 0 {
				i2 = 4
			}
			return i1 < i2
		})

		for _, semester := range sortedSemesters {
			payableGroup := semesters[semester]

			// Add school year and semester header
			header := fmt.Sprintf("\n📚 *%s* • %s\n", schoolYear, semester)
			if currentMessage.Len()+len(header) > 1800 {
				if err := sendCurrentBatch(); err != nil {
					log.Printf("Error sending batch: %v", err)
				}
			}
			currentMessage.WriteString(header)

			// Add payables for this term
			for _, p := range payableGroup {
				// Format payable details
				details := fmt.Sprintf(
					"➤ *%s*\n"+
						"   SOA ID: %s\n"+
						"   Amount: ₱%.2f\n"+
						"   Type: %s\n\n",
					p.Particulars,
					p.SOAID,
					p.TotalAmountToPay,
					func() string {
						if p.Type != nil && *p.Type != "" {
							return *p.Type
						}
						return "Not specified"
					}(),
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
				if currentCount%maxPayablesPerMessage == 0 {
					if err := sendCurrentBatch(); err != nil {
						log.Printf("Error sending batch: %v", err)
					}
				}
			}
		}
	}

	// Send any remaining payables
	if currentMessage.Len() > 0 {
		if err := sendCurrentBatch(); err != nil {
			log.Printf("Error sending final batch: %v", err)
		}
	}

	// Send final message with navigation options
	quickReplies := helpers.GetPaymentReplies()
	return h.fbSvc.SendQuickReplies(senderID, "What would you like to do next?", quickReplies)
}
