package models

import (
	"encoding/json"
	"time"
)

type BannerModel struct {
	TagIds    []int64
	FeatureId int64
	Content   json.RawMessage
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}