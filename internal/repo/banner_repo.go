package repo

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mBayzigitov/dynamic-content-service/internal/models"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type BannerRepository struct {
	p *pgxpool.Pool
	l *zap.SugaredLogger
}

func NewBannerRepository(p *pgxpool.Pool) *BannerRepository {
	logger, _ := zap.NewDevelopment()

	return &BannerRepository{
		p: p,
		l: logger.Sugar(),
	}
}

func (br *BannerRepository) DoesFeatureExist(featureID int64) (bool, error) {
	query := "SELECT COUNT(*) FROM features WHERE id = $1"

	var count int
	err := br.p.QueryRow(context.Background(), query, featureID).Scan(&count)
	if err != nil {
		return false, err
	}

	// check if the feature ID exists
	return count > 0, nil
}

func (br *BannerRepository) DoTagsExist(tagsIds []int64) (bool, error) {
	// idea is to compare count of tags from db and the actual slice len
	// if count(*) from tags where... = len(tagsIds) -> return true
	query := "SELECT COUNT(*) FROM tags WHERE id IN ("
	params := make([]string, len(tagsIds))
	for i, id := range tagsIds {
		params[i] = strconv.FormatInt(id, 10)
	}
	query += strings.Join(params, ",") + ")"

	var count int
	err := br.p.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		return false, err
	}

	// check if all tag IDs exist
	return count == len(tagsIds), nil
}

func (br *BannerRepository) GetBanner(tagId int64, featureId int64) (models.BannerModel, error) {
	var banner models.BannerModel

	// query with JOIN to select banner based on tagId, featureId, is_active=true, and to_delete=false
	query := `
		SELECT 
			b.id,
			b.content,
			b.feature_id,
			b.is_active,
			b.created_at,
			b.updated_at
		FROM 
			banners b
		JOIN 
			banners_tags bt ON b.id = bt.banner_id
		WHERE 
			bt.tag_id = $1
			AND b.feature_id = $2 
			AND b.is_active = true 
			AND b.to_delete = false
	`

	err := br.p.QueryRow(
		context.Background(),
		query,
		tagId,
		featureId,
	).Scan(
		&banner.Id,
		&banner.Content,
		&banner.FeatureId,
		&banner.IsActive,
		&banner.CreatedAt,
		&banner.UpdatedAt,
	)
	if err != nil {
		return models.BannerModel{}, err
	}

	return banner, nil
}

func (br *BannerRepository) CreateBanner(banner *models.BannerTagsModel) (int64, error) {
	// start a transaction
	tx, err := br.p.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	defer func() {
		if pm := recover(); pm != nil {
			tx.Rollback(context.Background())
			panic(pm)
		} else if err != nil {
			tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()

	// insert into banners table
	var bannerID int64
	var createdAt time.Time
	err = tx.QueryRow(
		context.Background(),
		"INSERT INTO banners(content, feature_id, is_active) VALUES ($1, $2, $3) RETURNING id, created_at",
		banner.Content,
		banner.FeatureId,
		banner.IsActive,
	).Scan(&bannerID, &createdAt)
	if err != nil {
		return 0, err
	}

	// insert into banner_versions table
	tags, _ := json.Marshal(banner.TagIds)
	fTags := strings.Trim(string(tags), "[]")
	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO banner_version(feature_id, banner_id, version, content, created_at, tags) VALUES ($1, $2, $3, $4, $5, $6)",
		banner.FeatureId,
		bannerID,
		1, // Version 1
		banner.Content,
		createdAt,
		fTags,
	)
	if err != nil {
		return 0, err
	}

	// map created banner with every tag id specified
	for _, tagID := range banner.TagIds {
		_, err = tx.Exec(
			context.Background(),
			"INSERT INTO banners_tags(banner_id, tag_id) VALUES ($1, $2)",
			bannerID,
			tagID,
		)
		if err != nil {
			return 0, err
		}
	}

	return bannerID, nil
}



