package repository

import (
	"camControl/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type cameraRepository struct {
	DB     *pgxpool.Pool
	logger *slog.Logger
}

type CameraRepository interface {
	GetCameraByID(uint) (*models.Camera, error)
}

func NewCameraRepository(db *pgxpool.Pool, logger *slog.Logger) CameraRepository {
	return &cameraRepository{DB: db, logger: logger}
}

func (r *cameraRepository) GetCameraByID(cameraId uint) (*models.Camera, error) {
	r.logger.Debug("cameraRepository | GetCameraByID", "cameraId", cameraId)
	ctx := context.Background()

	camera := &models.Camera{}

	query := `SELECT id, ip, login, password FROM cameras WHERE id=$1`
	err := r.DB.QueryRow(ctx, query, cameraId).Scan(
		&camera.ID,
		&camera.IP,
		&camera.Login,
		&camera.Password,
	)
	if err != nil {
		r.logger.Error("Error getting camera by id", "id", cameraId)
		return nil, err
	}

	return camera, nil
}
