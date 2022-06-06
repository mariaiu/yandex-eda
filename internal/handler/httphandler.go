package handler

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/mariaiu/yandex-eda-app/internal/config"
	m "github.com/mariaiu/yandex-eda-app/internal/models"
	"github.com/mariaiu/yandex-eda-app/internal/validator"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type HTTPHandler struct {
	srv Service
	auth *Auth
	logger *logrus.Logger
	validator *validating.Validator
	latitude float64
	longitude float64
	maxWorkers int
}

type Auth struct {
	name string
	password string
}

func NewHTTPHandler(srv Service, logger *logrus.Logger, app *config.App, auth *config.Auth, validator *validating.Validator) *HTTPHandler {
	return &HTTPHandler{
		srv: srv,
		auth: &Auth{
			name: auth.Name,
			password: auth.Password,
		},
		logger: logger,
		validator: validator,
		latitude: app.Latitude,
		longitude: app.Longitude,
		maxWorkers: app.MaxWorkers,
	}
}

func (h *HTTPHandler) ConfigureRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(h.setRequestID)
	router.Use(h.logRequest)
	router.HandleFunc("/restaurant", h.handleGetRestaurants).Methods(http.MethodGet)
	router.HandleFunc("/restaurant/{id:[0-9]+}", h.handleGetRestaurant).Methods(http.MethodGet)
	router.HandleFunc("/parse", h.basicAuth(h.handleParseRestaurants)).Methods(http.MethodGet)

	return router
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type ctxKey struct{}

func (h *HTTPHandler) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, id)))
	})
}

func (h *HTTPHandler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := h.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKey{}),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}

		logger.Logf(level, "completed with %d %s in %v", rw.code, http.StatusText(rw.code), time.Now().Sub(start))
	})
}

func (h *HTTPHandler) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, password, ok := r.BasicAuth()
		if ok {
			nameHash := sha256.Sum256([]byte(name))
			passwordHash := sha256.Sum256([]byte(password))
			expectedNameHash := sha256.Sum256([]byte(h.auth.name))
			expectedPasswordHash := sha256.Sum256([]byte(h.auth.password))

			nameMatch := subtle.ConstantTimeCompare(nameHash[:], expectedNameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if nameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		h.error(w, http.StatusUnauthorized, errors.New("not authenticated"))
	}
}

func (h *HTTPHandler) handleGetRestaurants(w http.ResponseWriter, r *http.Request)  {
	resp, err := h.srv.GetRestaurants(r.Context()); if err != nil {
		h.error(w, http.StatusInternalServerError, err)
		return
	}

	if len(resp) == 0 {
		h.error(w, http.StatusNotFound, errors.New("no restaurants found"))
		return
	}

	h.respond(w, http.StatusOK, resp)
}

func (h *HTTPHandler) handleGetRestaurant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"]); if err != nil {
		h.error(w, http.StatusBadRequest, err)
	}

	resp, err := h.srv.GetPositions(r.Context(), id); if err != nil {
		h.error(w, http.StatusInternalServerError, err)
	}

	if len(resp) == 0 {
		h.error(w, http.StatusNotFound, errors.New("restaurant not found"))
		return
	}

	h.respond(w, http.StatusOK, resp)
}

func (h *HTTPHandler) handleParseRestaurants(w http.ResponseWriter, r *http.Request) {
	req := m.ParseRequest{
		Longitude: h.longitude,
		Latitude:  h.latitude,
		Workers:   h.maxWorkers,
	}

	decoder := schema.NewDecoder()

	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		h.error(w, http.StatusBadRequest, err)
		return
	}

	if err := h.validator.ValidateStruct(req); err != nil {
		h.error(w, http.StatusUnprocessableEntity, err)
		return
	}

	count, err := h.srv.ParseRestaurants(r.Context(), &req); if err != nil {
		switch err.Error() {
		case "process timeout":
			h.error(w, http.StatusGatewayTimeout, err)
			return
		default:
			h.error(w, http.StatusInternalServerError, err)
			return
		}
	}

	h.respond(w, http.StatusOK, map[string]int{"were_processed": count})

}

func (h *HTTPHandler) error(w http.ResponseWriter, code int, err error) {
	h.respond(w, code, map[string]string{"error": err.Error()})
}

func(h *HTTPHandler) respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

