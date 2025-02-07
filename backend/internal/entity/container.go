package entity

import "time"

type Container struct {
	IpAddr         string    `json:"ip"`
	PingTime       int       `json:"ping_time"`
	LastSuccessful time.Time `json:"last_successful"`
}

type PingContainer struct {
	IpAddr         string    `json:"ip"`
	PingTime       int       `json:"ping_time"`
	IsSuccessful   bool      `json:"is_successful"`
	LastSuccessful time.Time `json:"last_successful"`
}
