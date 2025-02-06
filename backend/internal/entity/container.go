package entity

import "time"

type Container struct {
	IpAddr         string    `json:"ip"`
	PingTime       time.Time `json:"ping_time"`
	LastSuccessful time.Time `json:"last_successful"`
}
