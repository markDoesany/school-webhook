package menu

import (
	"fmt"
	"log"
	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/state"
	"strings"
	"time"

	"school-assistant-wh/internal/models"
	"school-assistant-wh/internal/services/facebook"
	"school-assistant-wh/internal/services/helpers"
)

const (
	dtrPerPage = 10
)

// HandleViewDTR handles viewing DTR records with pagination
func (h *MenuHandler) HandleViewDTR(senderID string, pageNum int) error {
	if pageNum < 1 {
		pageNum = 1
	}

	user, err := h.repo.GetUserByPSID(senderID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsActive {
		return h.utils.SendResponseWithQuickReplies(senderID, constants.AccountDeactivatedMessage)
	}

	profile, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil {
		return fmt.Errorf("failed to get student profile: %w", err)
	}

	totalCount, err := h.dtrRepo.GetDTRRecordsCount(
		time.Now().Year(),
		time.Now().Month(),
		profile.Student.School.SchoolID,
		profile.Student.StudentID,
	)
	if err != nil {
		return fmt.Errorf("failed to get DTR records count: %w", err)
	}

	totalPages := int(totalCount+dtrPerPage-1) / dtrPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	if pageNum > totalPages {
		pageNum = totalPages
	}

	offset := (pageNum - 1) * dtrPerPage
	if offset < 0 {
		offset = 0
	}
	records, err := h.dtrRepo.GetDTRRecords(
		time.Now().Year(),
		time.Now().Month(),
		profile.Student.School.SchoolID,
		profile.Student.StudentID,
		offset,
		dtrPerPage,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch DTR records: %w", err)
	}

	if len(records) == 0 {
		if pageNum == 1 {
			quickReplies := helpers.GetBack()
			return h.fbSvc.SendQuickReplies(senderID, "No attendance records found.", quickReplies)
		} else {
			quickReplies := helpers.GetBack()
			return h.fbSvc.SendQuickReplies(senderID, "No more records to show.", quickReplies)
		}
	}

	recordsByDate := make(map[string][]models.DTRRecord)
	for _, record := range records {
		dateKey := record.DateTimeIN.Format("2006-01-02")
		recordsByDate[dateKey] = append(recordsByDate[dateKey], record)
	}

	var sb strings.Builder
	sb.WriteString("📋 *Your Attendance Records*\n\n")
	sb.WriteString(fmt.Sprintf("📅 *%s %d*\n\n", time.Now().Month().String(), time.Now().Year()))

	for date, dateRecords := range recordsByDate {
		d, _ := time.Parse("2006-01-02", date)
		sb.WriteString(fmt.Sprintf("📅 *%s, %s*\n", d.Format("Monday"), d.Format("Jan 02, 2006")))

		for _, record := range dateRecords {
			timeStr := record.DateTimeIN.Format("15:04:05")
			sb.WriteString(fmt.Sprintf("  %s %s\n", timeStr, record.Type))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("_Page %d of %d\n", pageNum, totalPages))

	// Update state with comprehensive pagination details
	if err := h.stateManager.SetState(senderID, state.StateViewDTR, map[string]interface{}{
		state.KeyPageOffset:      offset,
		state.KeyPaginationItems: len(records),
		state.KeyPaginationTotal: int(totalCount),
		state.KeyPaginationPage:  pageNum,
		state.KeyPaginationSize:  dtrPerPage,
		state.KeyPaginationPages: totalPages,
	}); err != nil {
		log.Printf("Error updating state: %v", err)
	}

	var quickReplies []facebook.QuickReply
	hasMore := int64(offset+dtrPerPage) < totalCount
	if hasMore {
		quickReplies = helpers.GetViewMoreReplies()
	} else {
		quickReplies = helpers.GetBack()
	}

	return h.fbSvc.SendQuickReplies(senderID, sb.String(), quickReplies)
}
