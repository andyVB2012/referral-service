package http

import "github.com/andyVB2012/referral-service/model"

type AttributionResponse struct {
	Address      string                 `json:"address"`
	ReferralCode string                 `json:"referralCode"`
	Stats        model.AttributionStats `json:"stats"`
	Attributors  []model.AttributorData `json:"attributors"`
}
