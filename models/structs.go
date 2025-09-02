package models

import "time"

type Watch struct {
	ID             int       `json:"id"`
	URL            string    `json:"url"`
	IntervalMS     int       `json:"interval_ms"`
	ExpectedStatus int       `json:"expected_status"`
	CreatedAt      time.Time `json:"created_at"`
}

type Result struct {
	ID             int
	WatchID        int
	Status         int
	ResponseTimeMS int
	CreatedAt      time.Time
}

type Alert struct {
	ID        string
	WatchID   string
	Message   string
	CreatedAt time.Time
}
