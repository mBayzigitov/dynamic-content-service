package service

import (
	"github.com/mBayzigitov/dynamic-content-service/internal/models"
	"github.com/mBayzigitov/dynamic-content-service/internal/repo"
	"github.com/mBayzigitov/dynamic-content-service/internal/util/serverr"
	"go.uber.org/zap"
)

type BannerService struct {
	l     *zap.SugaredLogger
	br    *repo.BannerRepository
	redis *repo.CacheRepo
}

func NewBannerService(br *repo.BannerRepository, redis *repo.CacheRepo) *BannerService {
	loginst, _ := zap.NewDevelopment()
	return &BannerService{
		br:    br,
		l:     loginst.Sugar(),
		redis: redis,
	}
}

func (bs *BannerService) GetBanner() {
	// TODO
}

func (bs *BannerService) CreateBanner(banner *models.BannerModel) (int64, *serverr.ApiError) {
	// check if feature is present
	featExists, err := bs.br.DoesFeatureExist(banner.FeatureId)
	if err != nil {
		bs.l.Error(err.Error())
		return -1, serverr.StorageError
	}

	if !featExists {
		return -1, serverr.NewInvalidRequestError("Указанный feature_id не существует")
	}

	// check if tags are present
	tagsExist, err := bs.br.DoTagsExist(banner.TagIds)
	if err != nil {
		bs.l.Error(err.Error())
		return -1, serverr.StorageError
	}

	if !tagsExist {
		return -1, serverr.NewInvalidRequestError("Все/некоторые tag_id не существуют")
	}

	// create banner
	bs.l.Infof("Passed is_active: %t", banner.IsActive)
	createdId, err := bs.br.CreateBanner(banner)
	if err != nil {
		bs.l.Error(err.Error())
		return -1, serverr.StorageError
	}

	return createdId, nil
}
