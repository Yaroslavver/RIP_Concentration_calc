package repository

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db            *gorm.DB
	minio         *minio.Client
	minioBucket   string
	minioEndpoint string
	minioUseSSL   bool
}

type Settings struct {
	PostgresDSN    string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

func New(settings *Settings) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(settings.PostgresDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(settings.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.MinioAccessKey, settings.MinioSecretKey, ""),
		Secure: settings.MinioUseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:            db,
		minio:         minioClient,
		minioBucket:   settings.MinioBucket,
		minioEndpoint: settings.MinioEndpoint,
		minioUseSSL:   settings.MinioUseSSL,
	}, nil
}