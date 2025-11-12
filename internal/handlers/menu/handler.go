package menu

import "school-assistant-wh/internal/services/helpers"

// MenuHandler processes the user's menu selection
func (h *MenuHandler) MenuHandler(senderID, selection string) error {
	switch selection {
	case "1":
		return h.HandleViewGrades(senderID)
	case "2":
		return h.HandleViewPayables(senderID)
	case "3":
		return h.HandleViewBulletin(senderID, 1)
	case "4":
		return h.HandleViewDTR(senderID, 1)
	case "5":
		return h.ShowProfileMenu(senderID)
	case "6":
		return h.AskSupport(senderID)
	default:
		quickReplies := helpers.GetBack()
		return h.fbSvc.SendQuickReplies(
			senderID,
			"⚠️ Invalid selection. Please choose a valid option from the menu.",
			quickReplies,
		)
	}

}
