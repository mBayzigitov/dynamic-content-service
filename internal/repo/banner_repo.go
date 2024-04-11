package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type BannerRepository struct {
	p *pgxpool.Pool
	l *zap.SugaredLogger
}

func NewBannerRepository(p *pgxpool.Pool) *BannerRepository {
	logger, _ := zap.NewDevelopment()

	return &BannerRepository{
		p: p,
		l: logger.Sugar(),
	}
}

func (br *BannerRepository) DoesFeatureExist(featureId int64) bool {
	return true
}

func (br *BannerRepository) DoTagsExist(tagsIds []int64) bool {
	return true
}

func (br *BannerRepository) CreateBanner()  {
	// TODO
}