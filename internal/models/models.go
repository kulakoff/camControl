package models

type Camera struct {
	ID       int    `json:"id"`
	IP       string `json:"ip"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type PTZRequest struct {
	CameraID int    `json:"camera_id"`
	Action   string `json:"action"`
	Preset   string `json:"preset,omitempty"`
}
