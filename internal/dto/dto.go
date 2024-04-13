package dto

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/mBayzigitov/dynamic-content-service/internal/models"
	"github.com/mBayzigitov/dynamic-content-service/internal/util/serverr"
	"time"
)

// ///////////////////// TYPES ///////////////////////
type ValidationEntity interface {
	Validate(v *validator.Validate) *serverr.ApiError
}

type CreateBannerDto struct {
	TagIds    []int64         `json:"tag_ids" validate:"required"`
	FeatureId int64           `json:"feature_id" validate:"required"`
	Content   json.RawMessage `json:"content" validate:"required"`
	IsActive  bool            `json:"is_active"`
}

type ChangeBannerDto struct {
	TagIds    []int64          `json:"tag_ids"`
	FeatureId *int64           `json:"feature_id"`
	Content   *json.RawMessage `json:"content"`
	IsActive  *bool            `json:"is_active"`
}

type CreateBannerResponseDto struct {
	BannerId int64 `json:"banner_id"`
}

type GetBannerResponseDto struct {
	Content json.RawMessage `json:"content"`
}

type FilterBannersResponseDto struct {
	BannerId  int64           `json:"banner_id"`
	TagIds    []int64         `json:"tag_ids"`
	FeatureId int64           `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  bool            `json:"is_active"`
	ToDelete  bool            `json:"to_delete"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// ///////////////////// TYPES INIT METHODS ///////////////////////
func NewGetBannerResponse(banner *models.BannerModel) *GetBannerResponseDto {
	return &GetBannerResponseDto{
		Content: banner.Content,
	}
}

func NewCreateBannerResponse(banner_id int64) *CreateBannerResponseDto {
	return &CreateBannerResponseDto{
		BannerId: banner_id,
	}
}

func NewFilterBannersResponseDto(b models.BannerTagsModel) FilterBannersResponseDto {
	return FilterBannersResponseDto{
		BannerId:  b.Id,
		TagIds:    b.TagIds,
		FeatureId: b.FeatureId,
		Content:   b.Content,
		IsActive:  b.IsActive,
		ToDelete:  b.ToDelete,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

// ///////////////////// HELPER FUNCTIONS ///////////////////////

func (cbd *CreateBannerDto) Validate(v *validator.Validate) *serverr.ApiError {
	if err := v.Struct(cbd); err != nil {
		var verrs validator.ValidationErrors
		errors.As(err, &verrs)

		var errBody string
		if verrs != nil && len(verrs) > 0 {
			f := verrs[0]
			errBody = "field '" + f.Field() + "' validation failed: '" + f.ActualTag() + "' is violated"
		} else {
			errBody = "validation error"
		}

		apierr := serverr.NewInvalidRequestError(errBody)

		return apierr
	}

	return nil
}

func (cbd *ChangeBannerDto) Validate(v *validator.Validate) *serverr.ApiError {
	if err := v.Struct(cbd); err != nil {
		var verrs validator.ValidationErrors
		errors.As(err, &verrs)

		var errBody string
		if verrs != nil && len(verrs) > 0 {
			f := verrs[0]
			errBody = "field '" + f.Field() + "' validation failed: '" + f.ActualTag() + "' is violated"
		} else {
			errBody = "validation error"
		}

		apierr := serverr.NewInvalidRequestError(errBody)

		return apierr
	}

	return nil
}

func JsonBody(dto any) string {
	resp, _ := json.Marshal(dto)
	return string(resp)
}

func (cbd *CreateBannerDto) ToModel() *models.BannerTagsModel {
	return &models.BannerTagsModel{
		TagIds:    cbd.TagIds,
		FeatureId: cbd.FeatureId,
		Content:   cbd.Content,
		IsActive:  cbd.IsActive,
	}
}
