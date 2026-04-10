package dto

type SIWEMessageRequest struct {
	Address string `json:"address"`
}

type SIWEMessageResponse struct {
	Message string `json:"message"`
}

type SIWEVerifyRequest struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

type SIWEVerifyResponse struct {
	Authenticated bool   `json:"authenticated"`
	Address       string `json:"address"`
}

type AuthMeResponse struct {
	Authenticated   bool   `json:"authenticated"`
	Address         string `json:"address,omitempty"`
	ChainID         string `json:"chainId,omitempty"`
	IsPlatformWallet bool  `json:"isPlatformWallet"`
}

type AuthLogoutResponse struct {
	Success bool `json:"success"`
}