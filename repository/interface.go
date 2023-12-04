package repository

import (
	"context"

	"github.com/parsaakbari1209/go-mongo-crud-rest-api/model"
)

type Repository interface {
	CreateReferralCode(ctx context.Context, traderAddr string) (model.Code, error)
	GetReferralCode(ctx context.Context, traderAddr string) (model.Code, error)
}
