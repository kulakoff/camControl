package camera

import (
	"camControl/internal/models"
	"context"
	"fmt"
	"github.com/use-go/onvif"
	"github.com/use-go/onvif/ptz"
	"github.com/use-go/onvif/xsd/onvif"
)

type CameraService struct{}

func NewCameraService() *CameraService {
	return &CameraService{}
}
func (s *CameraService) ControlCamera(ctx context.Context, camera *models.Camera, req *models.PTZRequest) error {
	dev, err := onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    camera.IP,
		Username: camera.Login,
		Password: camera.Password,
	})
	if err != nil {
		return err
	}
	switch req.Action {
	case "move_left":
		return s.movePTZ(dev, -0.1, 0, 0)
	case "move_right":
		return s.movePTZ(dev, 0.1, 0, 0)
	case "move_up":
		return s.movePTZ(dev, 0, 0.1, 0)
	case "move_down":
		return s.movePTZ(dev, 0, -0.1, 0)
	case "zoom_in":
		return s.movePTZ(dev, 0, 0, 0.1)
	case "zoom_out":
		return s.movePTZ(dev, 0, 0, -0.1)
	case "goto_preset":
		return s.gotoPreset(dev, req.Preset)
	default:
		return fmt.Errorf("unsupported action: %s", req.Action)
	}
}
func (s *CameraService) movePTZ(dev *onvif.Device, pan, tilt, zoom float64) error {
	ptzReq := ptz.ContinuousMove{
		ProfileToken: onvif.ReferenceToken("Profile_1"),
		Velocity: onvif.PTZSpeed{
			PanTilt: onvif.Vector2D{X: pan, Y: tilt},
			Zoom:    onvif.Vector1D{X: zoom},
		},
	}
	_, err := dev.CallMethod(ptzReq)
	return err
}
func (s *CameraService) gotoPreset(dev *onvif.Device, presetName string) error {
	presetReq := ptz.GotoPreset{
		ProfileToken: onvif.ReferenceToken("Profile_1"),
		PresetToken:  onvif.ReferenceToken(presetName),
	}
	_, err := dev.CallMethod(presetReq)
	return err
}
