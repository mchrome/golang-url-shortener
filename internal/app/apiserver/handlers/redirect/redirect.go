package redirect

import (
	"fmt"
	"log/slog"
	"net/http"

	resp "github.com/mchrome/url-compression-api/internal/app/lib/api/response"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/mchrome/url-compression-api/internal/app/lib/logger/sl"
)

//go:generate go run github.com/vektra/mockery/v2@v2.36.1 --name URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log = log.With(slog.String("request_id", middleware.GetReqID(r.Context())))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("URL did not contain an alias")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}

		log.Info("alias recieved from URL", slog.Any("alias", alias))

		url, err := urlGetter.GetURL(alias)
		if err != nil {
			log.Error("could not retrieve url from storage", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info(fmt.Sprintf("retrieved URL:%s by alias:%s from storage", url, alias))

		http.Redirect(w, r, url, http.StatusFound)

	}
}
