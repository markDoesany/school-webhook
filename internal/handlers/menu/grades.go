package menu

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/models"
	"school-assistant-wh/internal/services/helpers"
	"school-assistant-wh/internal/state"
)

// groupGradesBySchoolYear groups grades by their school year
func groupGradesBySchoolYear(grades []models.SubjectGrade) map[string][]models.SubjectGrade {
	yearMap := make(map[string][]models.SubjectGrade)

	for _, grade := range grades {
		yearMap[grade.SchoolYear] = append(yearMap[grade.SchoolYear], grade)
	}

	return yearMap
}

// HandleViewGrades handles the view grades menu option
func (h *MenuHandler) HandleViewGrades(senderID string) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.IsActive == false {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	currentProfileData, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil {
		return fmt.Errorf("failed to get user links: %w", err)
	}

	grades, err := h.gradeRepo.GetStudentSubjectGrades(
		currentProfileData.Student.StudentID,
		currentProfileData.Student.School.SchoolID,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch grades: %w", err)
	}

	yearMap := groupGradesBySchoolYear(grades)

	if len(yearMap) == 0 {
		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(senderID, "No grades found for this student.", quickReplies)
	}

	var years []string
	yearOptions := make(map[string]string)
	i := 1
	for year := range yearMap {
		years = append(years, year)
		yearOptions[fmt.Sprint(i)] = year
		i++
	}

	sort.Slice(years, func(i, j int) bool {
		return years[i] > years[j]
	})

	var sb strings.Builder
	sb.WriteString("📚 *View Grades by School Year*\n\n")
	sb.WriteString("Please select a school year to view grades:\n\n")

	for i, year := range years {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, year))
	}

	stateData := map[string]any{
		state.KeySchoolYearMap: yearOptions,
	}

	if err := h.stateManager.SetState(senderID, state.StateViewGrades, stateData); err != nil {
		log.Printf("Error setting view grades state: %v", err)
	}

	return h.fbSvc.SendQuickReplies(senderID, sb.String(), helpers.GetBack())
}

// HandleViewGradesByYear displays grades for a specific school year, grouped by semester
func (h *MenuHandler) HandleViewGradesByYear(senderID, year string) error {
	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.IsActive == false {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	currentProfileData, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil {
		return fmt.Errorf("failed to get user links: %w", err)
	}

	grades, err := h.gradeRepo.GetStudentSubjectGrades(
		currentProfileData.Student.StudentID,
		currentProfileData.Student.School.SchoolID,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch grades: %w", err)
	}

	// Filter and group grades by semester
	semesterGrades := make(map[string][]models.SubjectGrade)
	for _, grade := range grades {
		if grade.SchoolYear == year {
			semesterGrades[grade.Semester] = append(semesterGrades[grade.Semester], grade)
		}
	}

	if len(semesterGrades) == 0 {
		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(senderID,
			fmt.Sprintf("No grades found for school year %s.", year),
			quickReplies)
	}

	// Sort semesters
	var semesters []string
	for sem := range semesterGrades {
		semesters = append(semesters, sem)
	}
	sort.Strings(semesters)

	// Build messages for each semester
	var allMessages []string
	for _, sem := range semesters {
		semesterGrades := semesterGrades[sem]
		sort.Slice(semesterGrades, func(i, j int) bool {
			return semesterGrades[i].SubjectDescription < semesterGrades[j].SubjectDescription
		})

		var gradesList strings.Builder
		gradesList.WriteString(fmt.Sprintf("📚 *%s - %s*\n\n", sem, year))

		for _, grade := range semesterGrades {
			gradesList.WriteString(fmt.Sprintf(constants.GradeItemTemplate,
				grade.SubjectDescription,
				grade.StudentGrade,
				grade.ExamTerm,
			))
		}
		allMessages = append(allMessages, gradesList.String())
	}

	// Send messages
	for _, msg := range allMessages {
		if err := h.fbSvc.SendTextMessage(senderID, msg); err != nil {
			log.Printf("Error sending grade message: %v", err)
		}
	}

	quickReplies := helpers.GetBack()
	return h.fbSvc.SendQuickReplies(senderID, "Select another option:", quickReplies)
}
