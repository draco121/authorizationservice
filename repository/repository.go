package repository

import (
	"context"
	"github.com/draco121/common/clients"
	"github.com/draco121/common/constants"
	"github.com/draco121/common/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var actionRoleMap constants.ActionRoleMapping = map[constants.Role][]constants.Action{
	constants.Tenant:      {constants.Read, constants.Write},
	constants.TenantAdmin: {constants.Read, constants.Write, constants.Delete},
	constants.Root:        {constants.All},
}

type IAuthorizationRepository interface {
	GetActions(role constants.Role) []constants.Action
	GetTokenClaims(token string) (*models.JwtCustomClaims, error)
	InsertAuthorizationLog(ctx context.Context, log models.AuthorizationLog) error
}

type authorizationRepository struct {
	IAuthorizationRepository
	db                             *mongo.Database
	authenticationServiceApiClient clients.IAuthenticationServiceApiClient
}

func NewAuthorizationRepo(database *mongo.Database, authenticationServiceApiClient clients.IAuthenticationServiceApiClient) IAuthorizationRepository {
	repo := authorizationRepository{
		db:                             database,
		authenticationServiceApiClient: authenticationServiceApiClient,
	}
	return &repo
}

func (c *authorizationRepository) GetActions(role constants.Role) []constants.Action {
	return actionRoleMap[role]
}

func (c *authorizationRepository) GetTokenClaims(token string) (*models.JwtCustomClaims, error) {
	claims, err := c.authenticationServiceApiClient.Authenticate(token)
	if err != nil {
		return nil, err
	} else {
		return claims, nil
	}
}

func (c *authorizationRepository) InsertAuthorizationLog(ctx context.Context, log models.AuthorizationLog) error {
	_, err := c.db.Collection("authorization-log").InsertOne(ctx, log)
	return err
}
