package model

type AttributorData struct {
	Code
	AttributionAmount
	VaultsInvested int    `json:"vaultsInvested"`
	LastUpdated    string `json:"lastUpdated"`
	TotalRefferals int    `json:"totalRefferals"`
}
