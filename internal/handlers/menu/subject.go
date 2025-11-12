package menu

import (
	"fmt"
	"log"
	"school-assistant-wh/internal/services/helpers"
	"school-assistant-wh/internal/state"
	"sort"
	"strings"
)

// HandleViewSubjects shows the list of enrolled subjects for the current term
func (h *MenuHandler) HandleViewSubjects(senderID string) error {
	// Get user profile
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	profile, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil {
		return fmt.Errorf("failed to get user links: %w", err)
	}

	// Get all enrolled subjects to determine available school years
	subjects, err := h.gradeRepo.GetEnrolledSubjects(profile.StudentID, profile.Student.School.SchoolID)
	if err != nil {
		return fmt.Errorf("failed to get enrolled subjects: %w", err)
	}

	// Group subjects by school year
	yearMap := make(map[string]bool)
	yearOptions := make(map[string]string)
	var years []string

	for _, subj := range subjects {
		if !yearMap[subj.SchoolYear] {
			yearMap[subj.SchoolYear] = true
			years = append(years, subj.SchoolYear)
		}
	}

	// Sort years in descending order (newest first)
	sort.Slice(years, func(i, j int) bool {
		return years[i] > years[j]
	})

	// Build the message
	var message string
	if len(years) == 0 {
		message = "*No enrolled subjects found.*"
	} else {
		var sb strings.Builder
		sb.WriteString("*View Subjects by School Year*\n\n")
		sb.WriteString("Please select a school year to view subjects:\n\n")

		for i, year := range years {
			sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, year))
			yearOptions[fmt.Sprint(i+1)] = year
		}

		message = sb.String()

	}
	stateData := map[string]any{
		state.KeySchoolYearMap: yearOptions,
	}

	if err := h.stateManager.SetState(senderID, state.StateSelectSubject, stateData); err != nil {
		log.Printf("Error setting view subjects state: %v", err)
	}

	quickReplies := helpers.GetBack()
	return h.fbSvc.SendQuickReplies(senderID, message, quickReplies)
}

// HandleViewSubjectsByYear displays subjects for a specific school year
func (h *MenuHandler) HandleViewSubjectsByYear(senderID, year string) error {
	// Get user profile
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	profile, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil {
		return fmt.Errorf("failed to get user links: %w", err)
	}

	// Get subjects for the selected school year
	subjects, err := h.gradeRepo.GetEnrolledSubjectsByYear(profile.StudentID, profile.Student.School.SchoolID, year)
	if err != nil {
		return fmt.Errorf("failed to get enrolled subjects: %w", err)
	}

	// Build the message
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📚 *Subjects for School Year %s*\n\n"+
		"*Name:* %s %s\n"+
		"*School:* %s\n"+
		"*Course:* %s\n\n",
		year,
		profile.Student.FirstName,
		profile.Student.LastName,
		profile.Student.School.SchoolName,
		profile.Student.Course,
	))

	if len(subjects) == 0 {
		sb.WriteString("No subjects found for this school year.")
	} else {
		for _, subj := range subjects {
			sb.WriteString(fmt.Sprintf("• *%s*\n"+
				"   Subject Unit: %s\n"+
				"   Schedule: %s\n"+
				"   Room: %s\n\n",
				subj.SubjectDescription,
				subj.SubjectUnit,
				subj.SubjectSchedule,
				subj.SubjectRoom,
			))
		}
	}

	// Update state to go back to subjects view
	if err := h.stateManager.SetState(senderID, state.StateViewSubjects, nil); err != nil {
		log.Printf("Error setting state: %v", err)
	}

	return h.fbSvc.SendQuickReplies(senderID, sb.String(), helpers.GetBack())
}
