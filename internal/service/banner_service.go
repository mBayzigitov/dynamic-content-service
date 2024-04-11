package service

import (
	"encoding/json"
	"github.com/mBayzigitov/dynamic-content-service/internal/repo"
	"github.com/mBayzigitov/dynamic-content-service/internal/util/serverr"
	"go.uber.org/zap"
)

type BannerService struct {
	l  *zap.SugaredLogger
	br *repo.BannerRepository
}

func NewBannerService(br *repo.BannerRepository) *BannerService {
	loginst, _ := zap.NewDevelopment()
	return &BannerService{
		br: br,
		l:  loginst.Sugar(),
	}
}

func (bs *BannerService) CreateBanner(tagIds []int64, featureId int64, content json.RawMessage, isActive bool) *serverr.ApiError {
	// check if feature is present

	// check if tags are present

	// create banner

	return nil
}
