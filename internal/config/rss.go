package config

import "time"

type RSS struct {
	Links         []string
	RequestPeriod time.Duration
}
