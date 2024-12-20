package createj

import (
	"context"
	"errors"
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/lib/api/request"
	resp "github.com/Igorezka/shortener/internal/app/lib/api/response"
	ci "github.com/Igorezka/shortener/internal/app/lib/cipher"
	"github.com/Igorezka/shortener/internal/app/storage"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	resp.Response
	Result string `json:"result"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.49.0 --name=URLSaver
type URLSaver interface {
	SaveURL(ctx context.Context, url string, userID string) (string, error)
}

func New(log *zap.Logger, cfg *config.Config, cipher *ci.Cipher, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.create_json.New"
		log := log.With(
			zap.String("op", op),
		)

		token, err := r.Cookie("token")
		var userID string
		if err != nil {
			userID = ""
		} else {
			userID, err = cipher.Open(token.Value)
			if err != nil {
				log.Error("failed to decode token", zap.String("error", err.Error()))
				resp.Status(r, http.StatusInternalServerError)
				resp.JSON(w, r, resp.Error("Unauthorized"))
				return
			}
		}

		var req Request
		err = request.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", zap.String("error", err.Error()))
			resp.Status(r, http.StatusBadRequest)
			resp.JSON(w, r, resp.Error("internal server error"))
			return
		}

		if len(req.URL) <= 0 {
			log.Info("url field required")
			resp.Status(r, http.StatusBadRequest)
			resp.JSON(w, r, resp.Error("url required"))
			return
		}

		if _, err := url.ParseRequestURI(req.URL); err != nil {
			log.Info("invalid url", zap.String("url", req.URL))
			resp.Status(r, http.StatusBadRequest)
			resp.JSON(w, r, resp.Error("only valid url required"))
			return
		}

		id, err := urlSaver.SaveURL(r.Context(), req.URL, userID)
		if err != nil {
			log.Error("failed to store link", zap.String("error", err.Error()))
			if errors.Is(err, storage.ErrURLExist) {
				resp.Status(r, http.StatusConflict)
				resp.JSON(w, r, Response{
					Response: resp.OK(),
					Result:   cfg.BaseURL + "/" + id,
				})
				return
			}
			resp.Status(r, http.StatusInternalServerError)
			resp.JSON(w, r, resp.Error("internal server error"))
			return
		}

		resp.Status(r, http.StatusCreated)
		resp.JSON(w, r, Response{
			Response: resp.OK(),
			Result:   cfg.BaseURL + "/" + id,
		})
	}
}
