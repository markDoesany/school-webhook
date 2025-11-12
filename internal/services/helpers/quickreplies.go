package helpers

import (
	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/services/facebook"
	"strings"
)

// validQuickReplyPayloads is a map of all valid quick reply payloads
var validQuickReplyPayloads = map[string]bool{
	"MENU":           true,
	"SWITCH PROFILE": true,
	"MY SA-ID":       true,
	"ABOUT US":       true,
	"TALK TO HUMAN":  true,
	"REGISTER":       true,
	"VIEW PROFILE":   true,
	"MAIN MENU":      true,
	"CONTINUE":       true,
	"NO":             true,
}

// IsQuickReplyPayload checks if the given message is a valid quick reply payload
func IsQuickReplyPayload(message string) bool {
	_, exists := validQuickReplyPayloads[strings.ToUpper(message)]
	return exists
}

// GetQuickReplies returns the appropriate quick replies based on user status
func GetQuickReplies(status string) []facebook.QuickReply {
	switch status {
	case constants.UserStatusRegistered:
		return []facebook.QuickReply{
			{
				ContentType: "text",
				Title:       "View Profile",
				Payload:     "VIEW_PROFILE",
			},
			{
				ContentType: "text",
				Title:       "My SA-ID",
				Payload:     "MY_SA_ID",
			},
			{
				ContentType: "text",
				Title:       "About Us",
				Payload:     "ABOUT_US",
			},
			{
				ContentType: "text",
				Title:       "Talk to Human",
				Payload:     "TALK_TO_HUMAN",
			},
		}
	case constants.UserStatusLinkedPrimary:
		return []facebook.QuickReply{
			{
				ContentType: "text",
				Title:       "Menu",
				Payload:     "MENU",
			},
			{
				ContentType: "text",
				Title:       "View Profile",
				Payload:     "VIEW_PROFILE",
			},
			{
				ContentType: "text",
				Title:       "Switch Profile",
				Payload:     "SWITCH_PROFILE",
			},
			{
				ContentType: "text",
				Title:       "My SA-ID",
				Payload:     "MY_SA_ID",
			},
		}
	default: // Unregistered
		return []facebook.QuickReply{
			{
				ContentType: "text",
				Title:       "Register",
				Payload:     "REGISTER",
			},
			{
				ContentType: "text",
				Title:       "About Us",
				Payload:     "ABOUT_US",
			},
			{
				ContentType: "text",
				Title:       "Talk to Human",
				Payload:     "TALK_TO_HUMAN",
			},
		}
	}
}

// GetAccountManagementReplies returns quick replies for account management
func GetAccountManagementReplies() []facebook.QuickReply {
	return []facebook.QuickReply{
		{
			ContentType: "text",
			Title:       "Switch Profile",
			Payload:     "SWITCH_PROFILE",
		},
		{
			ContentType: "text",
			Title:       "View All Profiles",
			Payload:     "VIEW_PROFILE",
		},
		{
			ContentType: "text",
			Title:       "Back to Menu",
			Payload:     "MAIN_MENU",
		},
	}
}

// GetMainMenuReplies returns the main menu quick replies
func GetMainMenuReplies() []facebook.QuickReply {
	return []facebook.QuickReply{
		{
			ContentType: "text",
			Title:       "Switch Profile",
			Payload:     "SWITCH_PROFILE",
		},
		{
			ContentType: "text",
			Title:       "My SA-ID",
			Payload:     "MY_SA_ID",
		},
		{
			ContentType: "text",
			Title:       "Talk to Human",
			Payload:     "TALK TO HUMAN",
		},
	}
}

// GetProfileConfirmationReplies returns quick replies for profile confirmation
func GetProfileConfirmationReplies() []facebook.QuickReply {
	return []facebook.QuickReply{
		{
			ContentType: "text",
			Title:       "Continue",
			Payload:     "CONTINUE",
		},
		{
			ContentType: "text",
			Title:       "No",
			Payload:     "NO",
		},
	}
}

// GetProfileManagementReplies returns quick replies for profile management
func GetProfileManagementReplies() []facebook.QuickReply {
	return []facebook.QuickReply{
		{
			ContentType: "text",
			Title:       "Continue",
			Payload:     "CONTINUE",
		},
		{
			ContentType: "text",
			Title:       "Switch Profile",
			Payload:     "SWITCH PROFILE",
		},
	}
}

// GetBackToMainMenuReplies returns quick replies for going back to main menu
func GetBack() []facebook.QuickReply {
	return []facebook.QuickReply{
		{
			ContentType: "text",
			Title:       "Back",
			Payload:     "BACK",
		},
	}
}

// GetPaymentReplies returns quick replies for payment-related actions
func GetPaymentReplies() []facebook.QuickReply {
	return []facebook.QuickReply{
		{
			ContentType: "text",
			Title:       "Back",
			Payload:     "BACK",
		},
		// {
		// 	ContentType: "text",
		// 	Title:       "Pay Now",
		// 	Payload:     "PAY_NOW",
		// },
		{
			ContentType: "text",
			Title:       "Payment Logs",
			Payload:     "PAYMENT_LOGS",
		},
	}
}

// GetViewMoreReplies returns quick replies for viewing more items with pagination
func GetViewMoreReplies() []facebook.QuickReply {
	return []facebook.QuickReply{
		{
			ContentType: "text",
			Title:       "Back",
			Payload:     "MENU",
		},
		{
			ContentType: "text",
			Title:       "View More",
			Payload:     "VIEW MORE",
		},
	}
}

func GetConfirmProfileSwitch() []facebook.QuickReply {
	return []facebook.QuickReply{
		{
			ContentType: "text",
			Title:       "Back",
			Payload:     "BACK",
		},
		{
			ContentType: "text",
			Title:       "Proceed",
			Payload:     "PROCEED",
		},
	}
}

func GetAskSupportReplies() []facebook.QuickReply {
	return []facebook.QuickReply{
		{
			ContentType: "text",
			Title:       "Back",
			Payload:     "BACK",
		},
		{
			ContentType: "text",
			Title:       "View Tickets",
			Payload:     "VIEW_TICKETS",
		},
	}
}
