package onvif_client

import (
	"camControl/internal/models"
	"context"
	"encoding/xml"
	"fmt"
	goonvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/device"
	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/ptz"
	sdk_device "github.com/use-go/onvif/sdk/device"
	sdk_media "github.com/use-go/onvif/sdk/media"
	sdk_ptz "github.com/use-go/onvif/sdk/ptz"
	"github.com/use-go/onvif/xsd/onvif"
	"io"
	"log/slog"
	"time"
)

type Preset struct {
	Token string
	Name  string
}

type PTZPreset struct {
	Token string `xml:"token,attr"`
	Name  string `xml:"Name"`
}

type GetPresetsResponse struct {
	Presets []PTZPreset `xml:"Body>GetPresetsResponse>Preset"`
}

type PTZController struct {
	dev          *goonvif.Device
	profileToken string
	logger       *slog.Logger
	//minStep      float64
}

const SPEED time.Duration = 100

func New(ip, port, username, password string, minStep float64, logger *slog.Logger) (*PTZController, error) {
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

	logger.Debug("DEBUG | Initialized camera controller", "ip", ip, "port", port, "profile", profileToken)

	return &PTZController{
		dev,
		profileToken,
		logger,
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

func (c *PTZController) Move(direction models.PTZAction, speed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	if speed > 1000 {
		return fmt.Errorf("speed too high")
	}

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

	slog.Info("Move", "X", x, "Y", y)

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

	time.Sleep(time.Duration(speed) * time.Millisecond)
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

func (c *PTZController) GotoPreset(ctx context.Context, presetToken string) error {
	response, err := sdk_ptz.Call_GotoPreset(ctx, c.dev, ptz.GotoPreset{
		ProfileToken: onvif.ReferenceToken(c.profileToken),
		PresetToken:  onvif.ReferenceToken(presetToken),
	})
	if err != nil {
		return fmt.Errorf("call GotoPreset failed: %v", err)
	}
	slog.Info("GotoPreset response", "response", response)
	return nil
}

func (c *PTZController) SetPreset(ctx context.Context, presetToken string) (string, error) {
	response, err := sdk_ptz.Call_SetPreset(ctx, c.dev, ptz.SetPreset{
		ProfileToken: onvif.ReferenceToken(c.profileToken),
		PresetToken:  onvif.ReferenceToken(presetToken),
	})
	if err != nil {
		return "", fmt.Errorf("call SetPreset failed: %v", err)
	}

	return string(response.PresetToken), nil
}

func (c *PTZController) RemovePreset(ctx context.Context, presetToken string) error {
	slog.Info("onvif_client | RemovePreset", "presetToken", presetToken)
	_, err := sdk_ptz.Call_RemovePreset(ctx, c.dev, ptz.RemovePreset{
		ProfileToken: onvif.ReferenceToken(c.profileToken),
		PresetToken:  onvif.ReferenceToken(presetToken),
	})
	if err != nil {
		return fmt.Errorf("call RemovePreset failed: %v", err)
	}

	return nil
}

// GetPresets custom getPreset method
func (c *PTZController) GetPresets(ctx context.Context) ([]PTZPreset, error) {
	//type PTZPreset struct {
	//	Token string `xml:"token,attr"`
	//	Name  string `xml:"Name"`
	//}
	//
	//type GetPresetsResponse struct {
	//	Presets []PTZPreset `xml:"Body>GetPresetsResponse>Preset"`
	//}

	request := ptz.GetPresets{
		ProfileToken: onvif.ReferenceToken(c.profileToken),
	}

	rawResponse, err := c.dev.CallMethod(request)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetPresets: %v", err)
	}
	defer rawResponse.Body.Close()

	body, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	//slog.Info("GetPresets raw response", "response", string(body))

	var response GetPresetsResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse GetPresets response: %v", err)
	}

	//slog.Info("Parsed GetPresets response", "response", response)

	var presets []PTZPreset
	for _, preset := range response.Presets {
		presets = append(presets, PTZPreset{
			Token: preset.Token,
			Name:  preset.Name,
		})
	}

	return presets, nil
}
