package config

type RSS struct {
	Links         []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}
