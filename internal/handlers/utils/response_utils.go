package utils

import (
	"fmt"
	"school-assistant-wh/internal/constants"
	"school-assistant-wh/internal/repositories"
	"school-assistant-wh/internal/services/facebook"
	"school-assistant-wh/internal/services/helpers"
)

type ResponseUtils struct {
	repo     repositories.UserRepository
	linkRepo repositories.UserLinkRepository
	fbSvc    *facebook.Service
}

func NewResponseUtils(repo repositories.UserRepository, linkRepo repositories.UserLinkRepository, fbSvc *facebook.Service) *ResponseUtils {
	return &ResponseUtils{
		repo:     repo,
		linkRepo: linkRepo,
		fbSvc:    fbSvc,
	}
}

func (u *ResponseUtils) SendResponseWithQuickReplies(senderID, message string) error {
	user, _ := u.repo.GetUserByPSID(senderID)
	var currentProfileData interface{}

	if user != nil && !user.IsActive {
		quickReplies := helpers.GetQuickReplies(constants.UserStatusUnregistered)
		return u.fbSvc.SendQuickReplies(senderID, constants.AccountDeactivatedMessage, quickReplies)
	}

	if user != nil {
		var err error
		currentProfileData, err = u.linkRepo.GetPrimaryLink(int(user.ID))
		if err != nil {
			fmt.Printf("error getting user links: %v\n", err)
		}
	}

	status := constants.UserStatusUnregistered
	if currentProfileData != nil {
		status = constants.UserStatusLinkedPrimary
	} else if user != nil && user.IsActive {
		status = constants.UserStatusRegistered
	}
	quickReplies := helpers.GetQuickReplies(status)
	return u.fbSvc.SendQuickReplies(senderID, message, quickReplies)
}
