package service

import (
	"camControl/internal/models"
	"camControl/internal/repository"
)

type cameraService struct {
	Repo repository.CameraRepository
}

type CameraService interface {
	GetCamera(cameraId uint) (*models.Camera, error)
}

// GetCamera - TODO implement me
func (s *cameraService) GetCamera(cameraId uint) (*models.Camera, error) {
	camera, err := s.Repo.GetCameraByID(cameraId)
	if err != nil {
		return nil, err
	}

	return camera, nil
}

func New(repo repository.CameraRepository) CameraService {
	return &cameraService{Repo: repo}
}
