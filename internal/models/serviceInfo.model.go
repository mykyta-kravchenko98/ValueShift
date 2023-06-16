package models

import "time"

type ServiceInfo struct {
	Name        string
	UpTime      time.Time
	Environment string
	Version     string
}
