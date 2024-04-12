package service

import (
	"encoding/json"
	"fmt"
	"github.com/mBayzigitov/dynamic-content-service/internal/models"
	"github.com/mBayzigitov/dynamic-content-service/internal/repo"
	"github.com/mBayzigitov/dynamic-content-service/internal/util/serverr"
	"go.uber.org/zap"
	"time"
)

const RedisTtl = 5 * time.Minute

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

func (bs *BannerService) GetBanner(tagId int64, featureId int64, useLastRevision bool) (models.BannerModel, *serverr.ApiError) {
	// check use_last_revision flag
	// if TRUE -> get from database directly
	// if FALSE -> try to get from redis cache, if fails -> get from database directly
	key := fmt.Sprintf("%d_%d", featureId, tagId)

	if !useLastRevision {
		bc, err := bs.redis.Get(key)

		if err == nil {
			resp := models.BannerModel{
				Content: json.RawMessage(bc),
			}
			bs.l.Infof("get banner from cache with key '%s'", key)

			return resp, nil // return if key in cache is present
		} else {
			// just log if no such key found
			bs.l.Infof("redis: no key '%s' in cache found", key)
		}
	}

	// get from database
	banner, err := bs.br.GetBanner(
		tagId,
		featureId,
	)
	if err != nil {
		bs.l.Info(err.Error())
		return models.BannerModel{}, serverr.StorageError
	}

	if banner.Id == 0 {
		bs.l.Info(serverr.BannerNotFoundError)
		return banner, serverr.BannerNotFoundError
	}

	if !useLastRevision	{
		err = bs.redis.Set(
			key,
			string(banner.Content),
			RedisTtl,
		)
		if err != nil {
			bs.l.Fatal(err)
		}

		bs.l.Infof("Banner [%d] is cached, key: %s", banner.Id, key)
	}

	return banner, nil
}

func (bs *BannerService) CreateBanner(banner *models.BannerTagsModel) (int64, *serverr.ApiError) {
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

	createdId, err := bs.br.CreateBanner(banner)
	if err != nil {
		bs.l.Error(err.Error())
		return -1, serverr.StorageError
	}

	return createdId, nil
}
