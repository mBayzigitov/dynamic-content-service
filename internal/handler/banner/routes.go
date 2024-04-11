package banner

import (
	"encoding/json"
	"fmt"
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

		resp, _ := json.Marshal(apierr)
		http.Error(w, string(resp), apierr.HttpStatus)
		return
	}

	if err := bh.valid.Struct(rb); err != nil {
		errors := fmt.Sprintf("%s", err.(validator.ValidationErrors))
		apierr := serverr.NewInvalidRequestError(errors)
		resp, _ := json.Marshal(apierr)

		http.Error(w, string(resp), apierr.HttpStatus)
		return
	}

	if apierr := bh.service.CreateBanner(rb.TagIds, rb.FeatureId, rb.Content, rb.IsActive); apierr != nil {
		// write apierr to user response
	} else {
		// return 200
	}
}
