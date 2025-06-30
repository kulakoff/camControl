package service

import (
	"camControl/internal/models"
	"camControl/internal/repository"
	"log/slog"
)

type cameraService struct {
	Repo repository.CameraRepository
}

type CameraService interface {
	GetCamera(cameraId uint) (*models.Camera, error)
}

// GetCamera - TODO implement me
func (s *cameraService) GetCamera(cameraId uint) (*models.Camera, error) {

	slog.Info("cameraService | GetCamera")

	camera, err := s.Repo.GetCameraByID(cameraId)
	if err != nil {
		return nil, err
	}

	return camera, nil
}

func New(repo repository.CameraRepository) CameraService {
	return &cameraService{Repo: repo}
}
