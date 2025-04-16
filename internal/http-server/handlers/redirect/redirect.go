package redirect

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	resp "github.com/reybrally/REST-API-app/internal/lib/api/response"
	"github.com/reybrally/REST-API-app/internal/lib/logger/sl"
	"github.com/reybrally/REST-API-app/internal/storage"
)

// URLGetter is an interface for getting url by alias.
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			w.WriteHeader(http.StatusNotFound) // 404 Not Found
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("got url", slog.String("url", resURL))
		// redirect to found url with HTTP статусом 302
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
