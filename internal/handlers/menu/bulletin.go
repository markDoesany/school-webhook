package menu

import (
	"fmt"
	"log"
	"strings"
	"time"

	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/services/facebook"
	"school-assistant-wh/internal/services/helpers"
	"school-assistant-wh/internal/state"
)

const (
	bulletinsPerPage = 3
)

// HandleViewBulletin handles viewing bulletins with pagination
func (h *MenuHandler) HandleViewBulletin(senderID string, pageNum int) error {
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

	currentProfileData, err := h.linkRepo.GetPrimaryLink(int(user.ID))
	if err != nil {
		return fmt.Errorf("failed to get user links: %w", err)
	}

	schoolID := currentProfileData.Student.School.SchoolID

	// Get total count of bulletins
	totalCount, err := h.bulletinRepo.GetBulletinsCount(schoolID)
	if err != nil {
		return fmt.Errorf("failed to get bulletins count: %w", err)
	}

	// Calculate offset from page number
	offset := (pageNum - 1) * bulletinsPerPage
	if offset < 0 {
		offset = 0
	}

	// Get paginated bulletins
	bulletins, err := h.bulletinRepo.GetBulletins(schoolID, nil, offset, bulletinsPerPage)
	if err != nil {
		return fmt.Errorf("failed to fetch bulletins: %w", err)
	}

	if len(bulletins) == 0 {
		if offset == 0 {
			quickReplies := helpers.GetBack()
			return h.fbSvc.SendQuickReplies(senderID, "No active bulletins found.", quickReplies)
		} else {
			quickReplies := helpers.GetBack()
			return h.fbSvc.SendQuickReplies(senderID, "No more bulletins to show.", quickReplies)
		}
	}

	// Send header message
	headerMessage := "📰 *Bulletin Board*\n\nStay informed with the latest school updates."
	if err := h.fbSvc.SendTextMessage(senderID, headerMessage); err != nil {
		return fmt.Errorf("failed to send header message: %w", err)
	}

	// Send each bulletin as a separate message
	for i, bulletin := range bulletins {
		var message strings.Builder

		// Add bulletin title with label
		message.WriteString(fmt.Sprintf("*Title:* %s\n", bulletin.Title))

		// Add date with label if available and not the Unix epoch
		isEpoch := bulletin.PeriodStart.Year() == 1970 &&
			bulletin.PeriodStart.Month() == 1 &&
			bulletin.PeriodStart.Day() == 1 &&
			!bulletin.PeriodStart.IsZero()

		if !bulletin.PeriodStart.IsZero() && !isEpoch {
			message.WriteString(fmt.Sprintf("*Date:* %s\n",
				bulletin.PeriodStart.Format("January 02, 2006")))
		}

		// Send image first if available
		if bulletin.ImageURL != "" && bulletin.ImageURL != "." {
			if err := h.fbSvc.SendImage(senderID, bulletin.ImageURL); err != nil {
				log.Printf("Failed to send image (URL: %s): %v", bulletin.ImageURL, err)
			}
		}

		// Add description with label if available
		if bulletin.Description != nil && *bulletin.Description != "" {
			desc := *bulletin.Description
			if len(desc) > 200 { // Truncate long descriptions
				desc = desc[:200] + "..."
			}
			message.WriteString(fmt.Sprintf("\n*Description:*\n%s\n", desc))
		}

		// Add link with label if available in Notes1
		if bulletin.Notes1 != nil && *bulletin.Notes1 != "" {
			notes := strings.TrimSpace(*bulletin.Notes1)
			startTag := "<redirectionlink>"
			endTag := "</redirectionlink>"
			startIdx := strings.Index(notes, startTag)
			endIdx := strings.Index(notes, endTag)

			if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
				link := strings.TrimSpace(notes[startIdx+len(startTag) : endIdx])
				if link != "" {
					message.WriteString(fmt.Sprintf("\n*Link:* %s\n", link))
				}
			}
		}

		// Add separator between bulletins
		if i < len(bulletins)-1 {
			message.WriteString("\n━━━━━━━━━━━━━━━━━━\n")
		}

		// Send the bulletin message
		if message.Len() > 0 {
			if err := h.fbSvc.SendTextMessage(senderID, message.String()); err != nil {
				log.Printf("Failed to send bulletin message: %v", err)
			}
		}

		time.Sleep(100 * time.Millisecond) // Small delay between messages
	}

	// Calculate pagination details
	currentPage := (offset / bulletinsPerPage) + 1
	totalPages := (int(totalCount) + bulletinsPerPage - 1) / bulletinsPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	// Update state with comprehensive pagination details
	if err := h.stateManager.SetState(senderID, state.StateViewBulletin, map[string]interface{}{
		state.KeyPageOffset:      offset,
		state.KeyPaginationItems: len(bulletins),
		state.KeyPaginationTotal: int(totalCount),
		state.KeyPaginationPage:  currentPage,
		state.KeyPaginationSize:  bulletinsPerPage,
		state.KeyPaginationPages: totalPages,
	}); err != nil {
		log.Printf("Error updating state: %v", err)
	}

	// Send pagination info to user
	messageText := fmt.Sprintf("_Page %d of %d_", currentPage, totalPages)

	// Prepare and send quick replies
	var quickReplies []facebook.QuickReply
	hasMore := int64(offset+bulletinsPerPage) < totalCount
	if hasMore {
		quickReplies = helpers.GetViewMoreReplies()
	} else {
		quickReplies = helpers.GetBack()
	}

	return h.fbSvc.SendQuickReplies(senderID, messageText, quickReplies)
}
