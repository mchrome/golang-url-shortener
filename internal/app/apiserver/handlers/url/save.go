package save

import (
	"fmt"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/mchrome/url-compression-api/internal/app/lib/api/response"
	"github.com/mchrome/url-compression-api/internal/app/lib/logger/sl"
	"github.com/mchrome/url-compression-api/internal/app/lib/random"
)

const aliasLen = 8

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.36.1 --name URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log = log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("could not decode request's JSON body", sl.Err(err))
			render.JSON(w, r, resp.Error("could not decode request's JSON body"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {

			validateErr := err.(validator.ValidationErrors)
			log.Error("request invalid", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomLatinString(aliasLen)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			log.Error("could not add url", sl.Err(err))
			render.JSON(w, r, resp.Error("could not add url"))
			return
		}

		log.Info(fmt.Sprintf("url added (id: %d, url: %s, alias: %s)", id, req.URL, alias))

		render.JSON(w, r, &Response{
			Response: resp.OK(),
			Alias:    alias,
		})

	}
}
