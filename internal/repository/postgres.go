package repository

import (
	"camControl/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type cameraRepository struct {
	DB *pgxpool.Pool
}

type CameraRepository interface {
	GetCameraByID(uint) (*models.Camera, error)
}

func New(db *pgxpool.Pool) CameraRepository {
	return &cameraRepository{DB: db}
}

func (c *cameraRepository) GetCameraByID(cameraId uint) (*models.Camera, error) {
	//TODO implement me
	//panic("implement me")
	slog.Debug("cameraRepository | GetCameraByID")
	ctx := context.Background()

	camera := &models.Camera{}

	query := `SELECT id, ip, login, password FROM cameras WHERE id=$1`
	err := c.DB.QueryRow(ctx, query, cameraId).Scan(
		&camera.ID,
		&camera.IP,
		&camera.Login,
		&camera.Password,
	)
	if err != nil {
		slog.Error("Error getting camera by id", "id", cameraId)
		return nil, err
	}

	return camera, nil
}
