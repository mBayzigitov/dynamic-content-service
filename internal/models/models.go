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
	ToDelete  bool
}

type BannerTagsModel struct {
	Id           int64
	TagIds       []int64
	FeatureId    int64
	Content      json.RawMessage
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastRevision int64
	ToDelete     bool
}

type BannerVersion struct {
	BannerId  string          `json:"banner_id"`
	Version   int64           `json:"version"`
	FeatureId int64           `json:"feature_id"`
	Tags      string          `json:"tags"`
	Content   json.RawMessage `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
}
