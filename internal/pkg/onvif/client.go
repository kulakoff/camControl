package onvif_client

import (
	"camControl/internal/models"
	"context"
	"fmt"
	goonvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/device"
	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/ptz"
	sdk_device "github.com/use-go/onvif/sdk/device"
	sdk_media "github.com/use-go/onvif/sdk/media"
	sdk_ptz "github.com/use-go/onvif/sdk/ptz"
	"github.com/use-go/onvif/xsd/onvif"
	"log/slog"
	"time"
)

type PTZController struct {
	dev          *goonvif.Device
	profileToken string
	//minStep      float64
}

func New(ip, port, username, password string, minStep float64) (*PTZController, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dev, err := goonvif.NewDevice(goonvif.DeviceParams{
		Xaddr:    ip + ":" + port,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("ONVIF connection failed: %w", err)
	}

	profileToken, err := getProfileToken(ctx, dev)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile token: %w", err)
	}

	slog.Info("DEBUG | Initialized camera controller", "ip", ip, "port", port, "profile", profileToken)

	return &PTZController{
		dev,
		profileToken,
		//minStep,
	}, nil
}

// getProfileToken - get first profile token
func getProfileToken(ctx context.Context, dev *goonvif.Device) (string, error) {
	profiles, err := sdk_media.Call_GetProfiles(ctx, dev, media.GetProfiles{})
	if err != nil {
		return "", err
	}

	if len(profiles.Profiles) == 0 {
		return "", fmt.Errorf("no profiles found")
	}

	token := string(profiles.Profiles[0].Token)
	slog.Info("DEBUG | Retrieved profile token", "token", token)

	return token, nil
}

func (c *PTZController) GetDeviceInfo(ctx context.Context) error {
	timeResponse, err := sdk_device.Call_GetSystemDateAndTime(ctx, c.dev, device.GetSystemDateAndTime{})
	if err != nil {
		return err
	}
	slog.Info("GetDeviceInfo", "sysInfo", timeResponse)

	infoResp, err := sdk_device.Call_GetDeviceInformation(ctx, c.dev, device.GetDeviceInformation{})
	if err != nil {
		return err
	}
	fmt.Printf("Информация об устройстве: %+v\n", infoResp)

	return nil
}

func (c *PTZController) Move(direction models.PTZAction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := c.CheckPTZSupport(ctx)
	if err != nil {
		return err
	}

	var x, y float64

	switch direction {
	case models.PTZLeft:
		x = -models.MinStep
	case models.PTZRight:
		x = models.MinStep
	case models.PTZUp:
		y = models.MinStep
	case models.PTZDown:
		y = -models.MinStep
	}

	// cal ptz action
	_, err = sdk_ptz.Call_ContinuousMove(ctx, c.dev, ptz.ContinuousMove{
		ProfileToken: onvif.ReferenceToken(c.profileToken),
		Velocity: onvif.PTZSpeed{
			PanTilt: onvif.Vector2D{X: x, Y: y},
		},
	})
	if err != nil {
		slog.Error("Failed to move camera", "direction", direction, "err", err)
		return err
	}

	time.Sleep(500 * time.Millisecond)
	return c.Stop(ctx)
}

func (c *PTZController) Stop(ctx context.Context) error {
	slog.Info("PTZController | stopping device")
	_, err := sdk_ptz.Call_Stop(ctx, c.dev, ptz.Stop{
		ProfileToken: onvif.ReferenceToken(c.profileToken),
		PanTilt:      true,
		Zoom:         false,
	})
	return err
}

func (c *PTZController) CheckPTZSupport(ctx context.Context) error {
	// get device capabilities
	_, err := sdk_device.Call_GetCapabilities(ctx, c.dev, device.GetCapabilities{
		Category: "PTZ",
	})
	if err != nil {
		return fmt.Errorf("get capabilities failed %w", err)
	}
	return nil
}
