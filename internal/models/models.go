package models

import (
	"encoding/json"
	"time"
)

type BannerModel struct {
	Id        int64
	TagId     int64
	FeatureId int64
	Content   json.RawMessage
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BannerTagsModel struct {
	TagIds    []int64
	FeatureId int64
	Content   json.RawMessage
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
