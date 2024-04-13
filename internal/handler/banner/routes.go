package banner

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/mBayzigitov/dynamic-content-service/internal/dto"
	"github.com/mBayzigitov/dynamic-content-service/internal/service"
	"github.com/mBayzigitov/dynamic-content-service/internal/util/serverr"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	TagIdParam           = "tag_id"
	FeatureIdParam       = "feature_id"
	UseLastRevisionParam = "use_last_revision"
	BannerIdPathVariable = "bannerId"
)

type BannerHandler struct {
	valid   *validator.Validate
	l       *zap.SugaredLogger
	service *service.BannerService
}

func NewHandler(service *service.BannerService) *BannerHandler {
	loginst, _ := zap.NewDevelopment()
	return &BannerHandler{
		valid:   validator.New(),
		l:       loginst.Sugar(),
		service: service,
	}
}

func (bh *BannerHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user_banner", bh.handleBannerGetting).Methods("GET")
	router.HandleFunc("/banner", bh.handleBannerCreation).Methods("POST")
	router.HandleFunc("/banner/{bannerId}", bh.handleBannerDeletion).Methods("DELETE")
	router.HandleFunc("/banner/{bannerId}", bh.handleBannerChange).Methods("PATCH")
}

// -------- Helper functions --------
func (bh *BannerHandler) adminOnlyAccess(r *http.Request) *serverr.ApiError {
	isAdmin, ok := r.Context().Value("isAdmin").(bool)
	if !ok {
		bh.l.Fatal(serverr.TokenParsingError)
		return serverr.TokenParsingError
	}

	if !isAdmin {
		bh.l.Info(serverr.ForbiddenAccessError.Error())
		return serverr.ForbiddenAccessError
	}

	return nil
}

// -------- Handler functions --------
func (bh *BannerHandler) handleBannerGetting(w http.ResponseWriter, r *http.Request) {
	// parse params
	ti := r.URL.Query().Get(TagIdParam)
	fi := r.URL.Query().Get(FeatureIdParam)
	ulr := r.URL.Query().Get(UseLastRevisionParam)
	var err error

	var tagId, featureId int64
	var useLastRevision bool

	tagId, err = strconv.ParseInt(ti, 10, 64)
	if ti == "" || err != nil {
		apierror := serverr.NewInvalidRequestError("Некорректное значение tag_id")
		bh.l.Info(apierror.Error())
		http.Error(w, apierror.JsonBody(), apierror.HttpStatus)
		return
	}

	featureId, err = strconv.ParseInt(fi, 10, 64)
	if fi == "" || err != nil {
		apierror := serverr.NewInvalidRequestError("Некорректное значение feature_id")
		bh.l.Info(apierror.Error())
		http.Error(w, apierror.JsonBody(), apierror.HttpStatus)
		return
	}

	if ulr == "" {
		useLastRevision = false
	} else {
		useLastRevision, err = strconv.ParseBool(ulr)

		if err != nil {
			apierror := serverr.NewInvalidRequestError("Некорректное значение use_last_revision")
			bh.l.Info(apierror.Error())
			http.Error(w, apierror.JsonBody(), apierror.HttpStatus)
			return
		}
	}

	if resp, apierr := bh.service.GetBanner(tagId, featureId, useLastRevision); apierr != nil {
		http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
	} else {
		w.WriteHeader(http.StatusOK)
		jsonBody := dto.JsonBody(dto.NewGetBannerResponse(&resp))
		w.Write([]byte(jsonBody))
	}
}

func (bh *BannerHandler) handleBannerCreation(w http.ResponseWriter, r *http.Request) {
	accessErr := bh.adminOnlyAccess(r)
	if accessErr != nil {
		http.Error(w, accessErr.JsonBody(), accessErr.HttpStatus)
		return
	}

	var rb dto.CreateBannerDto
	if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
		apierr := serverr.InvalidRequestError
		bh.l.Info(apierr.Error())
		http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
		return
	}

	err := rb.Validate(bh.valid)
	if err != nil {
		bh.l.Info(err.Error())
		http.Error(w, err.JsonBody(), err.HttpStatus)
		return
	}

	if createdId, apierr := bh.service.CreateBanner(rb.ToModel()); apierr != nil {
		http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
	} else {
		w.WriteHeader(201)
		resp := dto.JsonBody(dto.NewCreateBannerResponse(createdId))
		w.Write([]byte(resp))
		bh.l.Infof("Banner [id=%d] is created", createdId)
	}
}

func (bh *BannerHandler) handleBannerDeletion(w http.ResponseWriter, r *http.Request) {
	accessErr := bh.adminOnlyAccess(r)
	if accessErr != nil {
		http.Error(w, accessErr.JsonBody(), accessErr.HttpStatus)
		return
	}

	qp := mux.Vars(r)
	var bannerId int64
	var err error

	// check whether path param exists & has correct value
	if bi, ok := qp[BannerIdPathVariable]; !ok {
		apierr := serverr.NewInvalidRequestError("Отсутствует параметр 'bannerId'")
		bh.l.Info(apierr)
		http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
		return
	} else {
		bannerId, err = strconv.ParseInt(bi, 10, 64)
		if err != nil {
			apierr := serverr.NewInvalidRequestError("Неверный формат параметра 'bannerId'")
			bh.l.Info(apierr)
			http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
			return
		}
	}

	// call service method and return response
	if apierr := bh.service.DeleteBanner(bannerId); apierr != nil {
		http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
	} else {
		w.WriteHeader(204)
	}
}

func (bh *BannerHandler) handleBannerChange(w http.ResponseWriter, r *http.Request) {
	accessErr := bh.adminOnlyAccess(r)
	if accessErr != nil {
		http.Error(w, accessErr.JsonBody(), accessErr.HttpStatus)
		return
	}

	qp := mux.Vars(r)
	var bannerId int64
	var err error

	// check whether path param exists & has correct value
	if bi, ok := qp[BannerIdPathVariable]; !ok {
		apierr := serverr.NewInvalidRequestError("Отсутствует параметр 'bannerId'")
		bh.l.Info(apierr)
		http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
		return
	} else {
		bannerId, err = strconv.ParseInt(bi, 10, 64)
		if err != nil {
			apierr := serverr.NewInvalidRequestError("Неверный формат параметра 'bannerId'")
			bh.l.Info(apierr)
			http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
			return
		}
	}

	var cb dto.ChangeBannerDto
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&cb); err != nil {
		apierr := serverr.InvalidRequestError
		bh.l.Info(err)
		http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
		return
	}

	// call service method and return response
	if apierr := bh.service.ChangeBanner(bannerId, cb); apierr != nil {
		http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
	} else {
		w.WriteHeader(200)
	}
}
