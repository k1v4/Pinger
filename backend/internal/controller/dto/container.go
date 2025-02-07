package dto

import "time"

type NewContainerResponse struct {
	Ip string `json:"ip"`
}

type DtoPingContainer struct {
	PingTime       int       `json:"ping_time"`
	IsSuccessful   bool      `json:"is_successful"`
	LastSuccessful time.Time `json:"last_successful"`
}

type UpdateContainerRequest struct {
	PingTime       int       `json:"ping_time"`
	LastSuccessful time.Time `json:"last_successful"`
}

type AddContainerRequest struct {
	Ip             string    `json:"ip"`
	PingTime       int       `json:"ping_time"`
	LastSuccessful time.Time `json:"last_successful"`
}

type DeleteContainerResponse struct {
	IsSuccess bool `json:"is_success"`
}

type Message struct {
}
