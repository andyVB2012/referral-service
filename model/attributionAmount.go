package model

type AttributionAmount struct {
	ManualInvestment       float32 `json:"manulaInvestment"`
	SubscriptionInvestment float32 `json:"subscriptionInvestment"`
	TotalAmount            float32 `json:"amount"`
}
