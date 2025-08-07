package models

import "time"

type PTZRequest struct {
	CameraID int    `json:"camera_id"`
	Action   string `json:"action"`
	Speed    int    `json:"speed,omitempty"`
	//Preset   string `json:"preset,omitempty"`
}
type PTZRequestPreset struct {
	CameraID    int    `json:"camera_id"`
	PresetToken string `json:"preset_token"`
}

type PTZAction string

const (
	PTZLeft  PTZAction = "left"
	PTZRight PTZAction = "right"
	PTZUp    PTZAction = "up"
	PTZDown  PTZAction = "down"
	PTZStop  PTZAction = "stop"
	MinStep  float64   = 0.5
	Duration           = time.Second * 3
)
