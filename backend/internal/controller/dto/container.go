package dto

import "time"

type NewContainerResponse struct {
	Ip string `json:"ip"`
}

type UpdateContainerRequest struct {
	PingTime       time.Time `json:"ping_time"`
	LastSuccessful time.Time `json:"last_successful"`
}

type DeleteContainerResponse struct {
	IsSuccess bool `json:"is_success"`
}
