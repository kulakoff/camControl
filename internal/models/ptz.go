package models

import "time"

type PTZRequest struct {
	CameraID int    `json:"camera_id"`
	Action   string `json:"action"`
	Preset   string `json:"preset,omitempty"`
}

type PTZAction string

const (
	PTZLeft  PTZAction = "left"
	PTZRight PTZAction = "right"
	PTZUp    PTZAction = "up"
	PTZDown  PTZAction = "down"
	PTZStop  PTZAction = "stop"
	MinStep  float64   = 1
	Duration           = 500 * time.Millisecond
)
