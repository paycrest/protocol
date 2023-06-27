package sender

import "github.com/paycrest/paycrest-protocol/services"

// LoginResponse dummy comment
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}


type ReceiveAddressController struct {
	service *services.ReceiveAddressService
}