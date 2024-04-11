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
)

type BannerHandler struct {
	valid   *validator.Validate
	l       *zap.SugaredLogger
	service *service.BannerService
}

func NewHandler(service *service.BannerService) *BannerHandler {
	loginst, _ := zap.NewDevelopment()
	return &BannerHandler{
		valid: validator.New(),
		l:     loginst.Sugar(),
		service: service,
	}
}

func (bh *BannerHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/banner", bh.handleBannerCreation).Methods("POST")
}

func (bh *BannerHandler) handleBannerCreation(w http.ResponseWriter, r *http.Request) {
	var rb dto.CreateBannerDto
	if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
		apierr := serverr.InvalidRequestError
		defer bh.l.Warn(apierr.ErrType)

		http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
		return
	}

	if err := bh.valid.Struct(rb); err != nil {
		verrs := err.(validator.ValidationErrors)

		var errBody string
		if verrs != nil && len(verrs) > 0 {
			f := verrs[0]
			errBody = "field '" + f.Field() + "' validation failed: '" + f.ActualTag() + "' is violated"
		} else {
			errBody = "validation error"
		}

		apierr := serverr.NewInvalidRequestError(errBody)

		http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
		return
	}

	if createdId, apierr := bh.service.CreateBanner(rb.ToModel()); apierr != nil {
		http.Error(w, apierr.JsonBody(), apierr.HttpStatus)
	} else {
		w.WriteHeader(http.StatusOK)
		resp := dto.NewCreateBannerResponse(createdId).JsonBody()
		w.Write([]byte(resp))
	}
}
