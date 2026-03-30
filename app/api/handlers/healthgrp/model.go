package healthgrp

import (
	"database/sql"
	"time"
)

type HealthResponse struct {
	Status        string      `json:"status,omitempty"`
	Uptime        string      `json:"uptime,omitempty"`
	Timestamp     time.Time   `json:"timestamp,omitempty"`
	Since         time.Time   `json:"since,omitempty"`
	Version       string      `json:"version,omitempty"`
	GOMAXPROCS    string      `json:"GOMAXPROCS,omitempty"`
	DatabaseStats sql.DBStats `json:"DatabaseStats,omitempty"`
}
