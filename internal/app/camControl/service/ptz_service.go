package service

import (
	"camControl/internal/app/camControl/models"
	"camControl/internal/app/camControl/repository"
	"camControl/pkg/onvif"
	"context"
	"log/slog"
	"sync"
)

type PTZService interface {
	MoveCamera(ctx context.Context, req *models.PTZRequest) error
	GetPresets(ctx context.Context, cameraID uint) ([]onvif_client.PTZPreset, error)
	GoToPreset(ctx context.Context, cameraID uint, presetToken string) error
	SetPreset(ctx context.Context, cameraID uint, presetToken string) error
	RemovePreset(ctx context.Context, cameraID uint, presetToken string) error
	//getController(ctx context.Context, cameraID uint) (*onvif_client.PTZController, error)
}

type ptzService struct {
	controllers sync.Map // cache controllers
	camRepo     repository.CameraRepository
	logger      *slog.Logger
}

func NewPTZService(repo repository.CameraRepository, logger *slog.Logger) PTZService {
	return &ptzService{
		camRepo: repo,
		logger:  logger,
	}
}

func (s *ptzService) getController(ctx context.Context, cameraID uint) (*onvif_client.PTZController, error) {
	// 01 - check from cache
	if ctrl, ok := s.controllers.Load(cameraID); ok {
		s.logger.Info("getController", "cam", cameraID)
		return ctrl.(*onvif_client.PTZController), nil
	}

	// 02 - get camera data
	camera, err := s.camRepo.GetCameraByID(cameraID)
	if err != nil {
		return nil, err
	}

	// 03 - make new camera controller
	// TODO: refactor port and minStep
	ctrl, err := onvif_client.New(camera.IP, "80", camera.Login, camera.Password, 1, s.logger)
	if err != nil {
		return nil, err
	}

	// 04 - store controller
	s.controllers.Store(cameraID, ctrl)

	return ctrl, nil
}

func (s *ptzService) MoveCamera(ctx context.Context, req *models.PTZRequest) error {
	ctrl, err := s.getController(ctx, uint(req.CameraID))
	if err != nil {
		return err
	}

	return ctrl.Move(models.PTZAction(req.Action), req.Speed)
}

func (s *ptzService) GetPresets(ctx context.Context, cameraID uint) ([]onvif_client.PTZPreset, error) {
	ctrl, err := s.getController(ctx, cameraID)
	if err != nil {
		return nil, err
	}
	return ctrl.GetPresets(ctx)
}

func (s *ptzService) GoToPreset(ctx context.Context, cameraID uint, presetToken string) error {
	s.logger.Debug("PTZService | GoToPreset", "cameraID", cameraID, "presetToken", presetToken)
	ctrl, err := s.getController(ctx, cameraID)
	if err != nil {
		return err
	}
	return ctrl.GotoPreset(ctx, presetToken)
}

func (s *ptzService) SetPreset(ctx context.Context, cameraID uint, presetToken string) error {
	s.logger.Debug("PTZService | SetPreset", "cameraID", cameraID, "presetToken", presetToken)
	ctrl, err := s.getController(ctx, cameraID)
	if err != nil {
		return err
	}

	_, err = ctrl.SetPreset(ctx, presetToken)
	if err != nil {
		return err
	}
	return nil
}

func (s *ptzService) RemovePreset(ctx context.Context, cameraID uint, presetToken string) error {
	s.logger.Debug("PTZService | RemovePreset", "cameraID", cameraID, "presetToken", presetToken)
	ctrl, err := s.getController(ctx, cameraID)
	if err != nil {
		return err
	}

	err = ctrl.RemovePreset(ctx, presetToken)
	if err != nil {
		return err
	}
	return nil
}
