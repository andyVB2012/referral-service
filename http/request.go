package http

type AttributionRequest struct {
	Address   string `json:"address"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}
