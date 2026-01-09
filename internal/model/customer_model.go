package model

type CustomerResponse struct {
	Name   string `json:"name,omitempty"`
	Points int    `json:"points,omitempty"`
}
