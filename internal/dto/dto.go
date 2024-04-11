package dto

import "encoding/json"

type CreateBannerDto struct {
	TagIds    []int64 `json:"tag_ids" validate:"required"`
	FeatureId int64 `json:"feature_id" validate:"required"`
	Content   json.RawMessage `json:"content" validate:"required"`
	IsActive  bool `json:"is_active"`
}