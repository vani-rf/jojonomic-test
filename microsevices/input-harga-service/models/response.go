package models

type InputHargaResponse struct {
	IsError bool   `json:"error"`
	Message string `json:"message,omitempty"`
	ReffId  string `json:"reff_id,omitempty"`
}
