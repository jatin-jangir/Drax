package models

import "time"

type Alert struct {
	AlertName   string    `bson:"alertname"`
	Status      string    `bson:"status"` // firing/resolved
	Severity    string    `bson:"severity"` // critical/warning
	Instance    string    `bson:"instance"`
	Summary     string    `bson:"summary"`
	Description string    `bson:"description"`
	StartsAt    time.Time `bson:"startsAt"`
	EndsAt      time.Time `bson:"endsAt,omitempty"`
	CreatedAt   time.Time `bson:"createdAt"`
}

type AlertFilter struct {
	AlertName   string    `form:"alertname"`
	Status      string    `form:"status"`
	Severity    string    `form:"severity"`
	Instance    string    `form:"instance"`
	Search      string    `form:"search"`
	StartTime   time.Time `form:"start_time" time_format:"2006-01-02T15:04:05Z"`
	EndTime     time.Time `form:"end_time" time_format:"2006-01-02T15:04:05Z"`
	SortBy      string    `form:"sort_by"` // "alertname", "startsAt", etc.
	SortOrder   int       `form:"sort_order"` // 1=asc, -1=desc
}