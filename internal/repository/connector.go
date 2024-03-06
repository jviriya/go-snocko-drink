package repository

import (
	"go-pentor-bank/internal/infra/mongodb"
	"go-pentor-bank/internal/repository/commondb"
)

type CommonDBRepository struct {
	CommonDBRepo *commondb.SecureRepository
}

func NewCommonDBRepository(secureClient, secureCliShard mongodb.SecureClient) *CommonDBRepository {
	commonDBRepo := commondb.NewRepositoryV2(secureClient, secureCliShard)
	return &CommonDBRepository{
		CommonDBRepo: commonDBRepo,
	}
}
