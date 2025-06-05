package storage

import (
	"camControl/internal/config"
	"camControl/internal/models"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB *pgxpool.Pool
}

func NewStorage(conf *config.DbConfig) (*Storage, error) {
	// format connection str
	connStr := formatPostgresURL(conf)
	slog.Info("TEST", "connStr", connStr)

	// create connection pool
	psqlConf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DSN: %w", err)
	}
	psqlConf.MaxConns = 10
	psqlConf.MaxConnLifetime = 5 * time.Second

	// Connect to db
	db, err := pgxpool.NewWithConfig(context.Background(), psqlConf)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Check connection
	if err := db.Ping(context.Background()); err != nil {
		slog.Error("Unable to ping database", "error", err)
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	slog.Info("Connected to database")

	return &Storage{DB: db}, nil
}
func (s *Storage) Close() {
	if s.DB != nil {
		s.DB.Close()
		slog.Info("Closed database connection")
	}
}

// formatPostgresURL config to connection URI
func formatPostgresURL(cfg *config.DbConfig) string {
	//urlExample := "postgres://username:password@localhost:5432/database_name"
	//return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}

func (s *Storage) GetCameraByID(ctx context.Context, cameraId int) (*models.Camera, error) {
	camera := &models.Camera{}
	query := `SELECT id, ip, login, password FROM cameras WHERE id=$1`

	err := s.DB.QueryRow(ctx, query, cameraId).Scan(
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
