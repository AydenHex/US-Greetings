// endpoint.go
package greeting

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all endpoits compose the greeting service
type GreetingEndpoints struct {
	GetAllEndpoint  endpoint.Endpoint
	GetByIDEndpoint endpoint.Endpoint
	AddEndpoint     endpoint.Endpoint
	UpdateEndpoint  endpoint.Endpoint
	DeleteEndpoint  endpoint.Endpoint
}

// MakeGreetingEndpoints returns an Endpoints struct where each endpoint invoke
// the corresponding method on the provided greeting
func MakeGreetingEndpoints(s GreetingService) GreetingEndpoints {
	return GreetingEndpoints{
		GetAllEndpoint:  MakeGetAllEndpoint(s),
		GetByIDEndpoint: MakeGetByIDEndpoint(s),
		AddEndpoint:     MakeAddEndpoint(s),
		UpdateEndpoint:  MakeUpdateEndpoint(s),
		DeleteEndpoint:  MakeDeleteEndpoint(s),
	}
}

type GetAllRequest struct {
}
type GetAllResponse struct {
	greetings []Greeting `json:"greetings"`
}

func MakeGetAllEndpoint(s GreetingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		greetings, err := s.GetAll(ctx)
		return GetAllResponse{greetings}, err
	}
}

type GetByIDRequest struct {
	ID string
}

type GetByIDResponse struct {
	greeting Greeting `json:"greeting"`
}

func MakeGetByIDEndpoint(s GreetingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		greeting, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{greeting}, err
	}
}

type AddRequest struct {
	greeting Greeting
}

type AddResponse struct {
	greeting Greeting `json:"greeting"`
}

func MakeAddEndpoint(s GreetingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddRequest)
		greeting, err := s.Add(ctx, req.greeting)
		fmt.Println(AddResponse{greeting})
		return AddResponse{greeting}, err
	}
}

type UpdateRequest struct {
	ID       string
	greeting Greeting
}

type UpdateResponse struct {
}

func MakeUpdateEndpoint(s GreetingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		err := s.Update(ctx, req.ID, req.greeting)
		return UpdateResponse{}, err
	}
}

type DeleteRequest struct {
	ID string
}

type DeleteResponse struct {
}

func MakeDeleteEndpoint(s GreetingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		err := s.Delete(ctx, req.ID)
		return DeleteResponse{}, err
	}
}
