package dto

import (
	"encoding/json"
	"github.com/mBayzigitov/dynamic-content-service/internal/models"
)

type CreateBannerDto struct {
	TagIds    []int64         `json:"tag_ids" validate:"required"`
	FeatureId int64           `json:"feature_id" validate:"required"`
	Content   json.RawMessage `json:"content" validate:"required"`
	IsActive  bool            `json:"is_active"`
}

type GetBannerDto struct {
	TagId           int64 `json:"tag_id"`
	FeatureId       int64 `json:"feature_id"`
	UseLastRevision int64 `json:"use_last_revision"`
}

type CreateBannerResponseDto struct {
	Description string `json:"description"`
	BannerId    int64  `json:"banner_id"`
}

func NewCreateBannerResponse(banner_id int64) *CreateBannerResponseDto {
	return &CreateBannerResponseDto{
		Description: "Created",
		BannerId:    banner_id,
	}
}

func (cbrd *CreateBannerResponseDto) JsonBody() string {
	resp, _ := json.Marshal(cbrd)
	return string(resp)
}

func (dto *CreateBannerDto) ToModel() *models.BannerModel {
	return &models.BannerModel{
		TagIds:    dto.TagIds,
		FeatureId: dto.FeatureId,
		Content:   dto.Content,
		IsActive:  dto.IsActive,
	}
}
