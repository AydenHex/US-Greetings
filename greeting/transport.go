// transport.go
package greeting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	httptransport "github.com/go-kit/kit/transport/http"
)

var ErrMissingParam = errors.New("Missing parameter")

func MakeHTTPHandler(endpoints GreetingEndpoints) http.Handler {

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.StripSlashes)

	greetingRouter := chi.NewRouter()

	/**greetingRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))*/

	greetingRouter.Get("/", httptransport.NewServer(
		endpoints.GetAllEndpoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	greetingRouter.Get("/{id}", httptransport.NewServer(
		endpoints.GetByIDEndpoint,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	greetingRouter.Post("/", httptransport.NewServer(
		endpoints.AddEndpoint,
		decodeAddRequest,
		encoreResponseTest,
		options...,
	).ServeHTTP)

	greetingRouter.Put("/{id}", httptransport.NewServer(
		endpoints.UpdateEndpoint,
		decodeUpdateRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	greetingRouter.Delete("/{id}", httptransport.NewServer(
		endpoints.DeleteEndpoint,
		decodeDeleteRequest,
		encodeResponse,
		options...,
	).ServeHTTP)

	r.Mount("/greetings", greetingRouter)

	return r
}

func decodeGetRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return GetAllRequest{}, err
}

func decodeGetByIDRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrMissingParam
	}
	return GetByIDRequest{id}, err
}

func decodeAddRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var greeting Greeting
	err = render.Decode(r, &greeting)
	if err != nil {
		return nil, err
	}
	return AddRequest{greeting}, err
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrMissingParam
	}
	var greeting Greeting
	err = render.Decode(r, &greeting)
	if err != nil {
		return nil, err
	}
	return UpdateRequest{id, greeting}, err
}

func decodeDeleteRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, ErrMissingParam
	}
	return DeleteRequest{id}, err
}

func encoreResponseTest(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(error); ok && err != nil {
		encodeError(ctx, err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	fmt.Println(response)
	return json.NewEncoder(w).Encode(&response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInconsistentIDs, ErrMissingParam:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
