package model

type AttributorData struct {
	TraderAddress        string `json:"traderAddr"`
	TraderAcct           string `json:"traderAcct"`
	Deposited            int    `json:"deposited"`
	Withdrawn            int    `json:"withdrawn"`
	ReferredCode         string `json:"referredCode"`
	VaultsCreated        int    `json:"vaultsCreated"`
	ManualInvested       int    `json:"manualInv"`
	SubscriptionInvested int    `json:"subscriptionInv"`
	VaultsInvestedIn     int    `json:"vaultsInvestedIn"`
	TotalInvested        int    `json:"totalInv"`
	Subscriptions        int    `json:"subscriptions"`
	StakedAmount         int    `json:"stakedAmount"`
}
