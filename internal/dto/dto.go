package dto

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/mBayzigitov/dynamic-content-service/internal/models"
	"github.com/mBayzigitov/dynamic-content-service/internal/util/serverr"
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

type CreateBannerResponseDto struct {
	BannerId    int64  `json:"banner_id"`
}

type GetBannerResponseDto struct {
	Content     json.RawMessage `json:"content"`
}

type DeleteBannerResponseDto struct {
	Description string `json:"description"`
}

// ///////////////////// TYPES INIT METHODS ///////////////////////
func NewGetBannerResponse(banner *models.BannerModel) *GetBannerResponseDto {
	return &GetBannerResponseDto{
		Content:     banner.Content,
	}
}

func NewCreateBannerResponse(banner_id int64) *CreateBannerResponseDto {
	return &CreateBannerResponseDto{
		BannerId:    banner_id,
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
