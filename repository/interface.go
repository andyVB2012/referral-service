package repository

import (
	"context"

	"github.com/andyVB2012/referral-service/model"
)

type Repository interface {
	CreateReferralCode(ctx context.Context, traderAddr string) (model.Referral, error)
	IsTraderAddrInDb(ctx context.Context, traderAddr string) bool
	IsCodeInDb(ctx context.Context, code string) bool
	GetCode(ctx context.Context, traderAddr string) (string, error)
	IsTraderAccountInDb(ctx context.Context, traderAcct string) bool
	AddTraderAccount(ctx context.Context, traderAcct string, traderAddr string) error
	GetTraderAddrFromTraderAcct(ctx context.Context, traderAcct string) (string, error)
	AddVariable(ctx context.Context, traderAddr string, variable string, value int) error
	Unsubscribe(ctx context.Context, traderAddr string) error
	AddTraderAcctIfNotExists(ctx context.Context, traderAcct string, traderAddr string) error
	AddAttributor(ctx context.Context, refCode string, traderAddr string) error
	GetStats(ctx context.Context, code string) (model.AttributionStats, error)
	GetAllAttributors(ctx context.Context, code string) ([]model.AttributorData, error)
}
