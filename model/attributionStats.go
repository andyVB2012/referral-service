package model

type AttributionStats struct {
	RefferedAmount            int `json:"refferedAmount"`
	TotalDeposited            int `json:"totalDeposited"`
	TotalWithdrawn            int `json:"totalWithdrawn"`
	TotalVaultCreated         int `json:"totalVaultCreated"`
	TotalManualInvested       int `json:"totalManualInvested"`
	TotalSubscriptionInvested int `json:"totalSubscriptionInvested"`
	TotalInvested             int `json:"totalInvested"`
	TotalVaultsInvestedIn     int `json:"totalVaultsInvestedIn"`
	TotalSubscriptions        int `json:"totalSubscriptions"`
	TotalStakedAmount         int `json:"totalStakedAmount"`
}
