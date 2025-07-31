package service

import (
	"camControl/internal/models"
	onvif_client "camControl/internal/pkg/onvif"
	"camControl/internal/repository"
	"context"
	"log/slog"
	"sync"
)

type PTZService interface {
	MoveCamera(ctx context.Context, req *models.PTZRequest) error
	GetPresets(ctx context.Context, cameraID uint) ([]onvif_client.PTZPreset, error)
	GoToPreset(ctx context.Context, cameraID uint, presetToken string) error
	//getController(ctx context.Context, cameraID uint) (*onvif_client.PTZController, error)
}

type ptzService struct {
	controllers sync.Map // cache controllers
	camRepo     repository.CameraRepository
}

func NewPTZService(repo repository.CameraRepository) PTZService {
	return &ptzService{
		camRepo: repo,
	}
}

func (s *ptzService) getController(ctx context.Context, cameraID uint) (*onvif_client.PTZController, error) {
	// 01 - check from cache
	if ctrl, ok := s.controllers.Load(cameraID); ok {
		slog.Info("getController", "cam", cameraID)
		return ctrl.(*onvif_client.PTZController), nil
	}

	// 02 - get camera data
	camera, err := s.camRepo.GetCameraByID(cameraID)
	if err != nil {
		return nil, err
	}

	// 03 - make new camera controller
	// TODO: refactor port and minStep
	ctrl, err := onvif_client.New(camera.IP, "80", camera.Login, camera.Password, 1)
	if err != nil {
		return nil, err
	}

	// 04 - store controller
	s.controllers.Store(cameraID, ctrl)

	return ctrl, nil
}

func (s *ptzService) MoveCamera(ctx context.Context, req *models.PTZRequest) error {
	//TODO implement me
	ctrl, err := s.getController(ctx, uint(req.CameraID))
	if err != nil {
		return err
	}

	// FIXME
	return ctrl.Move(models.PTZAction(req.Action), req.Speed)
}

func (s *ptzService) GetPresets(ctx context.Context, cameraID uint) ([]onvif_client.PTZPreset, error) {
	slog.Info("PTZService | GetPresets")
	ctrl, err := s.getController(ctx, cameraID)
	if err != nil {
		return nil, err
	}
	return ctrl.GetPresets(ctx)
}

func (s *ptzService) GoToPreset(ctx context.Context, cameraID uint, presetToken string) error {
	slog.Info("PTZService | GoToPreset", "cameraID", cameraID, "presetToken", presetToken)
	ctrl, err := s.getController(ctx, cameraID)
	if err != nil {
		return err
	}
	return ctrl.GotoPreset(ctx, presetToken)
}
