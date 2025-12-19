package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Terminal represents a terminal device in the system
type Terminal struct {
	ID           string    `db:"id"            json:"id"`
	SerialNumber string    `db:"serial_number" json:"serial_number"`
	Status       string    `db:"status"        json:"status"`
	LastSeen     time.Time `db:"last_seen"     json:"last_seen"`
	Metadata     JSONB     `db:"metadata"      json:"metadata"`
	CreatedAT    time.Time `db:"created_at"    json:"created_at"`
}

// Job represents a job assigned to a terminal
type Job struct {
	ID         string    `db:"id"          json:"id"`
	TerminalID string    `db:"terminal_id" json:"terminal_id"`
	Type       string    `db:"type"        json:"type"`
	Payload    JSONB     `db:"payload"     json:"payload"`
	Status     string    `db:"status"      json:"status"`
	Result     JSONB     `db:"result"      json:"result"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"  json:"updated_at"`
}

// JSONB is used for fields that store JSON data, allowing flexibility in the structure of the metadata, payload, and result fields.
// Allows you to store arbitrary JSON data without defining a strict schema for those fields.
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	result := make(JSONB)
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*j = result
	return nil
}

// Constants status of Terminal
const (
	TerminalStatusIdle    = "idle"    // waiting for jobs
	TerminalStatusRunning = "running" // executing a job
	TerminalStatusOffline = "offline" // not reachable
)

// Constants status of Job
const (
	JobStatusPending = "pending" // job is created but not yet started
	JobStatusRunning = "running" // job is currently being executed
	JobStatusDone    = "done"    // job completed successfully
	JobStatusFailed  = "failed"  // job execution failed
)
