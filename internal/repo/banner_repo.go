package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mBayzigitov/dynamic-content-service/internal/dto"
	"github.com/mBayzigitov/dynamic-content-service/internal/models"
	"github.com/mBayzigitov/dynamic-content-service/internal/util"
	"github.com/mBayzigitov/dynamic-content-service/internal/util/serverr"
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

func (br *BannerRepository) DeleteMarkedBanners() error {
	// start a transaction
	tx, err := br.p.Begin(context.Background())
	if err != nil {
		br.l.Fatal(err)
		return err
	}
	defer func() {
		if pm := recover(); pm != nil {
			tx.Rollback(context.Background())
			panic(pm)
		} else if err != nil {
			br.l.Fatal(err)
			tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()

	// delete banners marked as to_delete
	_, err = tx.Exec(context.Background(), `
		DELETE FROM banners
		WHERE to_delete = true
	`)
	if err != nil {
		br.l.Fatal(err)
		return err
	}

	br.l.Info("Banners marked as to_delete have been cleaned up successfully")

	return nil
}

// CheckIfDuplicates
// Method that checks if a bunch of key (banner_id-feature_id-tag_id) already
// exists to satisfy the condition of unambigious definition
func (br *BannerRepository) CheckIfDuplicates(featureId int64, tagsIds []int64) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM banners b
			 JOIN public.banners_tags bt on b.id = bt.banner_id
		WHERE
		     b.feature_id = $1
		 and bt.tag_id IN (`

	params := make([]string, len(tagsIds))
	for i, id := range tagsIds {
		params[i] = strconv.FormatInt(id, 10)
	}
	query += strings.Join(params, ",") + ")"

	var count int
	err := br.p.QueryRow(
		context.Background(),
		query,
		featureId,
	).Scan(&count)
	if err != nil {
		return false, err
	}

	// check if at least one record exists
	return count > 0, nil
}

func (br *BannerRepository) GetBannerByTagAndFeature(tagId int64, featureId int64) (models.BannerModel, error) {
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
		"INSERT INTO banners(content, feature_id) VALUES ($1, $2) RETURNING id, created_at",
		banner.Content,
		banner.FeatureId,
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

func (br *BannerRepository) DeleteBanner(bannerId int64) *serverr.ApiError {
	// start a transaction
	tx, err := br.p.Begin(context.Background())
	if err != nil {
		br.l.Fatal(err)
		return serverr.StorageError
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

	// perform the update operation to set to_delete=true
	result, err := tx.Exec(
		context.Background(),
		"UPDATE banners SET to_delete=true WHERE id = $1 AND is_active = true",
		bannerId,
	)
	if err != nil {
		br.l.Error(err)
		return serverr.StorageError
	}

	// check if the update affected any rows
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		// no rows were affected, indicating that the banner with the given id doesn't exist
		return serverr.BannerNotFoundError
	}

	br.l.Infof("Banner [id=%d] has been marked as deleted successfully", bannerId)

	return nil
}

func (br *BannerRepository) ChangeBannerByRequest(bannerId int64, chban dto.ChangeBannerDto) *serverr.ApiError {
	// key idea is to get existing banner and use it as pattern for changes
	// check if banner is present, if it is -> get banner pattern, change updated_at
	bannerPattern, err := br.GetBannerById(bannerId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	bannerPattern.UpdatedAt = time.Now()
	fmt.Println(*bannerPattern)

	// if featureId NOT NULL -> check featureId, if exists -> change it in banner pattern
	if chban.FeatureId != nil {
		featExists, err := br.DoesFeatureExist(*chban.FeatureId)
		if err != nil {
			br.l.Error(err.Error())
			return serverr.StorageError
		}

		if featExists {
			bannerPattern.FeatureId = *chban.FeatureId
		} else {
			return serverr.NewInvalidRequestError("Указанный feature_id не существует")
		}
	}

	// if tagIds NOT NULL -> check whether tagIds exist, if exists -> change it in banner pattern
	if len(chban.TagIds) != 0 {
		tagsExist, err := br.DoTagsExist(chban.TagIds)
		if err != nil {
			br.l.Error(err.Error())
			return serverr.StorageError
		}

		if tagsExist {
			bannerPattern.TagIds = chban.TagIds
		} else {
			return serverr.NewInvalidRequestError("Все/некоторые tag_id не существуют")
		}
	}

	// if content NOT NULL -> assign value
	if chban.Content != nil {
		bannerPattern.Content = *chban.Content
	}

	// if is_active NOT NULL -> assign value, else keep value
	if chban.IsActive != nil {
		bannerPattern.IsActive = *chban.IsActive
	}

	// start transaction, commit through defer
	tx, txerr := br.p.Begin(context.Background())
	if txerr != nil {
		return serverr.StorageError
	}
	defer func() {
		if pm := recover(); pm != nil {
			tx.Rollback(context.Background())
			panic(pm)
		} else if err != nil {
			tx.Rollback(context.Background())
		} else {
			txerr = tx.Commit(context.Background())
		}
	}()

	// create new version, get last revision param
	tags, _ := json.Marshal(bannerPattern.TagIds)
	fTags := strings.Trim(string(tags), "[]")
	_, txerr = tx.Exec(
		context.Background(),
		"INSERT INTO banner_version(feature_id, banner_id, version, content, created_at, tags) VALUES ($1, $2, $3, $4, $5, $6)",
		bannerPattern.FeatureId,
		bannerId,
		bannerPattern.LastRevision+1,
		bannerPattern.Content,
		bannerPattern.UpdatedAt, // because version is created when main banner is updated
		fTags,
	)

	// delete mapped tags, map new tags
	err = br.RewriteBannerTags(bannerId, bannerPattern.TagIds)
	if err != nil {
		return err
	}

	// change the banner itself
	bannerPattern.LastRevision = bannerPattern.LastRevision + 1
	err = br.ChangeBanner(bannerId, bannerPattern)
	if err != nil {
		return err
	}

	br.l.Infof("Banner [id=%d] is updated successfully", bannerId)

	return nil
}

func (br *BannerRepository) GetBannerById(bannerId int64) (*models.BannerTagsModel, *serverr.ApiError) {
	// query to get banner details from the banners table
	row := br.p.QueryRow(
		context.Background(),
		"SELECT feature_id, content, is_active, created_at, updated_at, last_revision, to_delete FROM banners WHERE id = $1",
		bannerId,
	)

	// Initialize variables to store banner details
	var banner models.BannerTagsModel

	// Scan the banner details into the struct
	err := row.Scan(
		&banner.FeatureId,
		&banner.Content,
		&banner.IsActive,
		&banner.CreatedAt,
		&banner.UpdatedAt,
		&banner.LastRevision,
		&banner.ToDelete,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, serverr.BannerNotFoundError
		}
		br.l.Error(err)
		return nil, serverr.StorageError
	}

	// get tag IDs associated with the banner from the banners_tags table
	rows, err := br.p.Query(
		context.Background(),
		`SELECT tag_id 
			 FROM banners_tags
			 WHERE banner_id = $1`,
		bannerId,
	)
	if err != nil {
		return nil, serverr.StorageError
	}
	defer rows.Close()

	// iterate through the rows and append tag IDs to the slice
	for rows.Next() {
		var tagID int64
		if err := rows.Scan(&tagID); err != nil {
			br.l.Error(err)
			return nil, serverr.StorageError
		}
		banner.TagIds = append(banner.TagIds, tagID)
	}

	// check for errors during row iteration
	if err := rows.Err(); err != nil {
		br.l.Error(err)
		return nil, serverr.StorageError
	}

	return &banner, nil
}

func (br *BannerRepository) RewriteBannerTags(bannerId int64, tagIds []int64) *serverr.ApiError {
	// start a transaction
	tx, txerr := br.p.Begin(context.Background())
	if txerr != nil {
		br.l.Error(txerr)
		return serverr.StorageError
	}
	defer func() {
		if pm := recover(); pm != nil {
			tx.Rollback(context.Background())
			panic(pm)
		} else if txerr != nil {
			tx.Rollback(context.Background())
		} else {
			txerr = tx.Commit(context.Background())
		}
	}()

	// delete existing banners_tags records for the given bannerId
	_, txerr = tx.Exec(
		context.Background(),
		"DELETE FROM banners_tags WHERE banner_id = $1",
		bannerId,
	)
	if txerr != nil {
		br.l.Error(txerr)
		return serverr.StorageError
	}

	// insert new banners_tags records
	for _, tagId := range tagIds {
		_, txerr = tx.Exec(
			context.Background(),
			"INSERT INTO banners_tags (banner_id, tag_id) VALUES ($1, $2)",
			bannerId,
			tagId,
		)
		if txerr != nil {
			br.l.Error(txerr)
			return serverr.StorageError
		}
	}

	return nil
}

func (br *BannerRepository) ChangeBanner(bannerId int64, chban *models.BannerTagsModel) *serverr.ApiError {
	// start a transaction
	tx, txerr := br.p.Begin(context.Background())
	if txerr != nil {
		br.l.Error(txerr)
		return serverr.StorageError
	}
	defer func() {
		if pm := recover(); pm != nil {
			tx.Rollback(context.Background())
			panic(pm)
		} else if txerr != nil {
			br.l.Error(txerr)
			tx.Rollback(context.Background())
		} else {
			txerr = tx.Commit(context.Background())
		}
	}()

	// update the fields in the banners table
	_, txerr = tx.Exec(
		context.Background(),
		`UPDATE banners 
			 SET content = $1, 
			     feature_id = $2, 
			     is_active = $3, 
			     updated_at = $4, 
			     to_delete = $5,
			     last_revision = $6
			 WHERE id = $7`,
		chban.Content,
		chban.FeatureId,
		chban.IsActive,
		chban.UpdatedAt,
		chban.ToDelete,
		chban.LastRevision,
		bannerId,
	)
	if txerr != nil {
		br.l.Error(txerr)
		return serverr.StorageError
	}

	return nil
}

func (br *BannerRepository) GetBannersByFilter(featureId int64, tagId int64, limit int64, offset int64) ([]models.BannerTagsModel, *serverr.ApiError) {
	// construct query logic
	var featureQp, andQp, tagIdQp string

	both := featureId != 0 && tagId != 0

	if featureId != 0 {
		featureQp = fmt.Sprintf("b.feature_id = %d", featureId)
	}

	if both {
		andQp = "and"
	}

	if tagId != 0 {
		tagIdQp = fmt.Sprintf(
			`b.id in (
        		select btt.banner_id
				from banners_tags btt
        		where btt.tag_id = %d
      		)`,
			tagId)
	}

	query := fmt.Sprintf(
		`
		SELECT b.id,
			   bt.tag_id,
			   b.feature_id,
			   b.content,
			   b.is_active,
			   b.to_delete,
			   b.created_at,
			   b.updated_at
		FROM banners b
		JOIN banners_tags bt on b.id = bt.banner_id
		WHERE
		%s
		%s
		%s
		`,
		featureQp,
		andQp,
		tagIdQp,
	)

	// exec query
	rows, err := br.p.Query(context.Background(), query)
	if err != nil {
		br.l.Error(err)
		return nil, serverr.StorageError
	}
	defer rows.Close()

	var banners []models.BannerTagsModel
	var curModel models.BannerTagsModel
	for rows.Next() {
		var banner models.BannerModel
		err := rows.Scan(
			&banner.Id,
			&banner.TagId,
			&banner.FeatureId,
			&banner.Content,
			&banner.IsActive,
			&banner.ToDelete,
			&banner.CreatedAt,
			&banner.UpdatedAt,
		)
		if err != nil {
			return nil, serverr.StorageError
		}

		if curModel.Id == 0 {
			curModel = models.BannerTagsModel{
				Id:        banner.Id,
				FeatureId: banner.FeatureId,
				Content:   banner.Content,
				IsActive:  banner.IsActive,
				ToDelete:  banner.ToDelete,
				CreatedAt: banner.CreatedAt,
				UpdatedAt: banner.UpdatedAt,
			}
			curModel.TagIds = append(curModel.TagIds, banner.TagId)
		} else if banner.Id != curModel.Id {
			banners = append(banners, curModel)

			curModel = models.BannerTagsModel{
				Id:        banner.Id,
				FeatureId: banner.FeatureId,
				Content:   banner.Content,
				IsActive:  banner.IsActive,
				ToDelete:  banner.ToDelete,
				CreatedAt: banner.CreatedAt,
				UpdatedAt: banner.UpdatedAt,
			}
			curModel.TagIds = append(curModel.TagIds, banner.TagId)
		} else {
			curModel.TagIds = append(curModel.TagIds, banner.TagId)
		}
	}

	if curModel.Id != 0 {
		banners = append(banners, curModel)
	}

	// create boundaries
	var start, end int64
	blength := int64(len(banners))

	// check if offset is greater than the length of the slice, set it to 0
	if offset >= blength {
		start = blength
	} else {
		start = offset
	}

	// extract the desired elements from the slice
	result := banners[start:]

	// check if limit is 0, set it to the length of the slice
	if limit == 0 || limit > int64(len(result)) {
		end = int64(len(result))
	} else {
		end = limit
	}

	result = result[:end]

	return result, nil
}

func (br *BannerRepository) DeleteBannersByTagOrFeatureId(featureId int64, tagId int64) *serverr.ApiError {
	tx, txerr := br.p.Begin(context.Background())
	if txerr != nil {
		br.l.Error(txerr)
		return serverr.StorageError
	}
	defer func() {
		if pm := recover(); pm != nil {
			tx.Rollback(context.Background())
			panic(pm)
		} else if txerr != nil {
			br.l.Error(txerr)
			tx.Rollback(context.Background())
		} else {
			txerr = tx.Commit(context.Background())
		}
	}()

	var query string
	var param int64
	if featureId != 0 {
		featureExists, err := br.DoesFeatureExist(featureId)
		if err != nil {
			br.l.Error(err)
			return serverr.StorageError
		}

		if !featureExists {
			return serverr.NewInvalidRequestError("Указанный feature_id не существует")
		}

		param = featureId
		query = `
			UPDATE banners
        	SET to_delete = true
        	WHERE feature_id = $1
        `
		defer br.l.Infof("Banners with feature_id=%d were marked as deleted", param)
	} else {
		tagsExist, err := br.DoTagsExist([]int64{tagId})
		if err != nil {
			br.l.Error(err.Error())
			return serverr.StorageError
		}

		if !tagsExist {
			return serverr.NewInvalidRequestError("Указанный tag_id не существует")
		}

		param = tagId
		query = `
			UPDATE banners
        	SET to_delete = true
        	WHERE id IN (
					SELECT banner_id
					FROM banners_tags
					WHERE tag_id = $1
        		)`
		defer br.l.Infof("Banners with tag_id=%d were marked as deleted", param)
	}

	_, txerr = tx.Exec(context.Background(), query, param)
	if txerr != nil {
		br.l.Error(txerr)
		return serverr.StorageError
	}

	return nil
}

func (br *BannerRepository) GetBannerVersions(bannerId int64) ([]models.BannerVersion, *serverr.ApiError) {
	rows, err := br.p.Query(
		context.Background(),
		`SELECT bv.banner_id,
       				bv.version,
       				bv.feature_id,
       				bv.tags,
       				bv.content,
       				bv.created_at
			 FROM banner_version bv
			 WHERE banner_id = $1`,
		bannerId,
	)
	if err != nil {
		return nil, serverr.StorageError
	}
	defer rows.Close()

	var versions []models.BannerVersion
	for rows.Next() {
		var c models.BannerVersion
		if err := rows.Scan(&c.BannerId, &c.Version, &c.FeatureId, &c.Tags, &c.Content, &c.CreatedAt); err != nil {
			br.l.Error(err)
			return nil, serverr.StorageError
		}
		versions = append(versions, c)
	}

	if len(versions) == 0 {
		return []models.BannerVersion{}, serverr.BannerNotFoundError
	}

	return versions, nil
}

func (br *BannerRepository) SetBannerVersion(bannerId int64, versionId int64) *serverr.ApiError {
	fmt.Println(bannerId, versionId)

	// check if the specified version of the banner exists
	query := `
        SELECT bv.banner_id,
               bv.version,
               bv.feature_id,
               bv.tags,
               bv.content,
               bv.created_at,
               b.is_active,
               b.to_delete,
               b.id,
               b.created_at
        FROM banner_version bv
        	 JOIN banners b on bv.banner_id = b.id
        WHERE banner_id = $1 AND version = $2
    `
	var version models.BannerVersion
	var chban models.BannerTagsModel
	err := br.p.QueryRow(
		context.Background(),
		query,
		bannerId,
		versionId).Scan(
			&version.BannerId,
			&version.Version,
			&version.FeatureId,
			&version.Tags,
			&version.Content,
			&version.CreatedAt,
			&chban.IsActive,
			&chban.ToDelete,
			&chban.Id,
			&chban.CreatedAt)
	if err != nil {
		br.l.Error(err)
		return serverr.BannerNotFoundError
	}

	tx, txerr := br.p.Begin(context.Background())
	if txerr != nil {
		br.l.Error(txerr)
		return serverr.StorageError
	}
	defer func() {
		if pm := recover(); pm != nil {
			tx.Rollback(context.Background())
			panic(pm)
		} else if txerr != nil {
			br.l.Error(txerr)
			tx.Rollback(context.Background())
		} else {
			txerr = tx.Commit(context.Background())
		}
	}()

	// retrieve a slice from string of tagIds
	chban.TagIds, err = util.StringToIntSlice(version.Tags)

	chban.Content = version.Content
	chban.FeatureId = version.FeatureId

	// change updated_at because technically its updated now
	chban.UpdatedAt = time.Now()
	chban.LastRevision = versionId

	apierr := br.ChangeBanner(bannerId, &chban)
	if apierr != nil {
		return serverr.StorageError
	}

	apierr = br.RewriteBannerTags(bannerId, chban.TagIds)
	if apierr != nil {
		return serverr.StorageError
	}

	apierr = br.DeleteVersionsGreaterThan(versionId)
	if apierr != nil {
		return apierr
	}

	return nil
}

func (br *BannerRepository) DeleteVersionsGreaterThan(versionId int64) *serverr.ApiError {
	tx, txerr := br.p.Begin(context.Background())
	if txerr != nil {
		br.l.Error(txerr)
		return serverr.StorageError
	}
	defer func() {
		if pm := recover(); pm != nil {
			tx.Rollback(context.Background())
			panic(pm)
		} else if txerr != nil {
			br.l.Error(txerr)
			tx.Rollback(context.Background())
		} else {
			txerr = tx.Commit(context.Background())
		}
	}()

	sqlStatement := `
        DELETE FROM banner_version
        WHERE version > $1
    `

	_, txerr = tx.Exec(context.Background(), sqlStatement, versionId)
	if txerr != nil {
		return serverr.StorageError
	}

	return nil
}