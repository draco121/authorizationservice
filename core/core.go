package core

import (
	"context"
	"github.com/draco121/authorizationservice/repository"
	"github.com/draco121/common/constants"
	"github.com/draco121/common/models"
	"time"
)

type IAuthorizationService interface {
	Authorize(ctx context.Context, authorizationInput *models.AuthorizationInput) *models.AuthorizationOutput
}

type authorizationService struct {
	IAuthorizationService
	repo repository.IAuthorizationRepository
}

func NewAuthorizationService(repository repository.IAuthorizationRepository) IAuthorizationService {
	us := authorizationService{
		repo: repository,
	}
	return &us
}

func (c *authorizationService) Authorize(ctx context.Context, authorizationInput *models.AuthorizationInput) *models.AuthorizationOutput {
	claims, err := c.repo.GetTokenClaims(authorizationInput.Token)
	if err != nil {
		log := models.AuthorizationLog{
			Grant:       constants.Rejected,
			RequestedAt: time.Now(),
			Actions:     authorizationInput.Actions,
			Reason:      err.Error(),
		}
		_ = c.repo.InsertAuthorizationLog(ctx, log)
		return &models.AuthorizationOutput{
			Grant:  constants.Rejected,
			UserId: "",
		}
	} else {
		allowedActions := c.repo.GetActions(claims.Role)
		if authorizationEngine(allowedActions, authorizationInput.Actions) {
			log := models.AuthorizationLog{
				Grant:       constants.Allowed,
				RequestedAt: time.Now(),
				Actions:     authorizationInput.Actions,
				UserId:      claims.UserId,
				Role:        claims.Role,
				Reason:      "permissions matched",
			}
			_ = c.repo.InsertAuthorizationLog(ctx, log)
			return &models.AuthorizationOutput{
				Grant:  constants.Allowed,
				UserId: claims.UserId,
			}
		} else {
			log := models.AuthorizationLog{
				Grant:       constants.Rejected,
				RequestedAt: time.Now(),
				Actions:     authorizationInput.Actions,
				UserId:      claims.UserId,
				Role:        claims.Role,
				Reason:      "permissions does not match",
			}
			_ = c.repo.InsertAuthorizationLog(ctx, log)
			return &models.AuthorizationOutput{
				Grant:  constants.Rejected,
				UserId: claims.UserId,
			}
		}
	}
}
