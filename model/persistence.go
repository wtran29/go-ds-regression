package model

import "time"

type ModelMetadata struct {
	SavedAt     time.Time `json:"saved_at"`
	Description string    `json:"description"`
	NumSamples  int       `json:"num_samples,omitempty"`
	Version     string    `json:"version,omitempty"`
}
